import "./App.css";
import Home from "./components/home/Home";
import Header from "./components/header/Header";
import Register from "./components/register/Register";
import Login from "./components/login/Login";
import { Route, Routes, useNavigate } from "react-router-dom";

function App() {
  return (
    <>
      <Header />
      <Routes>
        <Route element={<Home />} path="/"></Route>
        <Route element={<Register />} path="/register"></Route>
        <Route element={<Login />} path="/login"></Route>
      </Routes>
    </>
  );
}

export default App;
