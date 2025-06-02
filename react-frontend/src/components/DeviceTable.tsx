import { useState } from 'react'
import type { Device } from '../types/device'

interface DeviceTableProps {
  devices: Device[]
  onEdit: (id: number) => void
  onView: (id: number) => void
  onDelete: (id: number) => void
}

function DeviceTable({ devices, onEdit, onView, onDelete }: DeviceTableProps) {
  const [sort, setSort] = useState('name')

  const handleSort = (newSort: string) => {
    setSort(newSort)
    // Sorting is handled by parent via query params
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
              Site
            </th>
            <th className="px-6 py-3 text-left text-sm font-medium text-gray-900 dark:text-gray-100">
              Value Stream
            </th>
            <th className="px-6 py-3 text-right text-sm font-medium text-gray-900 dark:text-gray-100">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          {devices.map((device) => (
            <tr key={device.id} className="border-b dark:border-gray-700">
              <td
                className="px-6 py-4 text-sm text-primary-600 hover:text-primary-500 cursor-pointer"
                onClick={() => onView(device.id)}
              >
                {device.name}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {device.site?.name || 'N/A'}
              </td>
              <td className="px-6 py-4 text-sm text-gray-900 dark:text-gray-100">
                {device.valueStream?.name || 'N/A'}
              </td>
              <td className="px-6 py-4 text-right text-sm">
                <button
                  onClick={() => onEdit(device.id)}
                  className="text-indigo-600 hover:text-indigo-500 mr-4"
                >
                  Edit
                </button>
                <button
                  onClick={() => onDelete(device.id)}
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

export default DeviceTable