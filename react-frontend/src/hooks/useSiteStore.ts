import { create } from 'zustand'

interface SiteState {
  selectedSiteId: number | null
  isCreateModalOpen: boolean
  isEditModalOpen: boolean
  isDetailsModalOpen: boolean
  setSelectedSiteId: (id: number | null) => void
  openCreateModal: () => void
  closeCreateModal: () => void
  openEditModal: (id: number) => void
  closeEditModal: () => void
  openDetailsModal: (id: number) => void
  closeDetailsModal: () => void
}

export const useSiteStore = create<SiteState>((set) => ({
  selectedSiteId: null,
  isCreateModalOpen: false,
  isEditModalOpen: false,
  isDetailsModalOpen: false,
  setSelectedSiteId: (id) => set({ selectedSiteId: id }),
  openCreateModal: () => set({ isCreateModalOpen: true }),
  closeCreateModal: () => set({ isCreateModalOpen: false }),
  openEditModal: (id) => set({ selectedSiteId: id, isEditModalOpen: true }),
  closeEditModal: () => set({ selectedSiteId: null, isEditModalOpen: false }),
  openDetailsModal: (id) => set({ selectedSiteId: id, isDetailsModalOpen: true }),
  closeDetailsModal: () => set({ selectedSiteId: null, isDetailsModalOpen: false }),
}))