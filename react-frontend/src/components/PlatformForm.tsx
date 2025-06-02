import { useState } from "react";
import type {
  Platform,
  RESTMetadata,
  InfluxDBMetadata,
  ApiAuthMethodType,
} from "../types/platform";

interface PlatformFormProps {
  platform?: Platform;
  onSubmit: (platform: Partial<Platform>) => void;
  onCancel: () => void;
  error?: string;
  isLoading?: boolean;
}

function PlatformForm({
  platform = undefined,
  onSubmit,
  onCancel,
  error,
  isLoading,
}: PlatformFormProps) {
  const [name, setName] = useState(platform?.name || "");
  const [type, setType] = useState(platform?.type || "REST");
  const [restMetadata, setRestMetadata] = useState<Partial<RESTMetadata>>(
    platform?.type === "REST"
      ? JSON.parse(platform?.metadata || "{}")
      : { auth: { type: "none" } },
  );
  const [influxMetadata, setInfluxMetadata] = useState<
    Partial<InfluxDBMetadata>
  >(
    platform?.type === "InfluxDB"
      ? JSON.parse(platform?.metadata || "{}")
      : { timeout: 10 },
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const metadata =
      type === "REST"
        ? JSON.stringify(restMetadata)
        : JSON.stringify(influxMetadata);
    onSubmit({ name, type, metadata });
  };

  const handleRestAuthChange = (authType: ApiAuthMethodType) => {
    setRestMetadata((prev) => ({
      ...prev,
      auth: {
        type: authType,
        apiKey: "",
        bearerToken: "",
        basicAuth: undefined,
      },
    }));
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
          placeholder="Enter platform name"
          required
        />
      </div>
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Type
        </label>
        <select
          name="platformType"
          title="Select a platform type"
          value={type}
          onChange={(e) => setType(e.target.value as "REST" | "InfluxDB")}
          className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
        >
          <option value="REST">REST</option>
          <option value="InfluxDB">InfluxDB</option>
        </select>
      </div>
      {type === "REST" && (
        <>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Base Endpoint
            </label>
            <input
              type="text"
              value={restMetadata.baseEndpoint || ""}
              onChange={(e) =>
                setRestMetadata((prev) => ({
                  ...prev,
                  baseEndpoint: e.target.value,
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter base endpoint (e.g., https://api.example.com)"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Timeout (seconds)
            </label>
            <input
              title="Enter timeout in seconds"
              type="number"
              value={restMetadata.timeout || 10}
              onChange={(e) =>
                setRestMetadata((prev) => ({
                  ...prev,
                  timeout: parseInt(e.target.value),
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              min="1"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Authentication Type
            </label>
            <select
              title="Select authentication type"
              value={restMetadata.auth?.type || "none"}
              onChange={(e) => handleRestAuthChange(e.target.value as ApiAuthMethodType)}
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
            >
              <option value="NONE">None</option>
              <option value="API_KEY">API Key</option>
              <option value="BEARER">Bearer Token</option>
              <option value="BASIC">Basic Auth</option>
            </select>
          </div>
          {restMetadata.auth?.type === "API_KEY" && (
            <div>
              <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
                API Key
              </label>
              <input
                type="text"
                value={restMetadata.auth.apiKey || ""}
                onChange={(e) =>
                  setRestMetadata((prev) => ({
                    ...prev,
                    auth: { ...prev.auth!, apiKey: e.target.value },
                  }))
                }
                className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                placeholder="Enter API key"
                required
              />
            </div>
          )}
          {restMetadata.auth?.type === "BEARER" && (
            <div>
              <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
                Bearer Token
              </label>
              <input
                type="text"
                value={restMetadata.auth.bearerToken || ""}
                onChange={(e) =>
                  setRestMetadata((prev) => ({
                    ...prev,
                    auth: { ...prev.auth!, bearerToken: e.target.value },
                  }))
                }
                className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                placeholder="Enter bearer token"
                required
              />
            </div>
          )}
          {restMetadata.auth?.type === "BASIC" && (

            <>
              <div>
                <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
                  Username
                </label>
                <input
                  type="text"
                  value={restMetadata.auth.basicAuth?.username || ""}
                  onChange={(e) =>
                    setRestMetadata((prev) => ({
                      ...prev,
                      auth: {
                        ...prev.auth!,
                        basicAuth: {
                          ...prev.auth!.basicAuth!,
                          username: e.target.value,
                        },
                      },
                    }))
                  }
                  className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                  placeholder="Enter username"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
                  Password
                </label>
                <input
                  type="password"
                  value={restMetadata.auth.basicAuth?.password || ""}
                  onChange={(e) =>
                    setRestMetadata((prev) => ({
                      ...prev,
                      auth: {
                        ...prev.auth!,
                        basicAuth: {
                          ...prev.auth!.basicAuth!,
                          password: e.target.value,
                        },
                      },
                    }))
                  }
                  className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                  placeholder="Enter password"
                  required
                />
              </div>
            </>
          )}
        </>
      )}
      {type === "InfluxDB" && (
        <>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              URL
            </label>
            <input
              type="text"
              value={influxMetadata.url || ""}
              onChange={(e) =>
                setInfluxMetadata((prev) => ({ ...prev, url: e.target.value }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter InfluxDB URL (e.g., https://influxdb.example.com)"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Token
            </label>
            <input
              type="text"
              value={influxMetadata.token || ""}
              onChange={(e) =>
                setInfluxMetadata((prev) => ({
                  ...prev,
                  token: e.target.value,
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter InfluxDB token"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Organization
            </label>
            <input
              type="text"
              value={influxMetadata.org || ""}
              onChange={(e) =>
                setInfluxMetadata((prev) => ({ ...prev, org: e.target.value }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter organization"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Bucket
            </label>
            <input
              type="text"
              value={influxMetadata.bucket || ""}
              onChange={(e) =>
                setInfluxMetadata((prev) => ({
                  ...prev,
                  bucket: e.target.value,
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter bucket"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Timeout (seconds)
            </label>
            <input
              title="Timeout (seconds)"
              type="number"
              value={influxMetadata.timeout || 10}
              onChange={(e) =>
                setInfluxMetadata((prev) => ({
                  ...prev,
                  timeout: parseInt(e.target.value),
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              min="1"
              required
            />
          </div>
        </>
      )}
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
          {isLoading
            ? "Saving..."
            : platform?.id
              ? "Update Platform"
              : "Create Platform"}
        </button>
      </div>
    </form>
  );
}

export default PlatformForm;
