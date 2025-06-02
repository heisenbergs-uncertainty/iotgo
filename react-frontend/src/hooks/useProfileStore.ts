import { create } from "zustand";

interface ProfileState {
  selectedTab: string;
  isApiKeyModalOpen: boolean;
  setSelectedTab: (tab: string) => void;
  openApiKeyModal: () => void;
  closeApiKeyModal: () => void;
}

export const useProfileStore = create<ProfileState>((set) => ({
  selectedTab: "general",
  isApiKeyModalOpen: false,
  setSelectedTab: (tab) => set({ selectedTab: tab }),
  openApiKeyModal: () => set({ isApiKeyModalOpen: true }),
  closeApiKeyModal: () => set({ isApiKeyModalOpen: false }),
}));
