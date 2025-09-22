# 🔖 Bookmark Manager

A fast and lightweight bookmark management API built with Go, featuring PostgreSQL storage and full-text search capabilities.

## ✨ Features

- **Fast API** - Built with Go and Chi router for optimal performance
- **Full-text Search** - PostgreSQL trigram-based search across titles and URLs
- **Export Support** - Export bookmarks in Netscape HTML format
- **Health Monitoring** - Built-in health checks for database connectivity
- **Rate Limiting** - Protection against abuse with configurable limits
- **Docker Ready** - Complete containerization with Docker Compose

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.25+ (for local development)

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/bookmark-manager/bm-backend.git
   cd bookmark-manager
   ```

2. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your preferred settings
   ```

3. **Start with Docker Compose**
   ```bash
   # Using justfile (recommended)
   just up

   # Or directly with Docker Compose
   docker compose up -d
   ```

4. **Verify installation**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

### Local Development

1. **Start database only**
   ```bash
   docker compose up db -d
   ```

2. **Run migrations**
   ```bash
   docker compose run migrate
   ```

3. **Start the application**
   ```bash
   go run cmd/bookmark-manager/main.go
   ```

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/health` | Health check |
| `GET` | `/api/v1/bookmarks` | List bookmarks with pagination & search |
| `POST` | `/api/v1/bookmarks` | Create new bookmark |
| `PATCH` | `/api/v1/bookmarks/{id}` | Update existing bookmark |
| `DELETE` | `/api/v1/bookmarks/{id}` | Delete bookmark |
| `GET` | `/api/v1/bookmarks/exists?url=<url>` | Check if bookmark exists |
| `GET` | `/api/v1/bookmarks/export/html` | Export as Netscape HTML |

### Example Usage

```bash
# Create a bookmark
curl -X POST http://localhost:8080/api/v1/bookmarks \
  -H "Content-Type: application/json" \
  -d '{"title": "GitHub", "url": "https://github.com"}'

# Search bookmarks
curl "http://localhost:8080/api/v1/bookmarks?search=github&per_page=10&page=1"

# Export bookmarks
curl "http://localhost:8080/api/v1/bookmarks/export/html" -o bookmarks.html
```

## 🛠 Available Commands

Using the included `justfile`:

```bash
just up          # Start all services
just down        # Stop all services  
just logs        # View logs
just rebuild     # Rebuild and restart
just fresh-start # Clean restart with new volumes
```

## 🔧 Configuration

Environment variables (see `.env.example`):

- `BM_DB_*` - Database connection settings
- `BM_HTTP_*` - HTTP server configuration  
- `BM_DEBUG` - Enable debug logging
- `BM_NO_COLOR` - Disable colored logs

## 🏗 Architecture

- **Router**: Chi with middleware for logging, CORS, rate limiting
- **Database**: PostgreSQL with trigram search extension
- **Storage**: Clean architecture with interface-based design
- **Validation**: Request validation using go-playground/validator
- **Logging**: Structured logging with slog and tint
