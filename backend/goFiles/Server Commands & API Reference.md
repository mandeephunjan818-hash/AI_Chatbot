# Server Commands & API Reference

The HTTP server exposes several endpoints for health checks, testing, and logging.

**Base URL:** `http://localhost:8080`

## System Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Returns service version and status. |
| GET | `/health` | Basic liveness probe (returns 200 OK). |
| GET | `/db-info` | Returns masked database connection details. |

## Testing Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/test-db` | Quick database ping test. |
| GET | `/full-test` | Triggers the internal `RunFullTest()` function. |

## Operational Endpoints

### Log Conversation

**Endpoint:** `POST /test-log`

**Description:** Logs a simulated conversation entry to MongoDB.

```json
{
    "status": "success",
    "message": "Conversation logged to MongoDB",
    "user_id": "user_20231225103000"
}
```

