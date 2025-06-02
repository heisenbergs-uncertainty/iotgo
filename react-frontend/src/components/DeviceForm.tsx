import { useState } from 'react'
import type {Device} from '../types/device'
import type { Site } from '../types/site'
import type { ValueStream } from '../types/valueStream'

interface DeviceFormProps {
  device?: Device
  sites: Site[]
  valueStreams: ValueStream[]
  onSubmit: (device: Partial<Device>) => void
  onCancel: () => void
  error?: string
  isLoading?: boolean
}

function DeviceForm({ device = undefined, sites, valueStreams, onSubmit, onCancel, error, isLoading }: DeviceFormProps) {
  const [name, setName] = useState(device?.name || '')
  const [siteId, setSiteId] = useState(device?.siteId?.toString() || '')
  const [valueStreamId, setValueStreamId] = useState(device?.valueStreamId?.toString() || '')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit({
      name,
      siteId: siteId ? parseInt(siteId) : undefined,
      valueStreamId: valueStreamId ? parseInt(valueStreamId) : undefined,
    })
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {error && (
        <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
          <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
        </div>
      )}
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Device Name
        </label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
          placeholder="Enter device name"
          required
        />
      </div>
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Site
        </label>
        <select
          name="siteId"
          title='Select a site'
          value={siteId}
          onChange={(e) => setSiteId(e.target.value)}
          className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
        >
          <option value="">Select a site (optional)</option>
          {sites.map((site) => (
            <option key={site.id} value={site.id}>
              {site.name}
            </option>
          ))}
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
          Value Stream
        </label>
        <select
          value={valueStreamId}
          name='valueStreamId'
          title='Select a value stream'
          onChange={(e) => setValueStreamId(e.target.value)}
          className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
        >
          <option value="">Select a value stream (optional)</option>
          {valueStreams.map((vs) => (
            <option key={vs.id} value={vs.id}>
              {vs.name}
            </option>
          ))}
        </select>
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
          className={`px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 ${isLoading ? 'opacity-80 cursor-not-allowed' : ''}`}
        >
          {isLoading ? 'Saving...' : device?.id ? 'Update Device' : 'Create Device'}
        </button>
      </div>
    </form>
  )
}

export default DeviceForm