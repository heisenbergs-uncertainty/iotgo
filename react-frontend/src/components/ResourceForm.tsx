import { useState, useEffect } from "react";
import {
  type Resource,
  type RESTResourceDetails,
  type InfluxDBResourceDetails,
  type PlatformType,
  type ResourceType,
} from "../types/platform";

interface ResourceFormProps {
  resource?: { id: number; name: string; type: ResourceType; details: any };
  platformType: PlatformType;
  onSubmit: (resource: Partial<Resource>) => void;
  onCancel: () => void;
  error?: string;
  isLoading?: boolean;
}

function ResourceForm({
  resource = undefined,
  platformType,
  onSubmit,
  onCancel,
  error,
  isLoading,
}: ResourceFormProps) {
  const [name, setName] = useState(resource?.name || "");
  const [type, setType] = useState<ResourceType>(
    resource?.type || platformType 
  );
  const [restDetails, setRestDetails] = useState<Partial<RESTResourceDetails>>(
    resource?.type === "REST"
      ? resource?.details
      : { method: "GET", headers: {}, queryParams: {} },
  );
  const [influxDetails, setInfluxDetails] = useState<
    Partial<InfluxDBResourceDetails>
  >(resource?.type === "InfluxDB"? resource?.details : {});

  useEffect(() => {
    setType(platformType);
  }, [platformType]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const details =
      type === "REST"
        ? JSON.stringify(restDetails)
        : JSON.stringify(influxDetails);
    onSubmit({ name, type, details });
  };

  const handleAddHeader = () => {
    setRestDetails((prev) => ({
      ...prev,
      headers: { ...prev.headers, "": "" },
    }));
  };

  const handleAddQueryParam = () => {
    setRestDetails((prev) => ({
      ...prev,
      queryParams: { ...prev.queryParams, "": "" },
    }));
  };

  const handleUpdateHeader = (
    oldKey: string,
    newKey: string,
    value: string,
  ) => {
    setRestDetails((prev) => {
      const newHeaders = { ...prev.headers };
      if (newKey !== oldKey) {
        delete newHeaders[oldKey];
      }
      newHeaders[newKey] = value;
      return { ...prev, headers: newHeaders };
    });
  };

  const handleUpdateQueryParam = (
    oldKey: string,
    newKey: string,
    value: string,
  ) => {
    setRestDetails((prev) => {
      const newQueryParams = { ...prev.queryParams };
      if (newKey !== oldKey) {
        delete newQueryParams[oldKey];
      }
      newQueryParams[newKey] = value;
      return { ...prev, queryParams: newQueryParams };
    });
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
          placeholder="Enter resource name"
          required
        />
      </div>
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Type
        </label>
        <select
          title="Select resource type"
          value={type}
          onChange={(e) =>
            setType(e.target.value as PlatformType)
          }
          className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
          disabled
        >
          <option value="rest_endpoint">REST Endpoint</option>
          <option value="influxdb_query">InfluxDB Query</option>
        </select>
      </div>
      {type === "REST" && (
        <>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Method
            </label>
            <select
              title="Select HTTP method"
              value={restDetails.method || "GET"}
              onChange={(e) =>
                setRestDetails((prev) => ({
                  ...prev,
                  method: e.target.value as "GET" | "POST" | "PUT" | "DELETE",
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
            >
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Path
            </label>
            <input
              type="text"
              value={restDetails.path || ""}
              onChange={(e) =>
                setRestDetails((prev) => ({ ...prev, path: e.target.value }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter path (e.g., /data)"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Headers
            </label>
            {Object.entries(restDetails.headers || {}).map(
              ([key, value], index) => (
                <div key={index} className="flex space-x-2 mb-2">
                  <input
                    type="text"
                    value={key}
                    onChange={(e) =>
                      handleUpdateHeader(key, e.target.value, value)
                    }
                    className="block w-1/2 px-4 py-2 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                    placeholder="Key"
                  />
                  <input
                    type="text"
                    value={value}
                    onChange={(e) =>
                      handleUpdateHeader(key, key, e.target.value)
                    }
                    className="block w-1/2 px-4 py-2 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                    placeholder="Value"
                  />
                </div>
              ),
            )}
            <button
              type="button"
              onClick={handleAddHeader}
              className="mt-2 px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600"
            >
              Add Header
            </button>
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Query Parameters
            </label>
            {Object.entries(restDetails.queryParams || {}).map(
              ([key, value], index) => (
                <div key={index} className="flex space-x-2 mb-2">
                  <input
                    type="text"
                    value={key}
                    onChange={(e) =>
                      handleUpdateQueryParam(key, e.target.value, value)
                    }
                    className="block w-1/2 px-4 py-2 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                    placeholder="Key"
                  />
                  <input
                    type="text"
                    value={value}
                    onChange={(e) =>
                      handleUpdateQueryParam(key, key, e.target.value)
                    }
                    className="block w-1/2 px-4 py-2 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                    placeholder="Value"
                  />
                </div>
              ),
            )}
            <button
              type="button"
              onClick={handleAddQueryParam}
              className="mt-2 px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600"
            >
              Add Query Parameter
            </button>
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Body (JSON)
            </label>
            <textarea
              value={restDetails.body || ""}
              onChange={(e) =>
                setRestDetails((prev) => ({ ...prev, body: e.target.value }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter JSON body (optional)"
              rows={4}
            />
          </div>
        </>
      )}
      {type === "InfluxDB" && (
        <>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Bucket
            </label>
            <input
              type="text"
              value={influxDetails.bucket || ""}
              onChange={(e) =>
                setInfluxDetails((prev) => ({
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
              Measurement
            </label>
            <input
              type="text"
              value={influxDetails.measurement || ""}
              onChange={(e) =>
                setInfluxDetails((prev) => ({
                  ...prev,
                  measurement: e.target.value,
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter measurement"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Field
            </label>
            <input
              type="text"
              value={influxDetails.field || ""}
              onChange={(e) =>
                setInfluxDetails((prev) => ({ ...prev, field: e.target.value }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter field"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
              Time Range (e.g., -1h)
            </label>
            <input
              type="text"
              value={influxDetails.timeRange || ""}
              onChange={(e) =>
                setInfluxDetails((prev) => ({
                  ...prev,
                  timeRange: e.target.value,
                }))
              }
              className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
              placeholder="Enter time range (e.g., -1h)"
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
            : resource?.id
              ? "Update Resource"
              : "Create Resource"}
        </button>
      </div>
    </form>
  );
}

export default ResourceForm;
