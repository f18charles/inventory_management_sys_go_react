# Backend

Go-based REST API backend for the Inventory Management System.

## Overview

This is the server-side application that handles all business logic and data persistence for the inventory management system. It provides RESTful API endpoints for inventory operations.

## Features

- **RESTful API** - Clean REST API for inventory operations
- **Data Persistence** - Reliable data storage and retrieval
- **Request Validation** - Input validation and error handling
- **CORS Support** - Configured for frontend communication
- **Scalable Architecture** - Built with Go for high performance

## Tech Stack

- Language: Go
- Architecture: REST API
- Data Format: JSON

## Project Structure

```
.
├── main.go              # Application entry point
├── handlers/            # HTTP request handlers
├── models/              # Data models
├── database/            # Database layer
└── go.mod               # Go module definition
```

## API Endpoints

### Items
- `GET /api/items` - Get all inventory items
- `GET /api/items/:id` - Get specific item
- `POST /api/items` - Create new item
- `PUT /api/items/:id` - Update item
- `DELETE /api/items/:id` - Delete item

## Getting Started

### Prerequisites
- Go 1.16 or higher

### Installation

```bash
cd backend
go mod download
```

### Running

```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## Configuration

Configure the server through environment variables or config file:
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - Database connection string
- `CORS_ORIGIN` - Allowed CORS origin (default: http://localhost:3000)

## API Response Format

All responses are in JSON format:

```json
{
  "status": "success",
  "data": {},
  "error": null
}
```

## Testing

```bash
go test ./...
```

## License

MIT
