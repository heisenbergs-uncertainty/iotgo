import { useState } from "react";
import type { Resource } from "../types/platform";

interface ResourceTableProps {
  resources: Resource[];
  onEdit: (id: number) => void;
  onDelete: (id: number) => void;
  onTest: (id: number) => void;
}

function ResourceTable({
  resources,
  onEdit,
  onDelete,
  onTest,
}: ResourceTableProps) {
  const [sort, setSort] = useState("name");

  const handleSort = (newSort: string) => {
    setSort(newSort);
  };

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white dark:bg-gray-800 shadow-md rounded-lg">
        <thead>
          <tr className="bg-gray-100 dark:bg-gray-700">
            <th
              className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100 cursor-pointer"
              onClick={() => handleSort(sort === "name" ? "-name" : "name")}
            >
              Name {sort === "name" ? "↑" : sort === "-name" ? "↓" : ""}
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Type
            </th>
            <th className="px-6 py-3 text-right text-sm font-medium text-gray-900 dark:text-gray-100">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          {resources.map((resource) => (
            <tr key={resource.id} className="border-b dark:border-gray-700">
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {resource.name}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {resource.type}
              </td>
              <td className="px-6 py-4 text-right text-sm">
                <button
                  onClick={() => onEdit(resource.id)}
                  className="text-indigo-600 hover:text-indigo-500 mr-4"
                >
                  Edit
                </button>
                <button
                  onClick={() => onTest(resource.id)}
                  className="text-green-600 hover:text-green-500 mr-4"
                >
                  Test
                </button>
                <button
                  onClick={() => onDelete(resource.id)}
                  className="text-red-600 hover:text-red-500"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default ResourceTable;
