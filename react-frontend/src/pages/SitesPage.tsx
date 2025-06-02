import { useState } from "react";
import { useSiteStore } from "../hooks/useSiteStore";
import {
  useSites,
  useSite,
  useCreateSite,
  useUpdateSite,
  useDeleteSite,
} from "../hooks/useQueries";
import SiteTable from "../components/SiteTable";
import SiteModal from "../components/SiteModal";
import SiteDetailsModal from "../components/SiteDetailsModal";
import type { Site } from "../types/site";

function SitesPage() {
  const {
    isCreateModalOpen,
    isEditModalOpen,
    isDetailsModalOpen,
    selectedSiteId,
    openCreateModal,
    closeCreateModal,
    openEditModal,
    closeEditModal,
    openDetailsModal,
    closeDetailsModal,
  } = useSiteStore();
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [nameFilter, setNameFilter] = useState("");
  // const [sort, setSort] = useState("name");

  const {
    data: sitesData,
    isLoading: isSitesLoading,
    error: sitesError,
  } = useSites({
    limit,
    offset: (page - 1) * limit,
    name: nameFilter,
    sort: "name",
  });
  const { data: selectedSite, isLoading: isSiteLoading } = useSite(
    selectedSiteId || 0,
  );

  const createSiteMutation = useCreateSite();
  const updateSiteMutation = useUpdateSite();
  const deleteSiteMutation = useDeleteSite();

  const handleCreate = (site: Partial<Site>) => {
    createSiteMutation.mutate(site, {
      onSuccess: () => closeCreateModal(),
      onError: (err: any) => console.error("Create error:", err),
    });
  };

  const handleUpdate = (site: Partial<Site>) => {
    if (selectedSiteId) {
      updateSiteMutation.mutate(
        { id: selectedSiteId, site },
        {
          onSuccess: () => closeEditModal(),
          onError: (err: any) => console.error("Update error:", err),
        },
      );
    }
  };

  const handleDelete = (id: number) => {
    if (window.confirm("Are you sure you want to delete this site?")) {
      deleteSiteMutation.mutate(id, {
        onError: (err: any) => console.error("Delete error:", err),
      });
    }
  };

  // const handleSort = (newSort: string) => {
  //   setSort(newSort);
  // };

  const totalPages = sitesData ? Math.ceil(sitesData.total / limit) : 1;
  const isLoading = isSitesLoading || isSiteLoading;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-secondary-900 dark:text-secondary-50">
          Site Management
        </h2>
        <button
          onClick={openCreateModal}
          className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
        >
          Add Site
        </button>
      </div>

      {sitesError && (
        <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
          <p className="text-sm text-red-600 dark:text-red-400">
            {sitesError.message || "Failed to load sites"}
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
        <>
          <SiteTable
            sites={sitesData?.items || []}
            onEdit={openEditModal}
            onView={openDetailsModal}
            onDelete={handleDelete}
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
        </>
      )}

      <SiteModal
        isOpen={isCreateModalOpen}
        onClose={closeCreateModal}
        onSubmit={handleCreate}
        error={createSiteMutation.error?.message}
        isLoading={createSiteMutation.isPending}
      />

      <SiteModal
        isOpen={isEditModalOpen}
        onClose={closeEditModal}
        site={selectedSite}
        onSubmit={handleUpdate}
        error={updateSiteMutation.error?.message}
        isLoading={updateSiteMutation.isPending}
      />

      <SiteDetailsModal
        isOpen={isDetailsModalOpen}
        onClose={closeDetailsModal}
        site={selectedSite || null}
      />
    </div>
  );
}

export default SitesPage;
