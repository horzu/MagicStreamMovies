import useAxiosPrivate from "../../hook/useAxiosPrivate";
import { useEffect, useState } from "react";
import Movies from "../movies/Movies";

const Recommended = () => {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");
  const axiosPrivate = useAxiosPrivate();

  useEffect(() => {
    const fetchRecommended = async () => {
      setLoading(true);
      setMessage("");

      try {
        const response = await axiosPrivate.get("/recommendedmovies");
        setMovies(response.data);
		console.log(response.data);
      } catch (error) {
        console.error("Error fetching recommended movies:", error.response.data || error);
      } finally {
        setLoading(false);
      }
    };
    fetchRecommended();
  }, []);

  return (
    <>
      {loading ? (
        <h2>Loading recommended movies...</h2>
      ) : (
        <Movies movies={movies} message={message} />
      )}
    </>
  );
};

export default Recommended;
