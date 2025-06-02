import { useState } from 'react'
import { useAuthStore } from '../hooks/useAuthStore'
import { Link } from 'react-router-dom'

function LandingPage() {
  const { user } = useAuthStore()
  const [activeTab, setActiveTab] = useState('summary')

  // Mock data for demonstration
  const deviceStats = {
    total: 24,
    online: 21,
    offline: 3,
    alert: 2,
  }

  const recentActivity = [
    { id: 1, device: 'Temperature Sensor', action: 'Alert triggered', time: '2 mins ago', status: 'alert' },
    { id: 2, device: 'Humidity Sensor', action: 'Data updated', time: '15 mins ago', status: 'success' },
    { id: 3, device: 'Pressure Monitor', action: 'Device offline', time: '1 hour ago', status: 'warning' },
    { id: 4, device: 'Security Camera', action: 'Status changed', time: '3 hours ago', status: 'info' },
  ]

  return (
    <div>
      {/* Hero section with welcome message */}
      <div className="mb-8">
        <h1 className="text-3xl md:text-4xl font-bold mb-4 text-secondary-900 dark:text-secondary-50">
          Welcome{user ? `, ${user.role}` : ''}!
        </h1>
        <p className="text-lg text-secondary-600 dark:text-secondary-400">
          Your IoT device management dashboard
        </p>
      </div>

      {/* Stats overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5 mb-8">
        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-soft p-6">
          <div className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm font-medium text-secondary-500 dark:text-secondary-400">Total Devices</p>
              <h3 className="text-2xl font-bold mt-1 text-secondary-900 dark:text-secondary-50">{deviceStats.total}</h3>
            </div>
            <div className="p-3 rounded-full bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
            </div>
          </div>
          <div className="mt-4 flex items-center">
            <span className="text-sm text-green-600 dark:text-green-400 bg-green-100 dark:bg-green-900/30 px-2 py-0.5 rounded-full">
              +12% growth
            </span>
          </div>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-soft p-6">
          <div className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm font-medium text-secondary-500 dark:text-secondary-400">Online</p>
              <h3 className="text-2xl font-bold mt-1 text-secondary-900 dark:text-secondary-50">{deviceStats.online}</h3>
            </div>
            <div className="p-3 rounded-full bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z" />
              </svg>
            </div>
          </div>
          <div className="mt-4 flex items-center">
            <div className="w-full bg-secondary-100 dark:bg-secondary-700 rounded-full h-2">
              <div className="bg-green-500 h-2 rounded-full" style={{ width: `${Math.round((deviceStats.online / deviceStats.total) * 100)}%` }}></div>
            </div>
            <span className="ml-2 text-sm text-secondary-600 dark:text-secondary-400">
              {Math.round((deviceStats.online / deviceStats.total) * 100)}%
            </span>
          </div>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-soft p-6">
          <div className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm font-medium text-secondary-500 dark:text-secondary-400">Offline</p>
              <h3 className="text-2xl font-bold mt-1 text-secondary-900 dark:text-secondary-50">{deviceStats.offline}</h3>
            </div>
            <div className="p-3 rounded-full bg-secondary-100 dark:bg-secondary-700 text-secondary-600 dark:text-secondary-400">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M18.364 5.636a9 9 0 010 12.728m0 0l-2.829-2.829m2.829 2.829L21 21M15.536 8.464a5 5 0 010 7.072m0 0l-2.829-2.829m-4.243 2.829a4.978 4.978 0 01-1.414-2.83m-1.414 5.658a9 9 0 01-2.167-9.238m7.824 2.167a1 1 0 111.414 1.414m-1.414-1.414L3 3m8.293 8.293l1.414 1.414" />
              </svg>
            </div>
          </div>
          <div className="mt-4 flex items-center">
            <div className="w-full bg-secondary-100 dark:bg-secondary-700 rounded-full h-2">
              <div className="bg-secondary-500 h-2 rounded-full" style={{ width: `${Math.round((deviceStats.offline / deviceStats.total) * 100)}%` }}></div>
            </div>
            <span className="ml-2 text-sm text-secondary-600 dark:text-secondary-400">
              {Math.round((deviceStats.offline / deviceStats.total) * 100)}%
            </span>
          </div>
        </div>

        <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-soft p-6">
          <div className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm font-medium text-secondary-500 dark:text-secondary-400">Alerts</p>
              <h3 className="text-2xl font-bold mt-1 text-secondary-900 dark:text-secondary-50">{deviceStats.alert}</h3>
            </div>
            <div className="p-3 rounded-full bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
          </div>
          <div className="mt-4">
            <Link to="/alerts" className="text-sm font-medium text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300">
              View all alerts →
            </Link>
          </div>
        </div>
      </div>

      {/* Main content tabs */}
      <div className="bg-white dark:bg-secondary-800 rounded-xl shadow-soft mb-8">
        <div className="border-b border-secondary-200 dark:border-secondary-700">
          <nav className="flex">
            <button
              className={`px-4 py-4 font-medium text-sm border-b-2 ${
                activeTab === 'summary'
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-secondary-500 hover:text-secondary-700 dark:text-secondary-400 dark:hover:text-secondary-200'
              }`}
              onClick={() => setActiveTab('summary')}
            >
              Summary
            </button>
            <button
              className={`px-4 py-4 font-medium text-sm border-b-2 ${
                activeTab === 'activity'
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-secondary-500 hover:text-secondary-700 dark:text-secondary-400 dark:hover:text-secondary-200'
              }`}
              onClick={() => setActiveTab('activity')}
            >
              Recent Activity
            </button>
            <button
              className={`px-4 py-4 font-medium text-sm border-b-2 ${
                activeTab === 'analytics'
                  ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                  : 'border-transparent text-secondary-500 hover:text-secondary-700 dark:text-secondary-400 dark:hover:text-secondary-200'
              }`}
              onClick={() => setActiveTab('analytics')}
            >
              Analytics
            </button>
          </nav>
        </div>
        
        <div className="p-6">
          {activeTab === 'summary' && (
            <div>
              <h3 className="text-lg font-medium mb-4 text-secondary-900 dark:text-secondary-50">System Overview</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="border border-secondary-200 dark:border-secondary-700 rounded-lg p-5">
                  <h4 className="font-medium mb-2 text-secondary-900 dark:text-secondary-50">Device Health</h4>
                  <div className="mt-4 flex items-center justify-center">
                    <div className="h-36 w-36 rounded-full border-8 border-primary-100 dark:border-primary-900/30 flex items-center justify-center">
                      <div className="text-center">
                        <span className="block text-3xl font-bold text-primary-600 dark:text-primary-400">87%</span>
                        <span className="text-xs text-secondary-500 dark:text-secondary-400">Overall Health</span>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="border border-secondary-200 dark:border-secondary-700 rounded-lg p-5">
                  <h4 className="font-medium mb-2 text-secondary-900 dark:text-secondary-50">Data Transmitted</h4>
                  <div className="text-center mt-4">
                    <span className="text-3xl font-bold text-secondary-900 dark:text-secondary-50">1.45</span>
                    <span className="text-xl text-secondary-700 dark:text-secondary-300 ml-1">GB</span>
                    <div className="mt-2 text-xs text-green-600 dark:text-green-400">
                      ↑ 12% compared to yesterday
                    </div>
                  </div>
                  <div className="mt-4 h-12 flex items-end space-x-1">
                    {[45, 30, 60, 70, 55, 65, 75, 60, 80, 75, 65, 85].map((h, i) => (
                      <div 
                        key={i} 
                        className="flex-1 bg-primary-100 dark:bg-primary-900/30" 
                        style={{ height: `${h}%` }}
                      ></div>
                    ))}
                  </div>
                </div>

                <div className="border border-secondary-200 dark:border-secondary-700 rounded-lg p-5">
                  <h4 className="font-medium mb-2 text-secondary-900 dark:text-secondary-50">Quick Actions</h4>
                  <div className="mt-4 space-y-3">
                    <button className="w-full px-4 py-2 flex items-center text-sm text-left rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-700 transition-colors">
                      <svg className="w-4 h-4 mr-3 text-primary-500 dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                      </svg>
                      Add new device
                    </button>
                    <button className="w-full px-4 py-2 flex items-center text-sm text-left rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-700 transition-colors">
                      <svg className="w-4 h-4 mr-3 text-primary-500 dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 21H5v-6l1.293-1.293A6.001 6.001 0 1118 9z" />
                      </svg>
                      Generate report
                    </button>
                    <button className="w-full px-4 py-2 flex items-center text-sm text-left rounded-lg hover:bg-secondary-100 dark:hover:bg-secondary-700 transition-colors">
                      <svg className="w-4 h-4 mr-3 text-primary-500 dark:text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      </svg>
                      Configure settings
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}
          
          {activeTab === 'activity' && (
            <div>
              <h3 className="text-lg font-medium mb-4 text-secondary-900 dark:text-secondary-50">Recent Activities</h3>
              <div className="overflow-hidden">
                <table className="min-w-full">
                  <thead>
                    <tr className="border-b border-secondary-200 dark:border-secondary-700">
                      <th className="py-3 px-4 text-left text-xs font-medium text-secondary-500 dark:text-secondary-400 uppercase tracking-wider">Device</th>
                      <th className="py-3 px-4 text-left text-xs font-medium text-secondary-500 dark:text-secondary-400 uppercase tracking-wider">Action</th>
                      <th className="py-3 px-4 text-left text-xs font-medium text-secondary-500 dark:text-secondary-400 uppercase tracking-wider">Time</th>
                      <th className="py-3 px-4 text-left text-xs font-medium text-secondary-500 dark:text-secondary-400 uppercase tracking-wider">Status</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white dark:bg-secondary-800">
                    {recentActivity.map((item) => (
                      <tr key={item.id} className="border-b border-secondary-200 dark:border-secondary-700 last:border-0">
                        <td className="py-4 px-4 whitespace-nowrap text-sm font-medium text-secondary-900 dark:text-secondary-50">{item.device}</td>
                        <td className="py-4 px-4 whitespace-nowrap text-sm text-secondary-600 dark:text-secondary-400">{item.action}</td>
                        <td className="py-4 px-4 whitespace-nowrap text-sm text-secondary-600 dark:text-secondary-400">{item.time}</td>
                        <td className="py-4 px-4 whitespace-nowrap">
                          <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
                            ${item.status === 'alert' ? 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300' :
                              item.status === 'warning' ? 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300' :
                              item.status === 'success' ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300' :
                              'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300'}
                          `}>
                            {item.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <div className="mt-4 text-right">
                <Link to="/activity" className="text-sm font-medium text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300">
                  View all activity →
                </Link>
              </div>
            </div>
          )}
          
          {activeTab === 'analytics' && (
            <div className="text-center py-8">
              <svg className="mx-auto h-16 w-16 text-secondary-400 dark:text-secondary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              <h3 className="mt-4 text-lg font-medium text-secondary-900 dark:text-secondary-50">Analytics Dashboard</h3>
              <p className="mt-2 text-secondary-600 dark:text-secondary-400">Coming soon! Advanced analytics and data visualization.</p>
              <button className="mt-4 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors">
                Request early access
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default LandingPage