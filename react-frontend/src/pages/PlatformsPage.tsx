import { useState } from "react";
import { usePlatformStore } from "../hooks/usePlatformStore";
import {
  usePlatforms,
  usePlatform,
  useCreatePlatform,
  useUpdatePlatform,
  useDeletePlatform,
  useTestPlatformConnection,
  useResources,
  useResourceForEdit,
  useCreateResource,
  useUpdateResource,
  useDeleteResource,
  useTestResource,
  useFetchDeviceData,
  useDevices,
} from "../hooks/useQueries";
import PlatformTable from "../components/PlatformTable";
import PlatformModal from "../components/PlatformModal";
import ResourceTable from "../components/ResourceTable";
import ResourceModal from "../components/ResourceModal";
import TestConnectionModal from "../components/TestConnectionModal";
import FetchDeviceDataModal from "../components/FetchDeviceDataModal";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { dracula } from "react-syntax-highlighter/dist/esm/styles/prism";
import type { Platform, Resource, TestConnectionInput } from "../types/platform";

function PlatformsPage() {
  const {
    selectedPlatformId,
    selectedResourceId,
    activeTab,
    isPlatformCreateModalOpen,
    isPlatformEditModalOpen,
    isResourceCreateModalOpen,
    isResourceEditModalOpen,
    isTestConnectionModalOpen,
    isFetchDeviceDataModalOpen,
    setSelectedPlatformId,
    setActiveTab,
    openPlatformCreateModal,
    closePlatformCreateModal,
    openPlatformEditModal,
    closePlatformEditModal,
    openResourceCreateModal,
    closeResourceCreateModal,
    openResourceEditModal,
    closeResourceEditModal,
    openTestConnectionModal,
    closeTestConnectionModal,
    openFetchDeviceDataModal,
    closeFetchDeviceDataModal,
  } = usePlatformStore();

  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [nameFilter, setNameFilter] = useState("");
  // const [sort, setSort] = useState("name");

  const {
    data: platformsData,
    isLoading: isPlatformsLoading,
    error: platformsError,
  } = usePlatforms({
    limit,
    offset: (page - 1) * limit,
    name: nameFilter,
    sort: "name",
  });
  const { data: selectedPlatform, isLoading: isPlatformLoading } = usePlatform(
    selectedPlatformId || 0,
  );
  const {
    data: resourcesData,
    isLoading: isResourcesLoading,
    error: resourcesError,
  } = useResources(selectedPlatformId || 0, {
    limit: 100,
    offset: 0,
    name: "",
    sort: "name",
  });
  const { data: resourceForEdit, isLoading: isResourceEditLoading } =
    useResourceForEdit(selectedResourceId || 0);
  const { data: devicesData } = useDevices({ limit: 100, offset: 0 });

  const createPlatformMutation = useCreatePlatform();
  const updatePlatformMutation = useUpdatePlatform();
  const deletePlatformMutation = useDeletePlatform();
  const testConnectionMutation = useTestPlatformConnection();
  const createResourceMutation = useCreateResource();
  const updateResourceMutation = useUpdateResource();
  const deleteResourceMutation = useDeleteResource();
  const testResourceMutation = useTestResource();
  const fetchDeviceDataMutation = useFetchDeviceData();

  const handleCreatePlatform = (platform: Partial<Platform>) => {
    createPlatformMutation.mutate(platform, {
      onSuccess: () => closePlatformCreateModal(),
    });
  };

  const handleUpdatePlatform = (platform: Partial<Platform>) => {
    if (selectedPlatformId) {
      updatePlatformMutation.mutate(
        { id: selectedPlatformId, platform },
        {
          onSuccess: () => closePlatformEditModal(),
        },
      );
    }
  };

  const handleDeletePlatform = (id: number) => {
    if (window.confirm("Are you sure you want to delete this platform?")) {
      deletePlatformMutation.mutate(id, {
        onSuccess: () => setSelectedPlatformId(null),
      });
    }
  };

  const handleTestConnection = (input: TestConnectionInput) => {
    testConnectionMutation.mutate(input, {
      onSuccess: () => {
        // Do not close the modal to show the result
      },
    });
  };

  const handleCreateResource = (resource: Partial<Resource>) => {
    if (selectedPlatformId) {
      createResourceMutation.mutate(
        { platformId: selectedPlatformId, resource },
        {
          onSuccess: () => closeResourceCreateModal(),
        },
      );
    }
  };

  const handleUpdateResource = (resource: Partial<Resource>) => {
    if (selectedResourceId) {
      updateResourceMutation.mutate(
        { id: selectedResourceId, resource },
        {
          onSuccess: () => closeResourceEditModal(),
        },
      );
    }
  };

  const handleDeleteResource = (id: number) => {
    if (
      selectedPlatformId &&
      window.confirm("Are you sure you want to delete this resource?")
    ) {
      deleteResourceMutation.mutate(
        { id, platformId: selectedPlatformId },
        {
          onSuccess: () => closeResourceEditModal(),
        },
      );
    }
  };

  const handleTestResource = (id: number) => {
    testResourceMutation.mutate(id);
  };

  const handleFetchDeviceData = (platformId: number, deviceId: number) => {
    fetchDeviceDataMutation.mutate({ platformId, deviceId });
  };

  // const handleSort = (newSort: string) => {
  //   setSort(newSort);
  // };

  const totalPages = platformsData ? Math.ceil(platformsData.total / limit) : 1;
  const isLoading =
    isPlatformsLoading ||
    isPlatformLoading ||
    isResourcesLoading ||
    isResourceEditLoading;

  const tabs = [
    { id: "details", label: "Details" },
    { id: "resources", label: "Resources" },
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-secondary-900 dark:text-secondary-50">
          Platform Management
        </h2>
        <button
          onClick={openPlatformCreateModal}
          className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
        >
          Add Platform
        </button>
      </div>

      {platformsError && (
        <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
          <p className="text-sm text-red-600 dark:text-red-400">
            {platformsError.message || "Failed to load platforms"}
          </p>
        </div>
      )}

      <div className="flex items-center space-x-4">
        <input
          type="text"
          value={nameFilter}
          onChange={(e) => setNameFilter(e.target.value)}
          placeholder="Filter by name..."
          className="px-4 py-2 border border-secondary-200 dark:border-secondary-700 rounded-lg dark:bg-secondary-800 dark:text-secondary-100"
        />
      </div>

      {isLoading ? (
        <div>Loading...</div>
      ) : (
        <div className="flex flex-col md:flex-row gap-6">
          {/* Platform List */}
          <div className="w-full md:w-1/2">
            <PlatformTable
              platforms={platformsData?.items || []}
              onEdit={openPlatformEditModal}
              onDelete={handleDeletePlatform}
              onSelect={setSelectedPlatformId}
            />
            <div className="flex justify-between items-center mt-4">
              <button
                onClick={() => setPage((p) => Math.max(p - 1, 1))}
                disabled={page === 1}
                className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg disabled:opacity-50"
              >
                Previous
              </button>
              <span>
                Page {page} of {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(p + 1, totalPages))}
                disabled={page === totalPages}
                className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg disabled:opacity-50"
              >
                Next
              </button>
            </div>
          </div>

          {/* Platform Details and Resources */}
          {selectedPlatformId && (
            <div className="w-full md:w-1/2">
              <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-lg font-medium text-secondary-900 dark:text-secondary-50">
                    {selectedPlatform?.name}
                  </h3>
                  <div className="space-x-2">
                    <button
                      onClick={openTestConnectionModal}
                      className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
                    >
                      Test Connection
                    </button>
                    <button
                      onClick={openFetchDeviceDataModal}
                      className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                    >
                      Fetch Device Data
                    </button>
                  </div>
                </div>
                <nav className="flex space-x-4 mb-4">
                  {tabs.map((tab) => (
                    <button
                      key={tab.id}
                      onClick={() => setActiveTab(tab.id)}
                      className={`px-4 py-2 text-sm rounded-lg ${
                        activeTab === tab.id
                          ? "bg-primary-600 text-white"
                          : "text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                      }`}
                    >
                      {tab.label}
                    </button>
                  ))}
                </nav>
                {activeTab === "details" && (
                  <div className="space-y-4">
                    <dl className="space-y-4">
                      <div>
                        <dt className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                          Name
                        </dt>
                        <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
                          {selectedPlatform?.name}
                        </dd>
                      </div>
                      <div>
                        <dt className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                          Type
                        </dt>
                        <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
                          {selectedPlatform?.type}
                        </dd>
                      </div>
                      <div>
                        <dt className="text-sm font-medium text-secondary-700 dark:text-secondary-300">
                          Metadata
                        </dt>
                        <dd className="mt-1 text-sm text-gray-900 dark:text-gray-100">
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
                            {JSON.stringify(
                              JSON.parse(selectedPlatform?.metadata || "{}"),
                              null,
                              2,
                            )}
                          </SyntaxHighlighter>
                        </dd>
                      </div>
                    </dl>
                  </div>
                )}
                {activeTab === "resources" && (
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <h4 className="text-md font-medium text-secondary-900 dark:text-secondary-50">
                        Resources
                      </h4>
                      <button
                        onClick={openResourceCreateModal}
                        className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
                      >
                        Add Resource
                      </button>
                    </div>
                    {resourcesError && (
                      <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
                        <p className="text-sm text-red-600 dark:text-red-400">
                          {resourcesError.message || "Failed to load resources"}
                        </p>
                      </div>
                    )}
                    {isResourcesLoading ? (
                      <div>Loading...</div>
                    ) : (
                      <ResourceTable
                        resources={resourcesData?.items || []}
                        onEdit={openResourceEditModal}
                        onDelete={handleDeleteResource}
                        onTest={handleTestResource}
                      />
                    )}
                    {testResourceMutation.data && (
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
                          {JSON.stringify(testResourceMutation.data, null, 2)}
                        </SyntaxHighlighter>
                      </div>
                    )}
                    {testResourceMutation.error && (
                      <div className="mt-4 p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
                        <p className="text-sm text-red-600 dark:text-red-400">
                          {testResourceMutation.error.message ||
                            "Failed to test resource"}
                        </p>
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      )}

      <PlatformModal
        isOpen={isPlatformCreateModalOpen}
        onClose={closePlatformCreateModal}
        onSubmit={handleCreatePlatform}
        error={createPlatformMutation.error?.message}
        isLoading={createPlatformMutation.isPending}
      />

      <PlatformModal
        isOpen={isPlatformEditModalOpen}
        onClose={closePlatformEditModal}
        platform={selectedPlatform}
        onSubmit={handleUpdatePlatform}
        error={updatePlatformMutation.error?.message}
        isLoading={updatePlatformMutation.isPending}
      />

      <ResourceModal
        isOpen={isResourceCreateModalOpen}
        onClose={closeResourceCreateModal}
        platformType={selectedPlatform?.type || "REST"}
        onSubmit={handleCreateResource}
        error={createResourceMutation.error?.message}
        isLoading={createResourceMutation.isPending}
      />

      <ResourceModal
        isOpen={isResourceEditModalOpen}
        onClose={closeResourceEditModal}
        resource={resourceForEdit}
        platformType={selectedPlatform?.type || "REST"}
        onSubmit={handleUpdateResource}
        error={updateResourceMutation.error?.message}
        isLoading={updateResourceMutation.isPending}
      />

      <TestConnectionModal
        isOpen={isTestConnectionModalOpen}
        onClose={closeTestConnectionModal}
        platform={selectedPlatform}
        onSubmit={handleTestConnection}
        error={testConnectionMutation.error?.message}
        isLoading={testConnectionMutation.isPending}
        testResult={testConnectionMutation.data}
      />

      <FetchDeviceDataModal
        isOpen={isFetchDeviceDataModalOpen}
        onClose={closeFetchDeviceDataModal}
        platformId={selectedPlatformId || 0}
        devices={
          devicesData?.items.map((d) => ({ id: d.id, name: d.name })) || []
        }
        onSubmit={handleFetchDeviceData}
        result={fetchDeviceDataMutation.data}
        error={fetchDeviceDataMutation.error?.message}
        isLoading={fetchDeviceDataMutation.isPending}
      />
    </div>
  );
}

export default PlatformsPage;
