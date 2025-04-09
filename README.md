# TaskFlow Pro (Go Edition)

**TaskFlow Pro** is a team-based task management backend built with **Go**, featuring real-time chat and file uploads. It's perfect for collaborative project management platforms like Trello or Asana — but built from the ground up in Go.

---

## Features

- **User Authentication**
  - JWT-based login & registration
  - Secure password hashing
  - README [User Authentication README](https://github.com/johnson-oragui/TaskFlow-Pro/blob/main/USERAUTHENTICATION.md)
- **Team & Role Management**
  - Create/join teams
  - Role-based permissions (`admin`, `manager`, `member`)
  - README [Team & Role Management README](https://github.com/johnson-oragui/TaskFlow-Pro/blob/main/TEAM-AND-ROLE-MANAGEMENT.md)
- **Task Management**
  - Create/update/delete tasks
  - Assign members, set priority, due dates
  - Track task activity logs
  - README [Task Management README](https://github.com/johnson-oragui/TaskFlow-Pro/blob/main/TASK-MANAGEMENT.md)
- **File Uploads**
  - Attach files to tasks
  - Store locally or to S3-compatible storage
  - README [File Uploads](https://github.com/johnson-oragui/TaskFlow-Pro/blob/main/FILE-UPLOADS.md)
- **Real-time Chat**
  - WebSocket-based chat per team/task
  - Message persistence
- **Dashboard Analytics**
  - Task stats and member activity summaries

---

## Tech Stack

| Layer         | Tech                        |
|---------------|-----------------------------|
| Language      | Go                          |
| Framework     | Gin / Fiber                 |
| ORM           | GORM or SQLC                |
| Database      | PostgreSQL                  |
| Auth          | JWT (`golang-jwt/jwt`)      |
| WebSocket     | Gorilla WebSocket           |
| File Uploads  | Go file handling / S3 SDK   |
| Env Config    | `joho/godotenv`             |
| Testing       | `stretchr/testify`, `httptest` |

---

## Project Structure


taskflow-pro/ ├── cmd/ │ └── server/ # main.go entrypoint ├── config/ # env, config loader ├── controllers/ # HTTP route handlers ├── chat/ # WebSocket logic ├── middleware/ # JWT, logging, etc. ├── models/ # GORM or SQLC models ├── routes/ # Route registration ├── services/ # Business logic ├── uploads/ # Uploaded files ├── utils/ # Reusable helpers ├── tests/ # Unit & integration tests ├── go.mod └── README.md

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/johnson-oragui/TaskFlow-Pro.git
cd taskflow-pro

2. Set Up Environment Variables
Copy the example file and edit it:
cp .env.example .env

Fill in:
DB_URL
JWT_SECRET
UPLOAD_PATH or S3 credentials/cloudinary
3. Install Dependencies
go mod tidy

4. Run the App
go run cmd/server/main.go


API Endpoints (Sample)

WebSocket (Real-time Chat)
Connect to ws://localhost:8000/api/chat/ws/:teamId
JSON Message format:
{
  "sender": "user_id",
  "message": "Hello world",
  "timestamp": "2025-04-09T12:00:00Z"
}


File Uploads
Multipart/form-data request
Stored in /uploads/ (or to S3/cloudinary if enabled)

Testing
go test ./...


Deployment
Use Docker, Railway, Render, or DigitalOcean to deploy the backend. A Dockerfile and docker-compose.yml can be added for easy setup.

License
MIT License — Use it freely!

Author
Your Name
 GitHub: @johnson-oragui

Contributions
Open an issue or submit a PR — ideas & improvements welcome!

---
