import { useState } from "react";
import type { GenerateApiKeyRequest } from "../types/apikey";

interface ApiKeyFormProps {
  onSubmit: (request: GenerateApiKeyRequest) => void;
  onCancel: () => void;
  error?: string;
  isLoading?: boolean;
}

function ApiKeyForm({ onSubmit, onCancel, error, isLoading }: ApiKeyFormProps) {
  const [name, setName] = useState("");
  const [scopes, setScopes] = useState<string[]>([]);

  const handleScopeChange = (scope: string) => {
    setScopes((prev) =>
      prev.includes(scope) ? prev.filter((s) => s !== scope) : [...prev, scope],
    );
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ name, scopes });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {error && (
        <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
          <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
        </div>
      )}
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Name
        </label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
          placeholder="Enter API key name"
          required
        />
      </div>
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Scopes
        </label>
        <div className="space-y-2">
          <label className="flex items-center">
            <input
              type="checkbox"
              checked={scopes.includes("read")}
              onChange={() => handleScopeChange("read")}
              className="h-4 w-4 rounded border-secondary-300 dark:border-secondary-700 text-primary-600 focus:ring-primary-500"
            />
            <span className="ml-2 text-sm text-secondary-600 dark:text-secondary-400">
              Read
            </span>
          </label>
          <label className="flex items-center">
            <input
              type="checkbox"
              checked={scopes.includes("write")}
              onChange={() => handleScopeChange("write")}
              className="h-4 w-4 rounded border-secondary-300 dark:border-secondary-700 text-primary-600 focus:ring-primary-500"
            />
            <span className="ml-2 text-sm text-secondary-600 dark:text-secondary-400">
              Write
            </span>
          </label>
        </div>
      </div>
      <div className="flex justify-end space-x-4">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isLoading}
          className={`px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 ${isLoading ? "opacity-80 cursor-not-allowed" : ""}`}
        >
          {isLoading ? "Generating..." : "Generate API Key"}
        </button>
      </div>
    </form>
  );
}

export default ApiKeyForm;
