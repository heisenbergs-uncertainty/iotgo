
export interface ApiKey {
  id: number;
  keyId: string;
  name: string;
  token?: string; // Only returned on creation
  isActive: boolean;
  userId: number;
  metadata: string; // JSON string, e.g., {"scopes": ["read", "write"]}
  expiresAt?: string;
}

export interface GenerateApiKeyRequest {
  name: string;
  scopes: string[];
}
