import { Dialog, Transition } from '@headlessui/react'
import { Fragment } from 'react'
import DeviceForm from './DeviceForm'
import type {Device} from '../types/device'
import type { Site } from '../types/site'
import type { ValueStream } from '../types/valueStream'

interface DeviceModalProps {
  isOpen: boolean
  onClose: () => void
  device?: Device
  sites: Site[]
  valueStreams: ValueStream[]
  onSubmit: (device: Partial<Device>) => void
  error?: string
  isLoading?: boolean
}

function DeviceModal({ isOpen, onClose, device, sites, valueStreams, onSubmit, error, isLoading }: DeviceModalProps) {
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
                  {device?.id ? 'Edit Device' : 'Create Device'}
                </Dialog.Title>
                <div className="mt-4">
                   <DeviceForm
                    device={device}
                    sites={sites}
                    valueStreams={valueStreams}
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
  )
}

export default DeviceModal