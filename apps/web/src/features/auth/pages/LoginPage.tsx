import { useState } from "react";
import type { FormEvent } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { ApiError } from "../../../shared/lib/api";
import { useAuth } from "../AuthProvider";
import { loginSchema } from "../api";

export function LoginPage() {
  const auth = useAuth();
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  if (auth.loading) {
    return <p className="muted">Loading…</p>;
  }
  if (auth.user) {
    return <Navigate to="/" replace />;
  }

  async function onSubmit(e: FormEvent): Promise<void> {
    e.preventDefault();
    const parsed = loginSchema.safeParse({ username, password });
    if (!parsed.success) {
      setError(parsed.error.issues[0]?.message ?? "invalid input");
      return;
    }
    setError(null);
    setBusy(true);
    try {
      await auth.login(parsed.data);
      navigate("/");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "login failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="container narrow">
      <h1>Sign in</h1>
      <form className="card" onSubmit={onSubmit}>
        <label className="field" htmlFor="username">
          Username
          <input
            id="username"
            className="input"
            autoComplete="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </label>
        <label className="field" htmlFor="password">
          Password
          <input
            id="password"
            className="input"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>
        {error && (
          <p className="error" role="alert">
            {error}
          </p>
        )}
        <button type="submit" className="btn" disabled={busy}>
          {busy ? "Signing in…" : "Sign in"}
        </button>
        <p className="muted">Default dev account: admin / admin123</p>
      </form>
    </div>
  );
}
