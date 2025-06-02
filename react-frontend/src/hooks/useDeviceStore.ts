import { create } from 'zustand'

interface DeviceState {
  selectedDeviceId: number | null
  isCreateModalOpen: boolean
  isEditModalOpen: boolean
  isDetailsModalOpen: boolean
  setSelectedDeviceId: (id: number | null) => void
  openCreateModal: () => void
  closeCreateModal: () => void
  openEditModal: (id: number) => void
  closeEditModal: () => void
  openDetailsModal: (id: number) => void
  closeDetailsModal: () => void
}

export const useDeviceStore = create<DeviceState>((set) => ({
  selectedDeviceId: null,
  isCreateModalOpen: false,
  isEditModalOpen: false,
  isDetailsModalOpen: false,
  setSelectedDeviceId: (id) => set({ selectedDeviceId: id }),
  openCreateModal: () => set({ isCreateModalOpen: true }),
  closeCreateModal: () => set({ isCreateModalOpen: false }),
  openEditModal: (id) => set({ selectedDeviceId: id, isEditModalOpen: true }),
  closeEditModal: () => set({ selectedDeviceId: null, isEditModalOpen: false }),
  openDetailsModal: (id) => set({ selectedDeviceId: id, isDetailsModalOpen: true }),
  closeDetailsModal: () => set({ selectedDeviceId: null, isDetailsModalOpen: false }),
}))