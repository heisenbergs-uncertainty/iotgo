import { create } from 'zustand'

interface ValueStreamState {
  selectedValueStreamId: number | null
  isCreateModalOpen: boolean
  isEditModalOpen: boolean
  isDetailsModalOpen: boolean
  setSelectedValueStreamId: (id: number | null) => void
  openCreateModal: () => void
  closeCreateModal: () => void
  openEditModal: (id: number) => void
  closeEditModal: () => void
  openDetailsModal: (id: number) => void
  closeDetailsModal: () => void
}

export const useValueStreamStore = create<ValueStreamState>((set) => ({
  selectedValueStreamId: null,
  isCreateModalOpen: false,
  isEditModalOpen: false,
  isDetailsModalOpen: false,
  setSelectedValueStreamId: (id) => set({ selectedValueStreamId: id }),
  openCreateModal: () => set({ isCreateModalOpen: true }),
  closeCreateModal: () => set({ isCreateModalOpen: false }),
  openEditModal: (id) => set({ selectedValueStreamId: id, isEditModalOpen: true }),
  closeEditModal: () => set({ selectedValueStreamId: null, isEditModalOpen: false }),
  openDetailsModal: (id) => set({ selectedValueStreamId: id, isDetailsModalOpen: true }),
  closeDetailsModal: () => set({ selectedValueStreamId: null, isDetailsModalOpen: false }),
}))