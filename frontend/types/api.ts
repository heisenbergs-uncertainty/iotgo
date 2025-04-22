export interface Activity {
    Description: string;
    Time: string;
}


export interface User {
    Id: number;
    Username: string;
    Roles: any[]; // Changed from Role (string) to Roles (array)
}

export interface Device {
    Id: number;
    Name: string;
    Manufacturer: string;
    Type: string;
    Building: string;
}

export interface Integration {
    Id: number;
    DeviceId: number;
    IntegrationType: string;
    Identifier: string;
    Host: string;
    Port: string;
    Protocol: string;
}

export interface Snapshot {
    Id: number;
    DeviceId: number;
    IntegrationId: number;
    Timestamp: string;
    Nodes: string;
}

export interface Recording {
    Id: number;
    DeviceId: number;
    IntegrationId: number;
    NodeId: string;
    NodeName: string;
    Timestamp: string;
    Value: string;
}

export interface ErrorLog {
    Id: number;
    Message: string;
    Timestamp: string;
    UserId: number;
}

export interface AuthResponse {
    username: string;
    roles: any[]; // Changed from role (string) to roles (array)
    user_id: number;
}


export interface ErrorResponse {
    error: string;
}