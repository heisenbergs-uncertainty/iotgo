import type { User } from "./common";

export interface LoginResponse {
  user: User;
  token: string;
  error?: string;
}

