import {
  AuthResponse,
  Device,
  ErrorLog,
  ErrorResponse,
  Integration,
  Recording,
  Snapshot,
  User,
} from "@/types/api";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

async function fetchWithError<T>(
  url: string,
  options: RequestInit = {},
): Promise<T> {
  const headers = {
    "Content-Type": "application/json",
    ...options.headers,
  };
  console.log(
    `API request: ${url}, method: ${options.method || "GET"}, headers:`,
    headers,
  );

  // const response = await fetch(`${API_BASE_URL}${url}`, {
  //   ...options,
  //   credentials: "include",
  //   headers,
  // });
  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
    credentials: "include",
  });

  const responseHeaders = Object.fromEntries(response.headers.entries());
  console.log(
    `API response: ${url}, status: ${response.status}, headers:`,
    responseHeaders,
  );

  // Log all cookies from the response
  const cookies = responseHeaders["set-cookie"];
  if (cookies) {
    console.log("Received cookies:", cookies);
  }

  if (!response.ok) {
    const errorData: ErrorResponse = await response.json().catch(() => ({}));
    throw new Error(
      errorData.error || `HTTP error! status: ${response.status}`,
    );
  }

  return response.json() as Promise<T>;
}

export const api = {
  // Auth endpoints
  login: (username: string, password: string) =>
    fetchWithError<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    }),

  getMe: () =>
    fetchWithError<AuthResponse>("/auth/me", {
      method: "GET",
    }),

  logout: () =>
    fetchWithError<{ message: string }>("/auth/logout", {
      method: "GET",
    }),

  // Device endpoints
  getDevices: (
    params: { limit?: number; offset?: number; query?: string } = {},
  ) => {
    const url = new URLSearchParams();
    if (params.limit) url.append("limit", params.limit.toString());
    if (params.offset) url.append("offset", params.offset.toString());
    if (params.query) url.append("query", params.query);
    return fetchWithError<Device[]>("/devices");
  },

  getDevice: (id: number) => fetchWithError<Device>(`/devices/${id}`),

  createDevice: (device: Omit<Device, "Id">) =>
    fetchWithError<Device>("/devices", {
      method: "POST",
      body: JSON.stringify(device),
    }),

  updateDevice: (id: number, device: Partial<Device>) =>
    fetchWithError<string>(`/devices/${id}`, {
      method: "PUT",
      body: JSON.stringify(device),
    }),

  deleteDevice: (id: number) =>
    fetchWithError<string>(`/devices/${id}`, {
      method: "DELETE",
    }),

  // User endpoints
  getUsers: (params: { limit?: number; offset?: number } = {}) => {
    const url = new URLSearchParams();
    if (params.limit) url.append("limit", params.limit.toString());
    if (params.offset) url.append("offset", params.offset.toString());
    return fetchWithError<User[]>(`/users?${url.toString()}`);
  },

  getUser: (id: number) => fetchWithError<User>(`/users/${id}`),

  createUser: (user: Omit<User, "Id"> & { Password: string }) =>
    fetchWithError<User>("/users", {
      method: "POST",
      body: JSON.stringify(user),
    }),

  updateUser: (id: number, user: Partial<User> & { Password?: string }) =>
    fetchWithError<string>(`/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(user),
    }),

  deleteUser: (id: number) =>
    fetchWithError<string>(`/users/${id}`, {
      method: "DELETE",
    }),

  // Integration endpoints
  getIntegrations: (params: { limit?: number; offset?: number } = {}) => {
    const url = new URLSearchParams();
    if (params.limit) url.append("limit", params.limit.toString());
    if (params.offset) url.append("offset", params.offset.toString());
    return fetchWithError<Integration[]>(`/integration?${url.toString()}`);
  },

  // Snapshot endpoints
  getSnapshots: (
    params: { limit?: number; offset?: number; query?: string } = {},
  ) => {
    const url = new URLSearchParams();
    if (params.limit) url.append("limit", params.limit.toString());
    if (params.offset) url.append("offset", params.offset.toString());
    if (params.query) url.append("query", params.query);
    return fetchWithError<Snapshot[]>(`/snapshot?${url.toString()}`);
  },

  // Recording endpoints
  getRecordings: (
    params: { limit?: number; offset?: number; query?: string } = {},
  ) => {
    const url = new URLSearchParams();
    if (params.limit) url.append("limit", params.limit.toString());
    if (params.offset) url.append("offset", params.offset.toString());
    if (params.query) url.append("query", params.query);
    return fetchWithError<Recording[]>(`/recording?${url.toString()}`);
  },

  // ErrorLog endpoints
  getErrorLogs: (params: { limit?: number; offset?: number } = {}) => {
    const url = new URLSearchParams();
    if (params.limit) url.append("limit", params.limit.toString());
    if (params.offset) url.append("offset", params.offset.toString());
    return fetchWithError<ErrorLog[]>(`/error_log?${url.toString()}`);
  },
};

