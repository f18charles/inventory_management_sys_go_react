# Frontend

React-based user interface for the Inventory Management System.

## Overview

Modern, responsive web application built with React that provides an intuitive interface for managing inventory items. Features real-time updates and seamless integration with the backend API.

## Features

- **Responsive Design** - Works on desktop and mobile devices
- **Real-time Updates** - Instant synchronization with backend
- **User-Friendly Interface** - Intuitive dashboard and controls
- **Component-Based Architecture** - Reusable React components
- **State Management** - Efficient state management solution

## Tech Stack

- Language: JavaScript/JSX
- Framework: React
- Build Tool: Create React App
- HTTP Client: Fetch API or Axios

## Project Structure

```
.
├── src/
│   ├── components/          # React components
│   ├── pages/               # Page components
│   ├── services/            # API service layer
│   ├── App.jsx              # Root component
│   └── index.js             # Entry point
├── public/                  # Static files
├── package.json             # Dependencies
└── README.md                # This file
```

## Available Components

- **Dashboard** - Main inventory overview
- **ItemList** - Display and manage inventory items
- **ItemForm** - Create and edit items
- **SearchBar** - Search and filter items
- **Navigation** - App navigation

## Getting Started

### Prerequisites
- Node.js 14+ and npm/yarn

### Installation

```bash
cd frontend
npm install
```

### Development Server

```bash
npm start
```

Opens [http://localhost:3000](http://localhost:3000) in the browser.

### Build

```bash
npm run build
```

Creates a production-ready build in the `build/` directory.

## API Integration

The frontend communicates with the backend API at `http://localhost:8080/api`.

### Example Service Call

```javascript
// services/inventoryService.js
const getItems = async () => {
  const response = await fetch('http://localhost:8080/api/items');
  return response.json();
};
```

## Environment Variables

Create a `.env` file in the frontend directory:

```
REACT_APP_API_URL=http://localhost:8080/api
```

## Testing

```bash
npm test
```

## Available Scripts

- `npm start` - Run development server
- `npm run build` - Build for production
- `npm test` - Run tests
- `npm run eject` - Eject from Create React App (irreversible)

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## License

MIT
