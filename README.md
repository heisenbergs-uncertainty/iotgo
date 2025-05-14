# IoTGo - IoT API Gateway

IoTGo is a comprehensive IoT API gateway designed to manage and integrate multiple IoT devices across various platforms. It provides a unified interface for data access and device management through a web UI and REST API.

## Table of Contents

1. [Overview](#overview)
2. [Key Features](#key-features)
3. [Architecture](#architecture)
4. [Models](#models)
5. [Drivers](#drivers)
6. [API Endpoints](#api-endpoints)
7. [Web Interface](#web-interface)
8. [Setup and Installation](#setup-and-installation)
9. [Development](#development)
10. [Docker Deployment](#docker-deployment)
11. [Authentication](#authentication)
12. [License](#license)

## Overview

IoTGo acts as a gateway between your applications and various IoT platforms. It abstracts the complexities of different communication protocols and APIs, providing a consistent interface to interact with IoT devices across multiple systems.

The application allows you to:
- Register and manage IoT devices
- Connect devices to different platforms (REST APIs, OPC UA servers, custom SDKs)
- Associate devices with physical sites and value streams
- Fetch device data through a unified API
- Manage access through API keys

## Key Features

- **Multi-Platform Support**: Connect to REST APIs, OPC UA servers, and platform-specific SDKs
- **Device Management**: Register, organize, and monitor IoT devices
- **Site Organization**: Group devices by physical location
- **Value Streams**: Categorize devices by production flow or business function
- **API Key Authentication**: Secure API access with revocable keys
- **Web Dashboard**: User-friendly interface for platform management
- **Extensible Driver System**: Easily add support for new platform types

## Architecture

IoTGo follows a modular architecture built on the Beego web framework and GORM ORM:

- **Controllers**: Handle HTTP requests and implement business logic
- **Models**: Define data structures and relationships
- **Drivers**: Implement platform-specific communication protocols
- **Middleware**: Handle authentication and request filtering
- **Routers**: Define API endpoints and web routes
- **DAL (Data Access Layer)**: Manage database interactions

## Models

### Core Entities

1. **Device**: Represents an IoT device with metadata
   - Can be associated with multiple platforms
   - Can belong to a site and value stream

2. **Platform**: Represents an external system (REST API, OPC UA server, SDK-based system)
   - Has a type that determines which driver to use
   - Contains connection metadata (endpoints, credentials)
   - Can have multiple resources (endpoints, nodes)

3. **DevicePlatform**: Many-to-many relationship between devices and platforms
   - Contains a unique `deviceAlias` that identifies the device in the platform

4. **Site**: Represents a physical location where devices are installed

5. **ValueStream**: Represents a production line or business function

6. **Resource**: Represents a specific interaction point for a platform
   - For REST: endpoints, methods, parameters
   - For OPC UA: node IDs
   - For SDK: method names, parameters

7. **User**: System user with authentication

8. **ApiKey**: API access tokens for authentication

## Drivers

The system uses a driver interface to abstract communication with different platform types:

```go
type PlatformDriver interface {
    Connect(ctx context.Context) error
    FetchData(ctx context.Context, deviceAlias string) (interface{}, error)
    Disconnect(ctx context.Context) error
}
```

### Implemented Drivers

1. **RESTDriver**: For platforms with REST APIs
   - Configurable base URL, authentication, and timeout
   - Sends HTTP requests to specified endpoints
   - Supports various authentication methods (API key, bearer token)

2. **OPCUADriver**: For OPC UA servers
   - Connects to OPC UA endpoints
   - Reads values from specified nodes
   - Manages connection lifecycle

3. **SDKDriver**: For platforms with proprietary SDKs
   - Template for implementing SDK-specific logic
   - Can be extended for specific platform SDKs

## API Endpoints

### Device Management
- `GET /api/devices`: List all devices with filtering and pagination
- `POST /api/devices`: Create a new device
- `GET /api/devices/:id`: Get a specific device
- `PUT /api/devices/:id`: Update a device
- `DELETE /api/devices/:id`: Delete a device

### Platform Management
- `GET /api/platforms`: List all platforms
- `POST /api/platforms`: Create a new platform
- `GET /api/platforms/:id`: Get a specific platform
- `PUT /api/platforms/:id`: Update a platform
- `DELETE /api/platforms/:id`: Delete a platform

### Device-Platform Associations
- `GET /api/devices/:device_id/platforms`: List platforms associated with a device
- `POST /api/devices/:device_id/platforms`: Associate a device with a platform
- `DELETE /api/devices/:device_id/platforms/:platform_id`: Remove association

### Data Access
- `GET /api/platforms/:platform_id/devices/:device_id/data`: Fetch device data from a platform

### Site Management
- `GET /api/sites`: List all sites
- `POST /api/sites`: Create a new site
- `GET /api/sites/:id`: Get a specific site
- `PUT /api/sites/:id`: Update a site
- `DELETE /api/sites/:id`: Delete a site

### Value Stream Management
- `GET /api/value_streams`: List all value streams
- `POST /api/value_streams`: Create a new value stream
- `GET /api/value_streams/:id`: Get a specific value stream
- `PUT /api/value_streams/:id`: Update a value stream
- `DELETE /api/value_streams/:id`: Delete a value stream

## Web Interface

IoTGo provides a web dashboard for managing the system without using the API directly:

- `/login`: User authentication
- `/dashboard`: System overview and statistics
- `/devices`: Device management
- `/devices/:device_id/platforms`: Device-platform association management
- `/platforms`: Platform management
- `/sites`: Site management
- `/value_streams`: Value stream management
- `/api_key`: API key management

## Setup and Installation

### Prerequisites
- Go 1.16+
- PostgreSQL 12+
- Node.js and npm (for frontend assets)

### Database Setup
```bash
# Create database
createdb iotgo

# Configure connection in main.go or environment variables
```

### Build and Run
```bash
# Get dependencies
go mod download

# Build the application
go build -o iotgo

# Run the application
./iotgo
```

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: HTTP port (default: 8080)
- `SESSION_SECRET`: Secret for session encryption

## Development

### Code Generation
The project uses code generation for database models:

```bash
# Generate database models
go generate ./generator/generate.go
```

### CSS with Tailwind
Frontend styling uses Tailwind CSS:

```bash
# Install dependencies
npm install

# Build CSS
npm run build-css
```

## Docker Deployment

A Docker configuration is provided for easy deployment:

```bash
# Build Docker image
docker build -t iotgo .

# Run with Docker Compose
docker-compose up -d
```

## Authentication

IoTGo uses two authentication methods:

1. **Web UI**: Session-based authentication with username/password
2. **API**: API key authentication with tokens

API keys can be generated and managed in the web interface.

## License

MIT License