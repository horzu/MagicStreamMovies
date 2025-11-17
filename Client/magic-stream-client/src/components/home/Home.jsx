import { useState, useEffect } from "react";
import axiosClient from "../../api/axiosConfig";
import Movies from "../movies/Movies";

const Home = () => {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState();

  useEffect(() => {
    const fetchMovies = async () => {
      setLoading(true);
      setMessage("");

      try {
        const response = await axiosClient.get("/movies");
        setMovies(response.data);
        if (response.data.length === 0) {
          setMessage("There are currently no movies available");
        }
      } catch (error) {
        console.error("Error fetching the movies", error);
        setMessage("Error fetching the movies");
      } finally {
        setLoading(false);
      }
    }
	fetchMovies();
  }, []);

  return (
    <>
      {loading ? (
        <h2>Text loading...</h2>
      ) : (
        <Movies movies={movies} message={message} />
      )}
    </>
  );
};

export default Home;
