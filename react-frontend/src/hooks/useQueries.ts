import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "./useAuthStore";
import type { PaginatedResponse, ApiResponse, User } from "../types/common";
import type { Device } from "../types/device";
import type { Site } from "../types/site";
import type { ValueStream } from "../types/valueStream";
import type { ApiKey, GenerateApiKeyRequest } from "../types/apikey";
import type {
  Platform,
  Resource,
  TestConnectionInput,
  FetchDeviceDataResponse,
  ResourceType,
} from "../types/platform";

interface FetchParams {
  limit?: number;
  offset?: number;
  name?: string;
  sort?: string;
}

const getAuthHeaders = () => {
  const token = localStorage.getItem("token");
  return token ? { Authorization: `Bearer ${token}` } : {};
};

const handleAuthError = (
  err: any,
  logout: () => void,
  navigate: () => void,
) => {
  if (err.response?.status === 401) {
    localStorage.removeItem("token");
    localStorage.removeItem("user_id");
    localStorage.removeItem("user_role");
    logout();
    navigate();
    return new Error("Unauthorized: Please log in again");
  }
  return err;
};

// Device Queries
export const useDevices = (params: FetchParams) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["devices", params],
    queryFn: async () => {
      try {
        const response = await axios.get<
          ApiResponse<PaginatedResponse<Device>>
        >("/api/devices", {
          params,
          headers: getAuthHeaders(),
        });
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
  });
};

export const useDevice = (id: number) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["device", id],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<Device>>(
          `/api/devices/${id}`,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!id,
  });
};

export const useCreateDevice = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (device: Partial<Device>) => {
      try {
        const response = await axios.post<ApiResponse<Device>>(
          "/api/devices",
          device,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
    },
  });
};

export const useUpdateDevice = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      id,
      device,
    }: {
      id: number;
      device: Partial<Device>;
    }) => {
      try {
        const response = await axios.put<ApiResponse<Device>>(
          `/api/devices/${id}`,
          device,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
      queryClient.invalidateQueries({ queryKey: ["device", data.id] });
    },
  });
};

export const useDeleteDevice = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (id: number) => {
      try {
        await axios.delete(`/api/devices/${id}`, {
          headers: getAuthHeaders(),
        });
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["devices"] });
    },
  });
};

// Site Queries
export const useSites = (params: FetchParams = {}) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["sites", params],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<PaginatedResponse<Site>>>(
          "/api/sites",
          {
            params,
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
  });
};

export const useSite = (id: number) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["site", id],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<Site>>(
          `/api/sites/${id}`,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!id,
  });
};

export const useCreateSite = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (site: Partial<Site>) => {
      try {
        const response = await axios.post<ApiResponse<Site>>(
          "/api/sites",
          site,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sites"] });
    },
  });
};

export const useUpdateSite = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({ id, site }: { id: number; site: Partial<Site> }) => {
      try {
        const response = await axios.put<ApiResponse<Site>>(
          `/api/sites/${id}`,
          site,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["sites"] });
      queryClient.invalidateQueries({ queryKey: ["site", data.id] });
    },
  });
};

export const useDeleteSite = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (id: number) => {
      try {
        await axios.delete(`/api/sites/${id}`, {
          headers: getAuthHeaders(),
        });
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sites"] });
    },
  });
};

// Value Stream Queries
export const useValueStreams = (params: FetchParams = {}) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["valueStreams", params],
    queryFn: async () => {
      try {
        const response = await axios.get<
          ApiResponse<PaginatedResponse<ValueStream>>
        >("/api/value-streams", {
          params,
          headers: getAuthHeaders(),
        });
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
  });
};

export const useValueStream = (id: number) => {
  const {logout} = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["valueStream", id],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<ValueStream>>(
          `/api/value-streams/${id}`,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!id,
  });
};

export const useCreateValueStream = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (valueStream: Partial<ValueStream>) => {
      try {
        const response = await axios.post<ApiResponse<ValueStream>>(
          "/api/value-streams",
          valueStream,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["valueStreams"] });
    },
  });
};

export const useUpdateValueStream = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      id,
      valueStream,
    }: {
      id: number;
      valueStream: Partial<ValueStream>;
    }) => {
      try {
        const response = await axios.put<ApiResponse<ValueStream>>(
          `/api/value-streams/${id}`,
          valueStream,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["valueStreams"] });
      queryClient.invalidateQueries({ queryKey: ["valueStream", data.id] });
    },
  });
};

export const useDeleteValueStream = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (id: number) => {
      try {
        await axios.delete(`/api/value-streams/${id}`, {
          headers: getAuthHeaders(),
        });
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["valueStreams"] });
    },
  });
};

// API Key Queries
export const useApiKeys = (userId: string) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["apiKeys", userId],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<ApiKey[]>>(
          `/api/users/${userId}/apikeys`,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!userId,
  });
};

export const useGenerateApiKey = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (request: GenerateApiKeyRequest) => {
      try {
        const response = await axios.post<ApiResponse<{ token: string }>>(
          "/api/apikeys",
          request,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["apiKeys"] });
    },
  });
};

export const useRevokeApiKey = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (keyId: string) => {
      try {
        await axios.delete(`/api/apikeys/${keyId}`, {
          headers: getAuthHeaders(),
        });
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["apiKeys"] });
    },
  });
};

// User Queries
export const useUser = (id: number) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["user", id],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<User>>(
          `/api/users?id=${id}`,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!id,
  });
};

// Platform Queries
export const usePlatforms = (params: FetchParams) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["platforms", params],
    queryFn: async () => {
      try {
        const response = await axios.get<
          ApiResponse<PaginatedResponse<Platform>>
        >("/api/platforms", {
          params,
          headers: getAuthHeaders(),
        });
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
  });
};

export const usePlatform = (id: number) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["platform", id],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<Platform>>(
          `/api/platforms/${id}`,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!id,
  });
};

export const useCreatePlatform = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (platform: Partial<Platform>) => {
      try {
        const response = await axios.post<ApiResponse<Platform>>(
          "/api/platforms",
          platform,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platforms"] });
    },
  });
};

export const useUpdatePlatform = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      id,
      platform,
    }: {
      id: number;
      platform: Partial<Platform>;
    }) => {
      try {
        const response = await axios.put<ApiResponse<Platform>>(
          `/api/platforms/${id}`,
          platform,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["platforms"] });
      queryClient.invalidateQueries({ queryKey: ["platform", data.id] });
    },
  });
};

export const useDeletePlatform = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (id: number) => {
      try {
        await axios.delete(`/api/platforms/${id}`, {
          headers: getAuthHeaders(),
        });
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["platforms"] });
    },
  });
};

export const useTestPlatformConnection = () => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (input: TestConnectionInput) => {
      try {
        const response = await axios.post<ApiResponse<{ message: string }>>(
          "/api/platforms/test",
          input,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
  });
};

export const useFetchDeviceData = () => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      platformId,
      deviceId,
    }: {
      platformId: number;
      deviceId: number;
    }) => {
      try {
        const response = await axios.get<ApiResponse<FetchDeviceDataResponse>>(
          `/api/platforms/${platformId}/devices/${deviceId}/data`,
          { headers: getAuthHeaders() },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
  });
};

// Resource Queries
export const useResources = (platformId: number, params: FetchParams) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["resources", platformId, params],
    queryFn: async () => {
      try {
        const response = await axios.get<
          ApiResponse<PaginatedResponse<Resource>>
        >(`/api/platforms/${platformId}/resources`, {
          params,
          headers: getAuthHeaders(),
        });
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!platformId,
  });
};

export const useResource = (id: number) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["resource", id],
    queryFn: async () => {
      try {
        const response = await axios.get<ApiResponse<Resource>>(
          `/api/platforms/resources/${id}`,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!id,
  });
};

export const useResourceForEdit = (id: number) => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useQuery({
    queryKey: ["resourceEdit", id],
    queryFn: async () => {
      try {
        const response = await axios.get<
          ApiResponse<{ id: number; name: string; type: ResourceType; details: any }>
        >(`/api/resources/${id}/edit`, { headers: getAuthHeaders() });
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    enabled: !!id,
  });
};

export const useCreateResource = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      platformId,
      resource,
    }: {
      platformId: number;
      resource: Partial<Resource>;
    }) => {
      try {
        const response = await axios.post<ApiResponse<Resource>>(
          `/api/platforms/${platformId}/resources`,
          resource,
          { headers: getAuthHeaders() },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (_, { platformId }) => {
      queryClient.invalidateQueries({ queryKey: ["resources", platformId] });
    },
  });
};

export const useBulkCreateResources = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      platformId,
      resources,
    }: {
      platformId: number;
      resources: Partial<Resource>[];
    }) => {
      try {
        const response = await axios.post<
          ApiResponse<{ created: Resource[]; errors?: string[] }>
        >(`/api/platforms/${platformId}/resources/bulk`, resources, {
          headers: getAuthHeaders(),
        });
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (_, { platformId }) => {
      queryClient.invalidateQueries({ queryKey: ["resources", platformId] });
    },
  });
};

export const useUpdateResource = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      id,
      resource,
    }: {
      id: number;
      resource: Partial<Resource>;
    }) => {
      try {
        const response = await axios.put<ApiResponse<Resource>>(
          `/api/platforms/resources/${id}`,
          resource,
          {
            headers: getAuthHeaders(),
          },
        );
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: ["resources", data.platformId],
      });
      queryClient.invalidateQueries({ queryKey: ["resource", data.id] });
    },
  });
};

export const useDeleteResource = () => {
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({
      id,
      platformId,
    }: {
      id: number;
      platformId: number;
    }) => {
      try {
        await axios.delete(`/api/platforms/${platformId}/resources/${id}`, {
          headers: getAuthHeaders(),
        });
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
    onSuccess: (_, { platformId }) => {
      queryClient.invalidateQueries({ queryKey: ["resources", platformId] });
    },
  });
};

export const useTestResource = () => {
  const { logout } = useAuthStore();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async (id: number) => {
      try {
        const response = await axios.post<
          ApiResponse<{
            resource_id: number;
            name: string;
            type: string;
            platform_id: number;
            platform_type: string;
            result: any;
          }>
        >(`/api/resources/${id}/test`, {}, { headers: getAuthHeaders() });
        return response.data.data;
      } catch (err: any) {
        throw handleAuthError(err, logout, () => navigate("/login"));
      }
    },
  });
};
