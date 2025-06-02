import { Dialog, Transition } from "@headlessui/react";
import { Fragment, useState } from "react";
import type { FetchDeviceDataResponse } from "../types/platform";

interface FetchDeviceDataModalProps {
  isOpen: boolean;
  onClose: () => void;
  platformId: number;
  devices: { id: number; name: string }[];
  onSubmit: (platformId: number, deviceId: number) => void;
  result?: FetchDeviceDataResponse;
  error?: string;
  isLoading?: boolean;
}

function FetchDeviceDataModal({
  isOpen,
  onClose,
  platformId,
  devices,
  onSubmit,
  result,
  error,
  isLoading,
}: FetchDeviceDataModalProps) {
  const [deviceId, setDeviceId] = useState<number | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (deviceId) {
      onSubmit(platformId, deviceId);
    }
  };

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={onClose}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black bg-opacity-25" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4 text-center">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="w-full max-w-lg transform overflow-hidden rounded-2xl bg-white dark:bg-gray-800 p-6 text-left align-middle shadow-xl transition-all">
                <Dialog.Title
                  as="h3"
                  className="text-lg font-medium leading-6 text-secondary-900 dark:text-secondary-50"
                >
                  Fetch Device Data
                </Dialog.Title>
                <div className="mt-4">
                  <form onSubmit={handleSubmit} className="space-y-6">
                    {error && (
                      <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
                        <p className="text-sm text-red-600 dark:text-red-400">
                          {error}
                        </p>
                      </div>
                    )}
                    <div>
                      <label className="block text-sm font-medium mb-2 text-secondary-700 dark:text-secondary-300">
                        Device
                      </label>
                      <select
                        name="deviceId"
                        title="Select a device"
                        value={deviceId || ""}
                        onChange={(e) =>
                          setDeviceId(parseInt(e.target.value) || null)
                        }
                        className="block w-full px-4 py-3 rounded-lg border border-secondary-200 dark:border-secondary-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-secondary-800 dark:text-secondary-100"
                        required
                      >
                        <option value="">Select a device</option>
                        {devices.map((device) => (
                          <option key={device.id} value={device.id}>
                            {device.name}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div className="flex justify-end space-x-4">
                      <button
                        type="button"
                        onClick={onClose}
                        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600"
                      >
                        Cancel
                      </button>
                      <button
                        type="submit"
                        disabled={isLoading}
                        className={`px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 ${isLoading ? "opacity-80 cursor-not-allowed" : ""}`}
                      >
                        {isLoading ? "Fetching..." : "Fetch Data"}
                      </button>
                    </div>
                    {result && (
                      <div className="mt-4">
                        <h4 className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                          Results
                        </h4>
                        <pre className="mt-2 p-4 bg-gray-100 dark:bg-gray-900 rounded-lg text-sm text-gray-900 dark:text-gray-100 overflow-auto">
                          {JSON.stringify(result, null, 2)}
                        </pre>
                      </div>
                    )}
                  </form>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}

export default FetchDeviceDataModal;
