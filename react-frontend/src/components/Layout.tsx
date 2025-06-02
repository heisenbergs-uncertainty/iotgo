import { useEffect } from "react";
import { Outlet, Navigate, useLocation, Link } from "react-router-dom";
import { useAuthStore } from "../hooks/useAuthStore";
import { useThemeStore } from "../hooks/useThemeStore";

function Layout() {
  const { user, isCheckingAuth } = useAuthStore();
  const { theme, toggleTheme } = useThemeStore();
  const location = useLocation();

  useEffect(() => {
    useAuthStore.getState().checkAuth();
  }, []);

  if (isCheckingAuth) {
    return <div>Loading...</div>;
  }

  if (!user && location.pathname !== "/login") {
    return <Navigate to="/login" replace />;
  }

  if (
    user &&
    location.pathname.startsWith("/admin") &&
    user.role !== "superuser"
  ) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      <header className="bg-white dark:bg-gray-800 shadow">
        <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <div className="flex items-center space-x-4">
            <h1 className="text-xl font-bold">
              <Link to="/" className="text-gray-900 dark:text-gray-100">
                IoTGo
              </Link>
            </h1>
            <Link
              to="/devices"
              className="text-gray-600 dark:text-gray-300 hover:text-primary-600"
            >
              Devices
            </Link>
            <Link
              to="/sites"
              className="text-gray-600 dark:text-gray-300 hover:text-primary-600"
            >
              Sites
            </Link>
            <Link
              to="/value-streams"
              className="text-gray-600 dark:text-gray-300 hover:text-primary-600"
            >
              Value Streams
            </Link>
            <Link
              to="/platforms"
              className="text-gray-600 dark:text-gray-300 hover:text-primary-600"
            >
              Platforms
            </Link>
            <Link
              to="/profile"
              className="text-gray-600 dark:text-gray-300 hover:text-primary-600"
            >
              Profile
            </Link>
          </div>
          <div className="flex items-center space-x-4">
            {user ? (
              <>
                <span>{user.role}</span>
                <button
                  onClick={() => useAuthStore.getState().logout()}
                  className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                >
                  Logout
                </button>
              </>
            ) : null}
            <button
              onClick={toggleTheme}
              className="p-2 rounded-full bg-gray-200 dark:bg-gray-700"
            >
              {theme === "light" ? "🌙" : "☀️"}
            </button>
          </div>
        </nav>
      </header>
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>
    </div>
  );
}

export default Layout;
