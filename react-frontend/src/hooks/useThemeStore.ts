import { create } from 'zustand'

interface ThemeState {
  theme: 'light' | 'dark'
  toggleTheme: () => void
}

export const useThemeStore = create<ThemeState>((set) => {
  const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
  const initialTheme = savedTheme || 'light'

  return {
    theme: initialTheme,
    toggleTheme: () => {
      set((state) => {
        const newTheme = state.theme === 'light' ? 'dark' : 'light'
        localStorage.setItem('theme', newTheme)
        document.documentElement.classList.toggle('dark', newTheme === 'dark')
        return { theme: newTheme }
      })
    },
  }
})