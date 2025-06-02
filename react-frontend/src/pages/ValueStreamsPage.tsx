import { useState } from "react";
import { useValueStreamStore } from "../hooks/useValueStreamStore";
import {
  useValueStreams,
  useValueStream,
  useCreateValueStream,
  useUpdateValueStream,
  useDeleteValueStream,
} from "../hooks/useQueries";
import ValueStreamTable from "../components/ValueStreamTable";
import ValueStreamModal from "../components/ValueStreamModal";
import ValueStreamDetailsModal from "../components/ValueStreamDetailsModal";
import type { ValueStream } from "../types/valueStream";

function ValueStreamsPage() {
  const {
    isCreateModalOpen,
    isEditModalOpen,
    isDetailsModalOpen,
    selectedValueStreamId,
    openCreateModal,
    closeCreateModal,
    openEditModal,
    closeEditModal,
    openDetailsModal,
    closeDetailsModal,
  } = useValueStreamStore();
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [nameFilter, setNameFilter] = useState("");
  // const [sort, setSort] = useState("name");

  const {
    data: valueStreamsData,
    isLoading: isValueStreamsLoading,
    error: valueStreamsError,
  } = useValueStreams({
    limit,
    offset: (page - 1) * limit,
    name: nameFilter,
    sort: "name",
  });
  const { data: selectedValueStream, isLoading: isValueStreamLoading } =
    useValueStream(selectedValueStreamId || 0);

  const createValueStreamMutation = useCreateValueStream();
  const updateValueStreamMutation = useUpdateValueStream();
  const deleteValueStreamMutation = useDeleteValueStream();

  const handleCreate = (valueStream: Partial<ValueStream>) => {
    createValueStreamMutation.mutate(valueStream, {
      onSuccess: () => closeCreateModal(),
      onError: (err: any) => console.error("Create error:", err),
    });
  };

  const handleUpdate = (valueStream: Partial<ValueStream>) => {
    if (selectedValueStreamId) {
      updateValueStreamMutation.mutate(
        { id: selectedValueStreamId, valueStream },
        {
          onSuccess: () => closeEditModal(),
          onError: (err: any) => console.error("Update error:", err),
        },
      );
    }
  };

  const handleDelete = (id: number) => {
    if (window.confirm("Are you sure you want to delete this value stream?")) {
      deleteValueStreamMutation.mutate(id, {
        onError: (err: any) => console.error("Delete error:", err),
      });
    }
  };

  // const handleSort = (newSort: string) => {
  //   setSort(newSort);
  // };

  const totalPages = valueStreamsData
    ? Math.ceil(valueStreamsData.total / limit)
    : 1;
  const isLoading = isValueStreamsLoading || isValueStreamLoading;

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold text-secondary-900 dark:text-secondary-50">
          Value Stream Management
        </h2>
        <button
          onClick={openCreateModal}
          className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
        >
          Add Value Stream
        </button>
      </div>

      {valueStreamsError && (
        <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900/30">
          <p className="text-sm text-red-600 dark:text-red-400">
            {valueStreamsError.message || "Failed to load value streams"}
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
          <ValueStreamTable
            valueStreams={valueStreamsData?.items || []}
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

      <ValueStreamModal
        isOpen={isCreateModalOpen}
        onClose={closeCreateModal}
        onSubmit={handleCreate}
        error={createValueStreamMutation.error?.message}
        isLoading={createValueStreamMutation.isPending}
      />

      <ValueStreamModal
        isOpen={isEditModalOpen}
        onClose={closeEditModal}
        valueStream={selectedValueStream}
        onSubmit={handleUpdate}
        error={updateValueStreamMutation.error?.message}
        isLoading={updateValueStreamMutation.isPending}
      />

      <ValueStreamDetailsModal
        isOpen={isDetailsModalOpen}
        onClose={closeDetailsModal}
        valueStream={selectedValueStream || null}
      />
    </div>
  );
}

export default ValueStreamsPage;
