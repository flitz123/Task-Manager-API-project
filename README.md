# Task Manager

A responsive task manager with a Go REST API and dependency-free frontend.

## Local development

```bash
go run ./cmd/server
```

Open `http://localhost:8080`.

## API

- `GET /tasks`
- `POST /tasks`
- `GET /tasks/{id}`
- `PUT /tasks/{id}`
- `PATCH /tasks/{id}`
- `DELETE /tasks/{id}`
- `GET /health`

Task JSON fields include `title`, `description`, `completed`, `priority`, `category`, and `dueDate`.

## Vercel

Import the repository into Vercel. The included `vercel.json` routes `/tasks` and `/health` to the Go serverless function, while `public` is served as the frontend.
