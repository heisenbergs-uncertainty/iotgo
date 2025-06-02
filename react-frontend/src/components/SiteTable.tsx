import { useState } from 'react'
import type { Site } from '../types/site'

interface SiteTableProps {
  sites: Site[]
  onEdit: (id: number) => void
  onView: (id: number) => void
  onDelete: (id: number) => void
}

function SiteTable({ sites, onEdit, onView, onDelete }: SiteTableProps) {
  const [sort, setSort] = useState('name')

  const handleSort = (newSort: string) => {
    setSort(newSort)
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white dark:bg-gray-800 shadow-md rounded-lg">
        <thead>
          <tr className="bg-gray-100 dark:bg-gray-700">
            <th
              className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100 cursor-pointer"
              onClick={() => handleSort(sort === 'name' ? '-name' : 'name')}
            >
              Name {sort === 'name' ? '↑' : sort === '-name' ? '↓' : ''}
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Address
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              City
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              State
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Country
            </th>
            <th className="px-6 py-3 text-right text-sm font-medium text-gray-900 dark:text-gray-100">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          {sites.map((site) => (
            <tr key={site.id} className="border-b dark:border-gray-700">
              <td
                className="px-6 py-4 text-sm text-primary-600 hover:text-primary-500 cursor-pointer"
                onClick={() => onView(site.id)}
              >
                {site.name}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {site.address}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {site.city}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {site.state}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {site.country}
              </td>
              <td className="px-6 py-4 text-right text-sm">
                <button
                  onClick={() => onEdit(site.id)}
                  className="text-indigo-600 hover:text-indigo-500 mr-4"
                >
                  Edit
                </button>
                <button
                  onClick={() => onDelete(site.id)}
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
  )
}

export default SiteTable