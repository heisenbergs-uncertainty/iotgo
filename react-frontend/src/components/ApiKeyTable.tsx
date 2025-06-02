import type { ApiKey } from "../types/apikey";

interface ApiKeyTableProps {
  apiKeys: ApiKey[];
  onRevoke: (keyId: string) => void;
}

function ApiKeyTable({ apiKeys, onRevoke }: ApiKeyTableProps) {
  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white dark:bg-gray-800 shadow-md rounded-lg">
        <thead>
          <tr className="bg-gray-100 dark:bg-gray-700">
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Name
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Key ID
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Scopes
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Status
            </th>
            <th className="px-6 py-3 text-right text-sm font-medium text-gray-900 dark:text-gray-100">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          {apiKeys.map((key) => {
            const metadata = JSON.parse(key.metadata || "{}");
            const scopes = metadata.scopes?.join(", ") || "N/A";
            return (
              <tr key={key.keyId} className="border-b dark:border-gray-700">
                <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                  {key.name}
                </td>
                <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                  {key.keyId}
                </td>
                <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                  {scopes}
                </td>
                <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                  {key.isActive ? "Active" : "Revoked"}
                </td>
                <td className="px-6 py-4 text-right text-sm">
                  {key.isActive && (
                    <button
                      onClick={() => onRevoke(key.keyId)}
                      className="text-red-600 hover:text-red-500"
                    >
                      Revoke
                    </button>
                  )}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export default ApiKeyTable;
