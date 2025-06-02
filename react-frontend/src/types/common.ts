export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  limit: number;
  offset: number;
}

export interface ApiResponse<T> {
  data: T;
  code: number;
}

export interface User {
  id: number;
  email: string;
  role: string;
}

