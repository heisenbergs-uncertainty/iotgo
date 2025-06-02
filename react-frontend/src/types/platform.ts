export const ApiAuthMethod = {
  NONE: "none",
  API_KEY: "api_key",
  BEARER: "bearer",
  BASIC: "basic",
};

export type ApiAuthMethodType = keyof typeof ApiAuthMethod;


export const PlatformTypes = {
  REST: "REST",
  InfluxDB: "InfluxDB",
};

export type PlatformType = keyof typeof PlatformTypes;

export const ResourceTypes = {
  REST: "rest_endpoint",
  InfluxDB: "influxdb_query",
} as const;

export type ResourceType = keyof typeof ResourceTypes;

export interface Platform {
  id: number;
  name: string;
  type: PlatformType
  metadata: string; // JSON string, parsed into RESTMetadata or InfluxDBMetadata
}

export interface RESTMetadata {
  baseEndpoint: string;
  timeout: number;
  auth: RESTAuth;
}

export interface RESTAuth {
  type: ApiAuthMethodType;
  apiKey?: string;
  bearerToken?: string;
  basicAuth?: BasicAuth;
}

export interface BasicAuth {
  username: string;
  password: string;
}

export interface InfluxDBMetadata {
  url: string;
  token: string;
  org: string;
  bucket: string;
  timeout: number;
}

export interface Resource {
  id: number;
  platformId: number;
  name: string;
  type: ResourceType;
  details: string; // JSON string, parsed into RESTResourceDetails or InfluxDBResourceDetails
}

export interface RESTResourceDetails {
  method: "GET" | "POST" | "PUT" | "DELETE";
  path: string;
  headers?: Record<string, string>;
  queryParams?: Record<string, string>;
  body?: string;
}

export interface InfluxDBResourceDetails {
  bucket: string;
  measurement: string;
  field: string;
  timeRange: string;
}

export interface DevicePlatform {
  platformId: number;
  deviceId: number;
  deviceAlias: string;
}

export interface TestConnectionInput {
  type: "REST" | "InfluxDB";
  metadata: string;
}

export interface FetchDeviceDataResponse {
  device_id: number;
  platform_id: number;
  alias: string;
  data: Record<string, any>;
}
