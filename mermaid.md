```mermaid
erDiagram
    User {
        int ID PK
        string Email
        string Password
        string Name
        string Role
        string Avatar
        string Bio
        datetime LastLogin
        bool TwoFactorEnabled
        json Preferences
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    Site {
        int ID PK
        string Name
        string Description
        string Address
        string City
        string State
        string Country
        float Latitude
        float Longitude
        bool IsActive
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    Device {
        int ID PK
        string Name
        string Status
        int SiteID FK
        int ValueStreamID FK
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    ValueStream {
        int ID PK
        string Name
        string Description
        string Type
        bool IsActive
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    Platform {
        int ID PK
        string Name
        string Type
        string Direction
        string Endpoint
        string EndpointURL
        string Description
        string AuthMethod
        string AuthType
        string ConnectionState
        string ConnectionError
        datetime LastConnected
        int OrganizationID
        bool IsActive
        string Metadata
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    PlatformController {
        int ID PK
        string Name
        string Type
        int PlatformID FK
        string Version
        string Status
        datetime LastHeartbeat
        bool IsActive
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    SensorData {
        int ID PK
        int DeviceID FK
        datetime Timestamp
        float Value
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    ApiKey {
        int ID PK
        string Name
        string KeyID
        string KeyHash
        bool IsActive
        int OwnerID
        int UserID
        datetime LastUsed
        datetime ExpiresAt
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    UserInteraction {
        int ID PK
        int UserID FK
        string Type
        string Details
        datetime CreatedAt
        datetime UpdatedAt
        datetime DeletedAt
    }

    Site ||--o{ Device : "has"
    ValueStream ||--o{ Device : "has"
    Device ||--o{ SensorData : "generates"
    Platform ||--o{ PlatformController : "has"
    User ||--o{ UserInteraction : "performs"
    User ||--o{ ApiKey : "owns"
    Device }o--o{ Platform : "connects_to"
```