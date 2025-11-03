package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/horzu/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movies from database"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding movies"})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting movie into database"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := utils.GetRoleFromContext(c)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be part of the ADMIN role"})
			return 
		}

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movieId required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)
		if err != nil {
			log.Println("GetReviewRanking error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking", "details": err.Error()})
			return
		}

		filter := bson.M{"imdb_id": movieId}

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)

	}
}

func GetReviewRanking(admin_review string) (string, int, error) {
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	basePromptTemplate := os.Getenv("BASE_PROMPT_TEMPLATE")
	if basePromptTemplate == "" {
		return "", 0, errors.New("could not read BASE_PROMPT_TEMPLATE")
	}

	basePrompt := strings.Replace(basePromptTemplate, "{rankings}", sentimentDelimited, 1)

	openRouterKey := os.Getenv("OPENROUTER_API_KEY")
    openRouterBase := os.Getenv("OPENROUTER_BASE_URL")
    model := os.Getenv("AI_MODEL")
    if openRouterKey == "" {
        return "", 0, errors.New("could not read OPENROUTER_API_KEY")
    }
    if openRouterBase == "" {
        return "", 0, errors.New("could not read OPENROUTER_BASE_URL")
    }
    if model == "" {
        // fallback if AI_MODEL not set
        model = "gpt-4o-mini"
    }

    // build OpenRouter chat request
    reqBody := map[string]interface{}{
        "model": model,
        "messages": []map[string]string{
            {"role": "user", "content": basePrompt + admin_review},
        },
        "max_tokens": 800,
        "temperature": 0,
    }

    bodyBytes, err := json.Marshal(reqBody)
    if err != nil {
        return "", 0, err
    }

    endpoint := strings.TrimRight(openRouterBase, "/") + "/chat/completions"
    httpReq, err := http.NewRequest("POST", endpoint, bytes.NewReader(bodyBytes))
    if err != nil {
        return "", 0, err
    }
    httpReq.Header.Set("Authorization", "Bearer "+openRouterKey)
    httpReq.Header.Set("Content-Type", "application/json")

    client := &http.Client{Timeout: 20 * time.Second}
    httpResp, err := client.Do(httpReq)
    if err != nil {
        return "", 0, err
    }
    defer httpResp.Body.Close()

    respBody, _ := io.ReadAll(httpResp.Body)
    if httpResp.StatusCode >= 300 {
        return "", 0, errors.New("openrouter API error: " + httpResp.Status + " - " + string(respBody))
    }

    var orResp struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
            Text string `json:"text"`
        } `json:"choices"`
    }
    if err := json.Unmarshal(respBody, &orResp); err != nil {
        return "", 0, err
    }
    if len(orResp.Choices) == 0 {
        return "", 0, errors.New("openrouter: no choices returned")
    }

    var responseText string
    if orResp.Choices[0].Message.Content != "" {
        responseText = strings.TrimSpace(orResp.Choices[0].Message.Content)
    } else {
        responseText = strings.TrimSpace(orResp.Choices[0].Text)
    }

	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == responseText {
			rankVal = ranking.RankingValue
			break
		}
	}
	return responseText, rankVal, nil

}

func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}

		favorite_genres, err := GetUsersFavoriteGenres(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env file not found")
		}

		var recommendedMovieLimitVal int64 = 5
		recommendedMovieLimitValStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")
		if recommendedMovieLimitValStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitValStr, 10, 64)
		}

		findOptions := options.Find()

		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}}).SetLimit(recommendedMovieLimitVal)

		filter := bson.M{"genre.genre_name": bson.M{"$in": favorite_genres}}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedMovies)

	}
}

func GetUsersFavoriteGenres(userId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{
		"favorite_genres.genre_name": 1,
		"_id":                        0,
	}

	opts := options.FindOne().SetProjection(projection)

	var results bson.M
	err := userCollection.FindOne(ctx, filter, opts).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}

	favGenresArray, ok := results["favorite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favorite genres for the user")
	}

	var genreNames []string
	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil
}
