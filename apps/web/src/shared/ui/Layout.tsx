import { Link, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../../features/auth/AuthProvider";

export function Layout() {
  const auth = useAuth();
  const navigate = useNavigate();

  async function onLogout(): Promise<void> {
    await auth.logout();
    navigate("/login");
  }

  return (
    <div className="container">
      <header className="nav">
        <Link to="/" className="brand">
          Boilerplate
        </Link>
        <div className="nav-right">
          <span className="muted">{auth.user?.username}</span>
          <button type="button" className="btn btn-secondary" onClick={onLogout}>
            Logout
          </button>
        </div>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  );
}
