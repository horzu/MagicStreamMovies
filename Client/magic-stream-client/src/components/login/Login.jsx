import { useState } from "react";
import Container from "react-bootstrap/Container";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import axiosClient from "../../api/axiosConfig";
import { useNavigate, useLocation, Link } from "react-router-dom";
import logo from "../../assets/react.svg";
import useAuth from "../../hook/useAuth";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { setAuth } = useAuth();

  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await axiosClient.post("/login", { email, password });
      console.log(response.data);
      if (response.data.error) {
        setError(response.data.error);
        return;
      }
      setAuth(response.data);
      localStorage.setItem("user", JSON.stringify(response.data));
      //   navigate(from, {replace: true})
      navigate("/");
    } catch (err) {
      console.error(err);
      setError("Invalid email or password");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container className="login-container d-flex align-items-center justify-content-center min-vh-100">
      <div
        className="login-card shadow p-4 rounded bg-white"
        style={{ maxWidth: 400, width: "100%" }}
      >
        <div className="text-center mb-4">
          <img src={logo} alt="Logo" width={60} className="mb-2" />
          <h2 className="fw-bold">Login</h2>
          <p className="text-muted">Login your Magic Movie Stream Account</p>
          {error && <div className="alert alert-danger py-2">{error}</div>}
        </div>

        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3" controlId="formBasicEmail">
            <Form.Label>E-Mail</Form.Label>
            <Form.Control
              type="email"
              placeholder="Enter E-mail"
              value={email}
              required
              onChange={(e) => setEmail(e.target.value)}
            />
          </Form.Group>
          <Form.Group className="mb-3" controlId="formBasicPassword">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Enter password"
              value={password}
              required
              autoFocus
              onChange={(e) => setPassword(e.target.value)}
            />
          </Form.Group>
          <Button
            variant="primary"
            type="submit"
            className="w-100 my-2"
            disabled={loading}
            style={{ fontWeight: 600, letterSpacing: 1 }}
          >
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                >
                  Logging in...
                </span>
              </>
            ) : (
              "Login"
            )}
          </Button>
        </Form>
        <div className="text-center mt-3">
          <span className="text-muted">Don't have account?</span>
          <Link to="/register" className="fw-semibold">
            Register here
          </Link>
        </div>
      </div>
    </Container>
  );
};

export default Login;
