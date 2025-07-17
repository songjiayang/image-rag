# Image RAG Frontend

Vue.js 3 frontend for the Image RAG service with AI-powered vector search.

## Features

- **Dashboard**: Overview with statistics and recent records
- **Records Management**: Create, view, and manage image records with metadata
- **Image Upload**: Drag-and-drop image upload with preview
- **AI-Powered Search**: Find similar images using vector similarity search
- **Image Gallery**: Carousel and thumbnail views for multiple images
- **Responsive Design**: Works on desktop and mobile devices
- **Settings**: Configure service preferences and check system status

## Tech Stack

- **Vue.js 3** with Composition API
- **TypeScript** for type safety
- **Element Plus** UI components
- **Vue Router** for navigation
- **Axios** for API calls
- **Vite** for build tooling

## Getting Started

### Prerequisites

- Node.js 18+ 
- Backend API running on `localhost:8080`

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Variables

Create a `.env` file in the frontend directory:

```bash
VITE_API_URL=http://localhost:8080/api/v1
```

### API Integration

The frontend connects to the Go backend API at the configured URL. Make sure the backend is running before starting the frontend.

## Project Structure

```
src/
├── components/          # Reusable Vue components
├── views/              # Page components
│   ├── DashboardView.vue
│   ├── RecordsView.vue
│   ├── SearchView.vue
│   ├── RecordDetailView.vue
│   ├── SettingsView.vue
│   └── NotFoundView.vue
├── services/           # API service layer
│   └── api.ts
├── types/             # TypeScript type definitions
│   └── index.ts
├── router/            # Vue Router configuration
│   └── index.ts
├── App.vue            # Root component
├── main.ts            # Application entry point
└── style.css          # Global styles
```

## Usage

1. **Dashboard**: View statistics and recent records
2. **Records**: Manage image collections with metadata
3. **Search**: Upload images to find similar ones using AI
4. **Settings**: Configure preferences and check system status

## Development

```bash
# Run type checking
npm run type-check

# Run linting
npm run lint

# Fix linting issues
npm run lint -- --fix
```