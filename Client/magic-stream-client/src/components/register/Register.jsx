import { useState, useEffect } from "react";
import Container from "react-bootstrap/Container";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import axiosClient from "../../api/axiosConfig";
import { useNavigate, Link, replace } from "react-router-dom";
import logo from '../../assets/react.svg';

const Register = () => {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [favoriteGenres, setFavoriteGenres] = useState([]);
  const [genres, setGenres] = useState([]);

  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleGenreChange = (e) => {
    const options = Array.from(e.target.selectedOptions);
    setFavoriteGenres(
      options.map((option) => ({
        genre_id: Number(option.value),
        genre_name: option.label,
      }))
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    const defaultRole = "USER";

    console.log(defaultRole);

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    setLoading(true);

    try {
      const payload = {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
        role: defaultRole,
        favorite_genres: favoriteGenres,
      };
      const response = await axiosClient.post("/register", payload);
	  console.log(response.data)
      if (response.data.error) {
        setError(response.data.error);
		console.log(response.data.error)
        return;
      }
      //   Registration is successful, redirect to login
      navigate("/login", {replace: true});
    } catch (err) {
		console.error("Registration failed please try again", err.response?.data || err.message)
	} finally {
		setLoading(false)
	}
  };

  useEffect(() => {
    const fetchGenres = async () => {
      try {
        const response = await axiosClient.get("/genres");
        setGenres(response.data);
      } catch (err) {
        console.error("Error while fetching genres from the server");
      }
    };
    fetchGenres();
  }, []);

  return (
    <Container className="login-container d-flex align-items-center justify-content-center min-vh-100">
      <div
        className="login-card shadow p-4 rounded bg-white"
        style={{ maxWidth: 400, width: "100%" }}
      >
        <div className="text-center mb-4">
          <img src={logo} alt="Logo" width={60} className="mb-2" />
          <h2 className="fw-bold">Register</h2>
          <p className="text-muted">Create your Magic Movie Stream Account</p>
          {error && <div className="alert alert-danger py-2">{error}</div>}
        </div>
        <Form onSubmit={handleSubmit}>
          <Form.Group className="mb-3">
            <Form.Label>First Name</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter First Name"
              value={firstName}
              required
              onChange={(e) => setFirstName(e.target.value)}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Last Name</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter Last Name"
              value={lastName}
              required
              onChange={(e) => setLastName(e.target.value)}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>E-Mail</Form.Label>
            <Form.Control
              type="email"
              placeholder="Enter E-mail"
              value={email}
              required
              onChange={(e) => setEmail(e.target.value)}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Enter password"
              value={password}
              required
              onChange={(e) => setPassword(e.target.value)}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label>Confirm Password</Form.Label>
            <Form.Control
              type="password"
              placeholder="Confirm password"
              value={confirmPassword}
              required
              onChange={(e) => setConfirmPassword(e.target.value)}
              isInvalid={!!confirmPassword && password !== confirmPassword}
            />
            <Form.Control.Feedback type="invalid">
              Passwords dont match
            </Form.Control.Feedback>
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Select
              multiple
              value={favoriteGenres.map((genre) => String(genre.genre_id))}
              required
              onChange={handleGenreChange}
            >
              {genres.map((genre) => (
                <option
                  key={genre.genre_id}
                  value={genre.genre_id}
                  label={genre.genre_name}
                >
                  {genre.genre_name}
                </option>
              ))}
            </Form.Select>
          </Form.Group>
          <Form.Text className="text-muted mb-3">
            Hold Ctrl (Windows) or Cmd (Mac) to select multiple genres
          </Form.Text>
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
                  Registering...
                </span>
              </>
            ) : (
              "Register"
            )}
          </Button>
        </Form>
      </div>
    </Container>
  );
};

export default Register;
