import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { useAuthStore } from '../hooks/useAuthStore'
import type { LoginResponse } from '../types/auth'

function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()
  const { login } = useAuthStore()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')
    try {
      const response = await axios.post<LoginResponse>('/auth/login', {
        email,
        password,
      })
      if (response.status === 200 && response.data.token) {
        login(
          response.data.user.id,
          response.data.user.email,
          response.data.user.role,
          response.data.token,
        )
        navigate('/', { replace: true })
      } else {
        setError(response.data.error || 'Invalid email or password')
      }
    } catch (error: any) {
      console.error('Login failed:', error.response?.data || error.message)
      setError(
        error.response?.data?.error || 'Invalid email or password',
      )
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-secondary-50 dark:bg-secondary-900">
      {/* Left side - Branding and information */}
      <div className="w-full md:w-1/2 bg-primary-600 dark:bg-primary-800 flex flex-col justify-center p-8 md:p-12 lg:p-24">
        <div className="max-w-md mx-auto">
          <div className="flex items-center space-x-2 mb-8">
            <svg
              className="w-10 h-10 text-white"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
            <h1 className="text-2xl font-bold text-white">IoTGo</h1>
          </div>
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
            Welcome to IoT Device Management
          </h2>
          <p className="text-primary-100 text-lg mb-8">
            Securely manage and monitor your IoT devices from anywhere. Get
            real-time updates, analytics, and insights all in one place.
          </p>
          <div className="flex space-x-4">
            <div className="flex items-center space-x-2">
              <svg
                className="w-5 h-5 text-primary-200"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M5 13l4 4L19 7"
                />
              </svg>
              <span className="text-primary-100">Real-time monitoring</span>
            </div>
            <div className="flex items-center space-x-2">
              <svg
                className="w-5 h-5 text-primary-200"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M5 13l4 4L19 7"
                />
              </svg>
              <span className="text-primary-100">Secure access</span>
            </div>
          </div>
        </div>
      </div>

      {/* Right side - Login form */}
      <div className="w-full md:w-1/2 flex items-center justify-center p-8 md:p-12 lg:p-24">
        <div className="w-full max-w-md">
          <div className="mb-10">
            <h2 className="text-2xl font-bold mb-2 text-secondary-900 dark:text-secondary-50">
              Sign in to your account
            </h2>
            <p className="text-secondary-600 dark:text-secondary-400">
              Enter your credentials to access your dashboard
            </p>
          </div>

          {error && (
            <div className="mb-6 p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
              <div className="flex">
                <svg
                  className="h-5 w-5 text-red-500 dark:text-red-400 mr-3 mt-0.5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
              </div>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
                Email Address
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100 transition-colors"
                placeholder="you@example.com"
                required
              />
            </div>
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-sm font-medium text-secondary-700 dark:text-secondary-300">
                  Password
                </label>
                <a
                  href="#"
                  className="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-500 dark:hover:text-primary-300"
                >
                  Forgot password?
                </a>
              </div>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100 transition-colors"
                placeholder="••••••••••"
                required
              />
            </div>
            <div className="flex items-center">
              <input
                id="remember-me"
                name="remember-me"
                type="checkbox"
                className="h-4 w-4 rounded border-secondary-300 dark:border-secondary-700 text-primary-600 focus:ring-primary-500"
              />
              <label
                htmlFor="remember-me"
                className="ml-2 block text-sm text-secondary-600 dark:text-secondary-400"
              >
                Remember me
              </label>
            </div>
            <button
              type="submit"
              disabled={isLoading}
              className={`w-full flex justify-center items-center px-4 py-3 rounded-lg bg-primary-600 hover:bg-primary-700 text-white font-medium transition-colors ${
                isLoading ? 'opacity-80 cursor-not-allowed' : ''
              }`}
            >
              {isLoading ? (
                <>
                  <svg
                    className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  Signing in...
                </>
              ) : (
                'Sign in'
              )}
            </button>
          </form>

          <div className="mt-8">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-secondary-200 dark:border-secondary-700"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-secondary-50 dark:bg-secondary-900 text-secondary-500 dark:text-secondary-400">
                  Don't have an account?
                </span>
              </div>
            </div>
            <div className="mt-6">
              <a
                href="#"
                className="w-full flex justify-center items-center px-4 py-3 border border-secondary-300 dark:border-secondary-700 rounded-lg shadow-sm text-sm font-medium text-secondary-700 dark:text-secondary-300 bg-secondary-50 dark:bg-secondary-800 hover:bg-secondary-100 dark:hover:bg-secondary-700 transition-colors"
              >
                Contact your administrator
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginPage