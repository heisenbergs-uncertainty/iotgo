import { create } from 'zustand'
import axios from 'axios'
import type {  User } from '../types/common'
import type { ApiResponse } from '../types/common'

interface AuthState {
  user: User | null
  isCheckingAuth: boolean
  lastAuthCheck: number
  login: (userId: number, email: string, role: string, token: string) => void
  logout: () => void
  checkAuth: () => Promise<void>
}

// Set up Axios interceptor to include Bearer token
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  isCheckingAuth: false,
  lastAuthCheck: 0,
  login: (userId: number, email: string, role: string, token: string) => {
    localStorage.setItem('user_id', userId.toString())
    localStorage.setItem('user_email', email)
    localStorage.setItem('user_role', role)
    localStorage.setItem('token', token)
    set({ user: { id: userId, email, role }, isCheckingAuth: false })
  },
  logout: () => {
    localStorage.removeItem('user_id')
    localStorage.removeItem('user_email')
    localStorage.removeItem('user_role')
    localStorage.removeItem('token')
    set({ user: null, isCheckingAuth: false })
  },
  checkAuth: async () => {
    const now = Date.now()
    const { lastAuthCheck } = get()
    
    // Prevent checking auth too frequently (e.g., within 5 seconds)
    if (now - lastAuthCheck < 5000) {
      console.log('Skipping auth check: too soon since last check')
      return
    }

    set({ isCheckingAuth: true, lastAuthCheck: now })
    const userId = localStorage.getItem('user_id')
    const role = localStorage.getItem('user_role') || ''
    
    if (!userId) {
      set({ isCheckingAuth: false })
      return
    }

    try {
      const response = await axios.get<ApiResponse<User>>(`/api/users/${userId}`)

      if (response.data.code === 200 && response.data.data.id) {
        set({
          user: { id: parseInt(userId), email: localStorage.getItem('user_email') || '', role },
          isCheckingAuth: false
        })
      } else {
        console.error('Invalid user response:', response.data)
        localStorage.removeItem('user_id')
        localStorage.removeItem('user_email')
        localStorage.removeItem('user_role')
        localStorage.removeItem('token')
        set({ user: null, isCheckingAuth: false })
      }
    } catch (err: any) {
      console.error('Auth check failed:', err.response?.data || err.message)
      if (err.response?.status === 429) {
        // Handle rate limiting: pause for 10 seconds
        setTimeout(() => {
          set({ isCheckingAuth: false })
        }, 10000)
      } else if (err.response?.status === 401) {
        // Invalid token: clear auth state
        localStorage.removeItem('user_id')
        localStorage.removeItem('user_role')
        localStorage.removeItem('token')
        set({ user: null, isCheckingAuth: false })
      } else {
        // Other errors: keep user state but stop checking
        set({ isCheckingAuth: false })
      }
    }
  },
}))