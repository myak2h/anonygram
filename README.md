# Anonygram

Anonymous image sharing with real-time updates.

## Setup

### Running locally

#### Prerequisites

- Go 1.25+
- Node.js 22+
- npm

**Backend:**
```bash
cd backend
go mod download
go run cmd/server/main.go
```
Server starts at http://localhost:8080

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```
App runs at http://localhost:5173

### Docker

Make sure docker is running and then

```bash
docker compose up --build
```
Frontend: http://localhost:5173, Backend: http://localhost:8080

## Architecture

```
├── backend/
│   ├── cmd/server/          # Entry point
│   └── internal/
│       ├── api/             # Handlers, routing
│       ├── config/          # Environment config
│       ├── models/          # Data structures
│       ├── storage/         # Repository interfaces + in-memory impl
│       └── ws/              # WebSocket hub
└── frontend/
    └── src/
        ├── component/       # React components
        └── hooks/           # useImages, useWebSocket, useTheme
```

### Backend

Go with chi router. Standard `cmd/` + `internal/` layout.

**Key decisions:**

- **Repository interfaces** for storage (`ImageRepository`, `FileRepository`). Currently in-memory, but easy to swap in a database. Also makes mocking in tests straightforward.
- **WebSocket hub** using gorilla/websocket. Hub-and-spoke pattern - single goroutine broadcasts to all clients via channels. Slow clients get dropped rather than blocking others.
- **File validation** via magic bytes, not Content-Type headers. Only accepts PNG/JPEG/GIF.
- **Thread safety** with `sync.RWMutex` for the image store, channels for WebSocket coordination.

### Frontend

React 19, TypeScript, Vite, Tailwind v4.

**Key decisions:**

- **Custom hooks** instead of Redux. `useImages` manages the feed, `useWebSocket` handles connection lifecycle with auto-reconnect.
- **Optimistic updates** - uploaded images appear immediately, deduplication prevents doubles when WebSocket broadcast arrives.
- **New image notifications** - when scrolled down, shows a floating "N new images" banner instead of jumping to top.
- **Dark/light theme** synced with system preference.

## API

### GET /images

Returns all images, newest first.

**Query params:**
- `tag` - Filter by tag (repeatable)

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Sunset at the beach",
    "tags": ["nature", "sunset"],
    "filename": "abc123.jpg",
    "createdAt": "2026-03-27T10:30:00Z"
  }
]
```

### POST /uploads

Upload an image.

**Body:** `multipart/form-data`
- `title` (required) - Image title
- `tags` - Comma-separated tags
- `image` (required) - Image file (PNG, JPEG, GIF)

**Response:** `201 Created` with the new image object

**Errors:**
- `400` - Missing title or image, invalid form data
- `413` - File too large (default limit: 10MB)
- `500` - Storage failure

### GET /files/{filename}

Serves uploaded images.

### GET /ws

WebSocket endpoint. Sends new image objects as JSON when uploads occur.

## Tests

```bash
cd backend && go test ./...
```

Coverage includes handlers, storage, WebSocket hub, and utilities. Uses testify for assertions and mocking.

## Future improvements

Things I'd add with more time:

- **Persistent storage** - SQLite or Postgres instead of in-memory
- **Image processing** - Resize/compress uploads, generate thumbnails
- **Infinite scroll** - Paginate the feed preferably using cursor based 
- **Upload queue** - Process uploads async for better UX under load

## Environment variables

### Backend

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Server port |
| UPLOAD_PATH | ./uploads | Where to store files |
| ALLOWED_ORIGINS | * | CORS origins (comma-separated) |
| MAX_UPLOAD_SIZE | 10485760 | Max file size in bytes (10MB) |
| CLIENT_BUFFER_SIZE | 256 | WebSocket client send buffer size |
| HUB_BUFFER_SIZE | 16 | WebSocket hub channel buffer size |

### Frontend

| Variable | Default | Description |
|----------|---------|-------------|
| VITE_API_BASE_URL | http://localhost:8080 | Backend API URL |
| VITE_WS_URL | ws://localhost:8080/ws | WebSocket endpoint URL |