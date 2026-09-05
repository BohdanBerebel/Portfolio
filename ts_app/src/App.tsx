import { useState } from "react";

import AuthForm from "./components/AuthForm";
import AuthenticatedLayout from "./components/layout/AuthenticatedLayout";
import { getToken, removeToken } from "./api/authStorage";

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(
    () => getToken() !== null,
  );

  function handleAuthenticated() {
    setIsAuthenticated(true);
  }

  function handleLogout() {
    removeToken();
    setIsAuthenticated(false);
  }

  if (!isAuthenticated) {
    return (
      <main className="auth-page">
        <section className="auth-card">
          <h1>Notes App</h1>
          <p className="subtitle">
            Manage your notes and conversations securely.
          </p>

          <AuthForm onAuthenticated={handleAuthenticated} />
        </section>
      </main>
    );
  }

  return <AuthenticatedLayout onLogout={handleLogout} />;
}

export default App;
