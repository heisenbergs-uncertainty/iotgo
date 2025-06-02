import { useState } from "react";
import type { ValueStream } from "../types/valueStream";

interface ValueStreamFormProps {
  valueStream?: ValueStream;
  onSubmit: (valueStream: Partial<ValueStream>) => void;
  onCancel: () => void;
  error?: string;
  isLoading?: boolean;
}

function ValueStreamForm({
  valueStream = undefined,
  onSubmit,
  onCancel,
  error,
  isLoading,
}: ValueStreamFormProps) {
  const [name, setName] = useState(valueStream?.name || "");
  const [type, setType] = useState(valueStream?.type || "");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      name,
      type,
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
          placeholder="Enter value stream name"
          required
        />
      </div>
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Type
        </label>
        <input
          type="text"
          value={type}
          onChange={(e) => setType(e.target.value)}
          className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
          placeholder="Enter value stream type"
          required
        />
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
          {isLoading
            ? "Saving..."
            : valueStream?.id
              ? "Update Value Stream"
              : "Create Value Stream"}
        </button>
      </div>
    </form>
  );
}

export default ValueStreamForm;
