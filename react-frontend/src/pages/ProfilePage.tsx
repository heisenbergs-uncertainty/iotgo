import { useAuthStore } from "../hooks/useAuthStore";
import { useProfileStore } from "../hooks/useProfileStore";
import {
  useUser,
  useApiKeys,
  useGenerateApiKey,
  useRevokeApiKey,
} from "../hooks/useQueries";
import ApiKeyTable from "../components/ApiKeyTable";
import ApiKeyModal from "../components/ApiKeyModal";
import type { GenerateApiKeyRequest } from "../types/apikey";

function ProfilePage() {
  const { user } = useAuthStore();
  const {
    selectedTab,
    setSelectedTab,
    isApiKeyModalOpen,
    openApiKeyModal,
    closeApiKeyModal,
  } = useProfileStore();

  // Conditionally call useUser and useApiKeys only if user.id is defined
  const {
    data: userData,
    isLoading: isUserLoading,
    error: userError,
  } = useUser(user?.id as number); // Assuming useUser expects a number, cast it if user.id is present
  const {
    data: apiKeys,
    isLoading: isApiKeysLoading,
    error: apiKeysError,
  } = useApiKeys(user?.id as string); // Assuming useApiKeys expects a string, cast it if user.id is present

  const generateApiKeyMutation = useGenerateApiKey();
  const revokeApiKeyMutation = useRevokeApiKey();

  const handleGenerateApiKey = (request: GenerateApiKeyRequest) => {
    generateApiKeyMutation.mutate(request, {
      onSuccess: () => closeApiKeyModal(),
      onError: (err: any) => console.error("Generate error:", err),
    });
  };

  const handleRevokeApiKey = (keyId: string) => {
    if (window.confirm("Are you sure you want to revoke this API key?")) {
      revokeApiKeyMutation.mutate(keyId, {
        onError: (err: any) => console.error("Revoke error:", err),
      });
    }
  };

  const tabs = [
    { id: "general", label: "General" },
    { id: "developer", label: "Developer" },
  ];

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-secondary-900 dark:text-secondary-50">
        User Profile
      </h2>
      <div className="flex flex-col md:flex-row gap-6">
        {/* Sidebar */}
        <div className="w-full md:w-1/4">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md">
            <nav className="p-4">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setSelectedTab(tab.id)}
                  className={`block w-full text-left px-4 py-2 text-sm rounded-lg ${
                    selectedTab === tab.id
                      ? "bg-primary-600 text-white"
                      : "text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </nav>
          </div>
        </div>

        {/* Main Content */}
        <div className="w-full md:w-3/4">
          {isUserLoading ? (
            <div>Loading...</div>
          ) : userError ? (
            <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
              <p className="text-sm text-red-600 dark:text-red-400">
                {userError.message || "Failed to load user data"}
              </p>
            </div>
          ) : (
            <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
              {selectedTab === "general" && (
                <div className="space-y-4">
                  <h3 className="text-lg font-medium text-secondary-900 dark:text-secondary-50">
                    General Information
                  </h3>
                  <dl className="space-y-4">
                    <div>
                      <dt className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                        Email
                      </dt>
                      <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
                        {userData?.email || "N/A"}
                      </dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                        Role
                      </dt>
                      <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
                        {userData?.role || "N/A"}
                      </dd>
                    </div>
                  </dl>
                </div>
              )}
              {selectedTab === "developer" && (
                <div className="space-y-4">
                  <div className="flex justify-between items-center">
                    <h3 className="text-lg font-medium text-secondary-900 dark:text-secondary-50">
                      API Keys
                    </h3>
                    <button
                      onClick={openApiKeyModal}
                      className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
                    >
                      Generate API Key
                    </button>
                  </div>
                  {isApiKeysLoading ? (
                    <div>Loading...</div>
                  ) : apiKeysError ? (
                    <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
                      <p className="text-sm text-red-600 dark:text-red-400">
                        {apiKeysError.message || "Failed to load API keys"}
                      </p>
                    </div>
                  ) : (
                    <ApiKeyTable
                      apiKeys={apiKeys || []}
                      onRevoke={handleRevokeApiKey}
                    />
                  )}
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <ApiKeyModal
        isOpen={isApiKeyModalOpen}
        onClose={closeApiKeyModal}
        onSubmit={handleGenerateApiKey}
        error={generateApiKeyMutation.error?.message}
        isLoading={generateApiKeyMutation.isPending}
      />
    </div>
  );
}

export default ProfilePage;
