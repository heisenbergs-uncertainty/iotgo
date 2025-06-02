import { create } from "zustand";

interface PlatformState {
  selectedPlatformId: number | null;
  selectedResourceId: number | null;
  activeTab: string;
  isPlatformCreateModalOpen: boolean;
  isPlatformEditModalOpen: boolean;
  isResourceCreateModalOpen: boolean;
  isResourceEditModalOpen: boolean;
  isTestConnectionModalOpen: boolean;
  isFetchDeviceDataModalOpen: boolean;
  setSelectedPlatformId: (id: number | null) => void;
  setSelectedResourceId: (id: number | null) => void;
  setActiveTab: (tab: string) => void;
  openPlatformCreateModal: () => void;
  closePlatformCreateModal: () => void;
  openPlatformEditModal: (id: number) => void;
  closePlatformEditModal: () => void;
  openResourceCreateModal: () => void;
  closeResourceCreateModal: () => void;
  openResourceEditModal: (id: number) => void;
  closeResourceEditModal: () => void;
  openTestConnectionModal: () => void;
  closeTestConnectionModal: () => void;
  openFetchDeviceDataModal: () => void;
  closeFetchDeviceDataModal: () => void;
}

export const usePlatformStore = create<PlatformState>((set) => ({
  selectedPlatformId: null,
  selectedResourceId: null,
  activeTab: "details",
  isPlatformCreateModalOpen: false,
  isPlatformEditModalOpen: false,
  isResourceCreateModalOpen: false,
  isResourceEditModalOpen: false,
  isTestConnectionModalOpen: false,
  isFetchDeviceDataModalOpen: false,
  setSelectedPlatformId: (id) =>
    set({
      selectedPlatformId: id,
      selectedResourceId: null,
      activeTab: "details",
    }),
  setSelectedResourceId: (id) => set({ selectedResourceId: id }),
  setActiveTab: (tab) => set({ activeTab: tab }),
  openPlatformCreateModal: () => set({ isPlatformCreateModalOpen: true }),
  closePlatformCreateModal: () => set({ isPlatformCreateModalOpen: false }),
  openPlatformEditModal: (id) =>
    set({ selectedPlatformId: id, isPlatformEditModalOpen: true }),
  closePlatformEditModal: () => set({ isPlatformEditModalOpen: false }),
  openResourceCreateModal: () => set({ isResourceCreateModalOpen: true }),
  closeResourceCreateModal: () => set({ isResourceCreateModalOpen: false }),
  openResourceEditModal: (id) =>
    set({ selectedResourceId: id, isResourceEditModalOpen: true }),
  closeResourceEditModal: () => set({ isResourceEditModalOpen: false }),
  openTestConnectionModal: () => set({ isTestConnectionModalOpen: true }),
  closeTestConnectionModal: () => set({ isTestConnectionModalOpen: false }),
  openFetchDeviceDataModal: () => set({ isFetchDeviceDataModalOpen: true }),
  closeFetchDeviceDataModal: () => set({ isFetchDeviceDataModalOpen: false }),
}));
