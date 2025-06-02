import { Dialog, Transition } from "@headlessui/react";
import { Fragment } from "react";
import type { ValueStream } from "../types/valueStream";

interface ValueStreamDetailsModalProps {
  isOpen: boolean;
  onClose: () => void;
  valueStream: ValueStream | null;
}

function ValueStreamDetailsModal({
  isOpen,
  onClose,
  valueStream,
}: ValueStreamDetailsModalProps) {
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
              <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white dark:bg-gray-800 p-6 text-left align-middle shadow-xl transition-all">
                <Dialog.Title
                  as="h3"
                  className="text-lg font-medium leading-6 text-secondary-900 dark:text-secondary-50"
                >
                  Value Stream: {valueStream?.name || "N/A"}
                </Dialog.Title>
                <div className="mt-4">
                  {valueStream ? (
                    <dl className="space-y-4">
                      <div>
                        <dt className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                          Name
                        </dt>
                        <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
                          {valueStream.name}
                        </dd>
                      </div>
                      <div>
                        <dt className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                          Type
                        </dt>
                        <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
                          {valueStream.type}
                        </dd>
                      </div>
                    </dl>
                  ) : (
                    <p className="text-sm text-red-600 dark:text-red-400">
                      Value Stream not found
                    </p>
                  )}
                </div>
                <div className="mt-6 flex justify-end">
                  <button
                    type="button"
                    onClick={onClose}
                    className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-gray-100 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600"
                  >
                    Close
                  </button>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}

export default ValueStreamDetailsModal;
