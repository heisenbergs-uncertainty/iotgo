import { Dialog, Transition } from "@headlessui/react";
import { Fragment } from "react";
import ResourceForm from "./ResourceForm";
import type { PlatformType, Resource, ResourceType } from "../types/platform";

interface ResourceModalProps {
  isOpen: boolean;
  onClose: () => void;
  resource?: { id: number; name: string; type: ResourceType; details: any };
  platformType: PlatformType;
  onSubmit: (resource: Partial<Resource>) => void;
  error?: string;
  isLoading?: boolean;
}

function ResourceModal({
  isOpen,
  onClose,
  resource,
  platformType,
  onSubmit,
  error,
  isLoading,
}: ResourceModalProps) {
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
                  {resource?.id ? "Edit Resource" : "Create Resource"}
                </Dialog.Title>
                <div className="mt-4">
                  <ResourceForm
                    resource={resource}
                    platformType={platformType}
                    onSubmit={onSubmit}
                    onCancel={onClose}
                    error={error}
                    isLoading={isLoading}
                  />
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
}

export default ResourceModal;
