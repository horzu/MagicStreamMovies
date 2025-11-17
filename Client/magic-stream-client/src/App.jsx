import "./App.css";
import Home from "./components/home/Home";
import Header from "./components/header/Header";
import Recommended from "./components/recommended/Recommended";
import Review from "./components/review/Review";
import Register from "./components/register/Register";
import Login from "./components/login/Login";
import { Route, Routes, useNavigate } from "react-router-dom";
import Layout from "./components/Layout";
import RequiredAuth from "./components/RequiredAuth";
import axiosClient from "./api/axiosConfig";
import useAuth from './hook/useAuth';

function App() {
  const navigate = useNavigate();
  const { auth, setAuth } = useAuth();

  const updateMovieReview = (imdb_id) => {
    navigate(`/review/${imdb_id}`);
  };

  const handleLogout = async () => {
    try {
      const response = await axiosClient.post("/logout", {
        user_id: auth.user_id,
      });
      setAuth(null);
    //   localStorage.removeItem("user");
      console.log("user logged out");
    } catch (error) {
      console.error("Logout error:", error);
    }
  };

  return (
    <>
      <Header handleLogout={handleLogout} />
      <Routes path="/" element={Layout}>
        <Route
          element={<Home updateMovieReview={updateMovieReview} />}
          path="/"
        ></Route>
        <Route element={<Register />} path="/register"></Route>
        <Route element={<Login />} path="/login"></Route>
        <Route element={<RequiredAuth />}>
          <Route path="/recommended" element={<Recommended />} />
          <Route path="/review/:imdb_id" element={<Review />} />
        </Route>
      </Routes>
    </>
  );
}

export default App;
