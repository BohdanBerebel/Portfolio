import { useState } from "react";
import type { SubmitEvent } from "react";
import { login, register } from "../api/auth";
import { setToken, setUserId } from "../api/authStorage";

type AuthMode = "login" | "register";

interface AuthFormProps {
  onAuthenticated: () => void;
}

function AuthForm({ onAuthenticated }: AuthFormProps) {
  const [mode, setMode] = useState<AuthMode>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    setMessage("");
    setLoading(true);

    try {
      const result =
        mode === "login"
          ? await login({ email, password })
          : await register({ email, password });

      setToken(result.token);
      setUserId(result.id);

      onAuthenticated();
    } catch {
      setMessage("Authentication failed.");
    } finally {
      setLoading(false);
    }
  }

  function switchMode() {
    setMode(mode === "login" ? "register" : "login");
    setMessage("");
  }

  return (
    <div>
      <h2>{mode === "login" ? "Login" : "Register"}</h2>

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="email">Email</label>

          <input
            id="email"
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="password">Password</label>

          <input
            id="password"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Loading..." : mode === "login" ? "Login" : "Register"}
        </button>
      </form>

      {message && <p className="error-message">{message}</p>}

      <button className="auth-switch" type="button" onClick={switchMode}>
        {mode === "login"
          ? "Create an account"
          : "Already have an account? Login"}
      </button>
    </div>
  );
}

export default AuthForm;
