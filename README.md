

# Chirpy

Chirpy is a lightweight HTTP server written in Go. It provides an API for managing chirps—short text messages—and demonstrates concepts such as HTTP routing, JSON, middleware, authentication, and database-backed storage.

## Features

- Create and retrieve chirps
- Retrieve chirps by ID or author
- User registration and login
- Authentication with access tokens
- Update user credentials
- Health-check endpoint
- JSON request and response handling

## Requirements

- Go (use the version specified in `go.mod`)
- PostgreSQL, if your local configuration uses the database-backed API

## Getting started

Clone the repository and enter its directory:

```bash
git clone <repository-url>
cd <repository-directory>
```
Create a .env file if the project expects one. For example:

```bash
DB_URL=your_database_connection_string
PLATFORM=dev
JWT_SECRET=replace-with-a-long-random-secret
POLKA_KEY=your-webhook-key
```

Use the variable names and values expected by your project. Don’t commit real secrets or credentials.

Run database migrations if your setup requires them, then start the server:

```bash
go run .
```

The server’s address and port depend on your configuration. Check the startup output or project settings.
API

## the routes the server supports:

| Method | Endpoint | Functionality |
|---|---|---|
| `GET` | `/api/healthz` | Checks whether the server is ready. |
| `GET` | `/admin/metrics` | Returns the server’s request metrics. |
| `GET` | `/api/chirps` | Retrieves chirps. |
| `GET` | `/api/chirps/{chirpID}` | Retrieves a specific chirp by its ID. |
| `POST` | `/admin/reset` | Resets the server’s development data or metrics. |
| `POST` | `/api/chirps` | Creates a chirp. |
| `POST` | `/api/users` | Creates a user account. |
| `POST` | `/api/login` | Authenticates a user and returns login credentials. |
| `POST` | `/api/refresh` | Refreshes an access token using a refresh token. |
| `POST` | `/api/revoke` | Revokes a refresh token. |
| `POST` | `/api/polka/webhooks` | Handles a Polka webhook, such as upgrading a user’s account. |
| `PUT` | `/api/users` | Updates the authenticated user’s information. |
| `DELETE` | `/api/chirps/{chirpID}` | Deletes a specific chirp by its ID. |

Some routes may require an Authorization header with a bearer token. Confirm the route names, methods, request bodies, and authentication requirements against your implementation.

## Example request:

Create a chirp, replacing the URL and token with your server’s values:

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"body":"Hello from Chirpy!"}'
```

## Running tests

Run the project’s tests with:

```bash
go test ./...
```

## Project structure

directories:

    main.go — starts and configures the server
    internal/ — application logic and handlers
    sql/ — database queries or migrations
    go.mod — Go module and dependency information

Security notes

    Keep secrets, database credentials, and API keys out of version control.
    Use a strong JWT secret outside local development.
    Don’t expose development-only routes in a production deployment.
