# Inventory Management System

A full-stack inventory management system built with Go backend and React frontend.

## Overview

This project provides a comprehensive solution for managing inventory operations. It combines a robust Go backend API with a modern React user interface.

## Tech Stack

- **Backend**: Go
- **Frontend**: React
- **Architecture**: Client-Server REST API

## Project Structure

```
.
├── backend/        # Go backend application
└── frontend/       # React frontend application
```

## Backend

Located in `backend/` directory. Go-based REST API for inventory operations.

**Features:**
- RESTful API endpoints
- Data persistence layer
- Request validation
- Error handling

## Frontend

Located in `frontend/` directory. React-based user interface for inventory management.

**Features:**
- Responsive UI components
- Real-time inventory updates
- User-friendly dashboard
- Integration with backend API

## Getting Started

### Backend Setup
```bash
cd backend
# Install and run Go application
go mod download
go run main.go
```

### Frontend Setup
```bash
cd frontend
# Install dependencies
npm install
# Start development server
npm start
```

## Features

- Add, update, and delete inventory items
- Track inventory quantities
- Real-time synchronization between client and server
- Responsive design for mobile and desktop

## License

MIT
