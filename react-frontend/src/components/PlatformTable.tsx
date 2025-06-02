import { useState } from "react";
import type { Platform } from "../types/platform";
import { usePlatformStore } from "../hooks/usePlatformStore";

interface PlatformTableProps {
  platforms: Platform[];
  onEdit: (id: number) => void;
  onDelete: (id: number) => void;
  onSelect: (id: number) => void;
}

function PlatformTable({
  platforms,
  onEdit,
  onDelete,
  onSelect,
}: PlatformTableProps) {
  const [sort, setSort] = useState("name");
  const { setSelectedPlatformId } = usePlatformStore();

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
          {platforms.map((platform) => (
            <tr
              key={platform.id}
              className="border-b dark:border-gray-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700"
              onClick={() => {
                setSelectedPlatformId(platform.id);
                onSelect(platform.id);
              }}
            >
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {platform.name}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {platform.type}
              </td>
              <td className="px-6 py-4 text-right text-sm">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onEdit(platform.id);
                  }}
                  className="text-indigo-600 hover:text-indigo-500 mr-4"
                >
                  Edit
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onDelete(platform.id);
                  }}
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

export default PlatformTable;
