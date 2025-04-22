import Link from 'next/link';
import {Activity} from '@/types/api';
import ProtectedRoute from "@/components/ProtectedRoute";
import {api} from "@/lib/api";

interface HomeProps {
  deviceCount: number;
  recentActivities: Activity[];
}

export default async function Home() {
const deviceCount = (await api.getDevices()).length;
let recentActivities: Activity[] = [];

  try {
    //recentActivities = data.RecentActivities;
  } catch (error) {
    console.error('Failed to fetch main data:', error);
  }

  return (
      <ProtectedRoute>
        <div className="py-8 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto">
          <div className="mb-12">
            <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] text-center border-0 rounded-lg overflow-hidden">
              <div className="py-12 px-6 relative">
                <div className="absolute inset-0 bg-gradient-to-r from-blue-500/10 to-teal-500/10"></div>
                <div className="relative z-10">
                  <h2 className="text-4xl sm:text-5xl font-bold text-[--color-text-primary] mb-4">
                    Welcome to <span className="text-[--color-text-accent]">IoTGo</span> Platform
                  </h2>
                  <p className="text-lg text-[--color-text-secondary] max-w-2xl mx-auto">
                    Your comprehensive solution for IoT device management and monitoring.
                  </p>
                  <div className="mt-8">
                    <Link
                        href="/devices/new"
                        className="bg-[--color-text-accent] text-white font-semibold py-3 px-6 rounded-lg hover:bg-opacity-90 transition-all inline-flex items-center"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
                      </svg>
                      Add New Device
                    </Link>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] border-0 rounded-lg hover:-translate-y-1 hover:shadow-[0_12px_20px_var(--color-card-shadow)] transition-all duration-300">
              <div className="p-6">
                <div className="flex items-center mb-4">
                  <div className="p-3 bg-blue-500/10 rounded-lg mr-4">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-[--color-text-accent]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <h5 className="text-xl font-semibold text-[--color-text-primary]">
                    Devices
                  </h5>
                </div>
                <h2 className="text-4xl font-bold text-[--color-text-accent] mb-2">
                  {deviceCount}
                </h2>
                <p className="text-[--color-text-secondary] mb-4">Total connected devices</p>
                <Link
                    href="/devices"
                    className="inline-flex items-center text-[--color-text-accent] hover:underline font-medium"
                >
                  Manage Devices
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clipRule="evenodd" />
                  </svg>
                </Link>
              </div>
            </div>

            <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] border-0 rounded-lg hover:-translate-y-1 hover:shadow-[0_12px_20px_var(--color-card-shadow)] transition-all duration-300">
              <div className="p-6">
                <div className="flex items-center mb-4">
                  <div className="p-3 bg-green-500/10 rounded-lg mr-4">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-[--color-success]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                    </svg>
                  </div>
                  <h5 className="text-xl font-semibold text-[--color-text-primary]">
                    System Status
                  </h5>
                </div>
                <div className="flex items-center mb-4">
                  <svg xmlns="http://www.w3.org/2000/svg" className="text-[--color-success] h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  <span className="text-[--color-text-primary]">
                  All systems operational
                </span>
                </div>
                <div className="space-x-2">
                <span className="bg-[--color-success] text-white py-1 px-3 rounded-full text-sm font-medium">
                  API Online
                </span>
                  <span className="bg-[--color-success] text-white py-1 px-3 rounded-full text-sm font-medium">
                  Database Connected
                </span>
                </div>
              </div>
            </div>

            <div className="bg-[--color-background-secondary] shadow-[0_8px_16px_var(--color-card-shadow)] border-0 rounded-lg hover:-translate-y-1 hover:shadow-[0_12px_20px_var(--color-card-shadow)] transition-all duration-300">
              <div className="p-6">
                <div className="flex items-center mb-4">
                  <div className="p-3 bg-purple-500/10 rounded-lg mr-4">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <h5 className="text-xl font-semibold text-[--color-text-primary]">
                    Recent Activity
                  </h5>
                </div>
                <ul className="space-y-2 mt-2">
                  {recentActivities.length > 0 ? (
                      recentActivities.map((activity, index) => (
                          <li
                              key={index}
                              className="flex justify-between items-center py-2 border-b border-gray-700/30"
                          >
                      <span className="text-[--color-text-primary]">
                        {activity.Description}
                      </span>
                            <small className="text-[--color-text-secondary] ml-2">
                              {activity.Time}
                            </small>
                          </li>
                      ))
                  ) : (
                      <li className="py-3 text-[--color-text-secondary] text-center border border-dashed border-gray-700/30 rounded-lg bg-gray-700/10">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mx-auto mb-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        No recent activities
                      </li>
                  )}
                </ul>
              </div>
            </div>
          </div>
        </div>
      </ProtectedRoute>
  );
}