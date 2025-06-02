import { Dialog, Transition } from "@headlessui/react";
import { Fragment } from "react";
import PlatformForm from "./PlatformForm";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { dracula } from "react-syntax-highlighter/dist/esm/styles/prism";
import type { Platform, TestConnectionInput } from "../types/platform";

interface TestConnectionModalProps {
  isOpen: boolean;
  onClose: () => void;
  platform?: Platform;
  onSubmit: (input: TestConnectionInput) => void;
  error?: string;
  isLoading?: boolean;
  testResult?: { message: string }; // Add testResult prop for connection test response
}

function TestConnectionModal({
  isOpen,
  onClose,
  platform,
  onSubmit,
  error,
  isLoading,
  testResult,
}: TestConnectionModalProps) {
  const handleSubmit = (platformData: Partial<Platform>) => {
    onSubmit({
      type: platformData.type as "REST" | "InfluxDB",
      metadata: platformData.metadata || "",
    });
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
                  Test Platform Connection
                </Dialog.Title>
                <div className="mt-4">
                  <PlatformForm
                    platform={platform}
                    onSubmit={handleSubmit}
                    onCancel={onClose}
                    error={error}
                    isLoading={isLoading}
                  />
                  {testResult && (
                    <div className="mt-4">
                      <h4 className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                        Test Result
                      </h4>
                      <SyntaxHighlighter
                        language="json"
                        style={dracula}
                        customStyle={{
                          padding: "1rem",
                          borderRadius: "0.5rem",
                          backgroundColor: "#1a202c",
                          fontSize: "0.875rem",
                          maxHeight: "300px",
                          overflow: "auto",
                        }}
                      >
                        {JSON.stringify(testResult, null, 2)}
                      </SyntaxHighlighter>
                    </div>
                  )}
                  {error && !isLoading && (
                    <div className="mt-4 p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
                      <p className="text-sm text-red-600 dark:text-red-400">
                        {error}
                      </p>
                    </div>
                  )}
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}

export default TestConnectionModal;
