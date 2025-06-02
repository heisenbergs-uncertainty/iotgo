import { useEffect } from "react";
import "./App.css";
import { Navigate, Route, Routes } from "react-router-dom";
import { useThemeStore } from "./hooks/useThemeStore";
import Layout from "./components/Layout";
import LandingPage from "./pages/LandingPage";
import LoginPage from "./pages/LoginPage";
import { useAuthStore } from "./hooks/useAuthStore";
import DevicesPage from "./pages/DevicesPage";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import SitesPage from "./pages/SitesPage";
import ValueStreamsPage from "./pages/ValueStreamsPage";
import ProfilePage from "./pages/ProfilePage";
import PlatformsPage from "./pages/PlatformsPage";

const queryClient = new QueryClient();

function App() {
  const { checkAuth, isCheckingAuth } = useAuthStore();
  const { theme } = useThemeStore();

  useEffect(() => {
    checkAuth();
    // Apply dark mode class to html element based on theme
    document.documentElement.classList.toggle("dark", theme === "dark");
  }, []);

  if (isCheckingAuth) {
    return <div>Loading...</div>;
  }

  return (
    <QueryClientProvider client={queryClient}>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<LandingPage />} />
          <Route path="/devices" element={<DevicesPage />} />

          <Route path="/sites" element={<SitesPage />} />
          <Route path="/value-streams" element={<ValueStreamsPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/platforms" element={<PlatformsPage />} />
          <Route
            path="/analytics"
            element={
              <div className="p-8 text-center">Analytics page coming soon</div>
            }
          />
          <Route
            path="/activity"
            element={
              <div className="p-8 text-center">Activity page coming soon</div>
            }
          />
          <Route
            path="/alerts"
            element={
              <div className="p-8 text-center">Alerts page coming soon</div>
            }
          />
          <Route path="/admin/*" element={<Navigate to="/" replace />} />
        </Route>

        <Route path="/login" element={<LoginPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </QueryClientProvider>
  );
}

export default App;
