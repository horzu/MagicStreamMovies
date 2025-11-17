import "./App.css";
import Home from "./components/home/Home";
import Header from "./components/header/Header";
import Recommended from "./components/recommended/Recommended";
import Register from "./components/register/Register";
import Login from "./components/login/Login";
import { Route, Routes, useNavigate } from "react-router-dom";
import Layout from "./components/Layout";
import RequiredAuth from "./components/RequiredAuth";

function App() {
  return (
    <>
      <Header />
      <Routes path="/" element={Layout}>
        <Route element={<Home />} path="/"></Route>
        <Route element={<Register />} path="/register"></Route>
        <Route element={<Login />} path="/login"></Route>
        <Route element={<RequiredAuth />}>
          <Route path="/recommended" element={<Recommended />} />
        </Route>
      </Routes>
    </>
  );
}

export default App;
