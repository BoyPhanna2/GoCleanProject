# Go Clean Architecture Example: Auth + Todo API

This is a complete, runnable Go backend example project using:
- **Framework:** Gin
- **DI:** Uber Fx
- **Database:** SQLite (pure Go, via `modernc.org/sqlite`)
- **Auth:** JWT and bcrypt
- **Architecture:** Feature-first Clean Architecture

## Setup Instructions

1. Ensure you have Go installed (1.20+).
2. Clone or download this repository.
3. Install dependencies (should be automatically handled by `go run` or `go mod tidy`).
4. Start the server:
   ```bash
   go run cmd/api/main.go
   ```
5. The application will automatically create `myapp.db` and start serving on `:8080` (or `PORT` from `.env`).

## Configuration (.env)

You can customize the application behavior by creating a `.env` file in the root directory:

```env
PORT=8080
BASE_URL=http://localhost:8080   # Used dynamically for Swagger UI requests
DB_PATH=myapp.db
JWT_SECRET=supersecretkey
LOG_DIR=logs
```

**Note on Swagger UI:** The application serves interactive API documentation at `/swagger`. The server dynamically injects the `BASE_URL` from your `.env` into the `swagger.yaml` file so that the "Try it out" feature works in any environment (local, staging, or production).

## Docker Deployment

To build and run this application using Docker:

1. **Build the image:**
   ```bash
   docker build -t myapp-api .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 \
     -v $(pwd)/data:/data \
     -v $(pwd)/logs:/logs \
     -e JWT_SECRET=your_super_secret_key \
     myapp-api
   ```
   *Note: We mount `/data` and `/logs` as volumes so your SQLite database and log files persist even if the container stops.*

## Logs

The application implements **daily rotating file logs** via a custom `zap` integration.
- Logs are written to the `logs/` directory.
- A new log file is created every day, named `YYYY-MM-DD.log` (e.g., `2023-10-25.log`).
- **Standard requests (2xx/3xx)** are logged at `INFO` level.
- **Errors (4xx/5xx)** are logged at `ERROR` level and include full context: request body (with passwords sanitized), stack traces, and latency information. This makes debugging much easier by isolating a specific day's traffic and failures.

## API Usage Examples

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com", "password":"password123", "name":"Test User"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com", "password":"password123"}'
```
*Copy the `token` from the response for the next steps.*

### 3. Create a Todo
Replace `<YOUR_TOKEN>` with the token from the login response.
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -d '{"title":"Buy groceries", "description":"Milk, eggs, and bread"}'
```

### 4. List Todos
```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

### 5. Update a Todo
Replace `1` with the actual Todo ID from the list.
```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -d '{"title":"Buy groceries", "description":"Milk, eggs, and bread", "done":true}'
```

### 6. Delete a Todo
Replace `1` with the actual Todo ID.
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```
