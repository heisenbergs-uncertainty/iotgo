import { useState } from 'react'
import { useDeviceStore } from '../hooks/useDeviceStore'
import {
  useDevices,
  useSites,
  useValueStreams,
  useCreateDevice,
  useUpdateDevice,
  useDeleteDevice,
  useDevice,
} from '../hooks/useQueries'
import DeviceTable from '../components/DeviceTable'
import DeviceModal from '../components/DeviceModal'
import DeviceDetailsModal from '../components/DeviceDetailsModal'
import type { Device } from '../types/device'

function DevicesPage() {
  const {
    isCreateModalOpen,
    isEditModalOpen,
    isDetailsModalOpen,
    selectedDeviceId,
    openCreateModal,
    closeCreateModal,
    openEditModal,
    closeEditModal,
    openDetailsModal,
    closeDetailsModal,
  } = useDeviceStore()
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [nameFilter, setNameFilter] = useState('')
  // const [sort, setSort] = useState('name')

  const { data: devicesData, isLoading: isDevicesLoading, error: devicesError } = useDevices({
    limit,
    offset: (page - 1) * limit,
    name: nameFilter,
    sort: "name",
  })
  const { data: sites, isLoading: isSitesLoading } = useSites()
  const { data: valueStreams, isLoading: isValueStreamsLoading } = useValueStreams()
  const { data: selectedDevice, isLoading: isDeviceLoading } = useDevice(selectedDeviceId || 0)

  const createDeviceMutation = useCreateDevice()
  const updateDeviceMutation = useUpdateDevice()
  const deleteDeviceMutation = useDeleteDevice()

  const handleCreate = (device: Partial<Device>) => {
    createDeviceMutation.mutate(device, {
      onSuccess: () => closeCreateModal(),
      onError: (err: any) => console.error('Create error:', err),
    })
  }

  const handleUpdate = (device: Partial<Device>) => {
    if (selectedDeviceId) {
      updateDeviceMutation.mutate(
        { id: selectedDeviceId, device },
        {
          onSuccess: () => closeEditModal(),
          onError: (err: any) => console.error('Update error:', err),
        },
      )
    }
  }

  const handleDelete = (id: number) => {
    if (window.confirm('Are you sure you want to delete this device?')) {
      deleteDeviceMutation.mutate(id, {
        onError: (err: any) => console.error('Delete error:', err),
      })
    }
  }

  // const handleSort = (newSort: string) => {
  //   setSort(newSort)
  // }

  const totalPages = devicesData ? Math.ceil(devicesData.total / limit) : 1
  const isLoading = isDevicesLoading || isSitesLoading || isValueStreamsLoading || isDeviceLoading

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-secondary-900 dark:text-secondary-50">
          Device Management
        </h2>
        <button
          onClick={openCreateModal}
          className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
        >
          Add Device
        </button>
      </div>

      {devicesError && (
        <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
          <p className="text-sm text-red-600 dark:text-red-400">
            {devicesError.message || 'Failed to load devices'}
          </p>
        </div>
      )}

      <div className="flex items-center space-x-4">
        <input
          type="text"
          value={nameFilter}
          onChange={(e) => setNameFilter(e.target.value)}
          placeholder="Filter by name..."
          className="px-4 py-2 border border-secondary-200 dark:border-secondary-700 rounded-lg dark:bg-secondary-800 dark:text-secondary-100"
        />
      </div>

      {isLoading ? (
        <div>Loading...</div>
      ) : (
        <>
          <DeviceTable
            devices={devicesData?.items || []}
            onEdit={openEditModal}
            onView={openDetailsModal}
            onDelete={handleDelete}
          />
          <div className="flex justify-between items-center mt-4">
            <button
              onClick={() => setPage((p) => Math.max(p - 1, 1))}
              disabled={page === 1}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg disabled:opacity-50"
            >
              Previous
            </button>
            <span>
              Page {page} of {totalPages}
            </span>
            <button
              onClick={() => setPage((p) => Math.min(p + 1, totalPages))}
              disabled={page === totalPages}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </>
      )}

      <DeviceModal
        isOpen={isCreateModalOpen}
        onClose={closeCreateModal}
        sites={sites?.items || []}
        valueStreams={valueStreams?.items || []}
        onSubmit={handleCreate}
        error={createDeviceMutation.error?.message}
        isLoading={createDeviceMutation.isPending}
      />

      <DeviceModal
        isOpen={isEditModalOpen}
        onClose={closeEditModal}
        device={selectedDevice}
        sites={sites?.items || []}
        valueStreams={valueStreams?.items || []}
        onSubmit={handleUpdate}
        error={updateDeviceMutation.error?.message}
        isLoading={updateDeviceMutation.isPending}
      />

      <DeviceDetailsModal
        isOpen={isDetailsModalOpen}
        onClose={closeDetailsModal}
        device={selectedDevice || null}
      />
    </div>
  )
}

export default DevicesPage