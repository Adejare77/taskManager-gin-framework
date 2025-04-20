# Task Manager API

[![Go Report Card](https://goreportcard.com/badge/github.com/Adejare77/taskManager-gin-framework)](https://goreportcard.com/report/github.com/Adejare77/taskManager-gin-framework)

A robust task management system with user authentication, session management, and automatic task status updates.

## Key Features
- **User Authentication**: Session-based authentication using Redis for server-side session storage

- **Task Management**:
  - Create tasks with start/due dates
  - Automatic status transitions (pending → in-progress → overdue/completed)
  - Paginated task listing
- **Scheduled Updates**: Cron job for automatic status updates
- **RESTful API**: Full CRUD operations for tasks
- **Database**: PostgreSQL with connection pooling
- **Security**: Password hashing, session management, and validation

## Date Handling
### Task Creation Rules
- `due_date`: **Required** (format: `YYYY-MM-DD HH:MM`)
- `start_date`: **Optional** (defaults to current time if omitted)

Example JSON:
```json
{
  "start_date": "2025-03-15 09:00",
  "due_date": "2025-03-20 17:00"
}
```

## API Endpoints

### Authentication
| Method | Endpoint     | Description       |
|--------|--------------|-------------------|
| POST   | /register    | User registration |
| POST   | /login       | User login        |
| GET    | /user/logout | User logout       |
| DELETE | /user        | Delete user       |

### Task Management
| Method | Endpoint          | Description        |
|--------|-------------------|--------------------|
| POST   | /tasks            | Create new task    |
| GET    | /tasks            | List all tasks     |
| GET    | /tasks/:task_id   | Get task details   |
| PATCH  | /tasks/:task_id   | Update task        |
| DELETE | /tasks/:task_id   | Delete task        |

### System
| Method | Endpoint  | Description  |
|--------|-----------|--------------|
| GET    | /health   | Health check |

## Setup Instructions

### Prerequisites
- Go 1.24+
- PostgreSQL
- Redis

### Environment Variables
Create a `.env` file:

```env
# Cron Job Scheduler
CRON_SCHEDULE="60"

# Database
DB_HOST="localhost"
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=task_manager
DB_PORT=5432
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME="30m"

# Session
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=""
SECRET_KEY=your_secret_key
REDIS_SIZE=10
SESSION_MAX_AGE=600


# Server
SERVER_PORT=3000
```

### Installation
```bash
# clone task manager
git clone https://github.com/adejare77/taskManager-gin-framework
cd taskManager

# Set up environment
cp envsample .env
nano .env  # Update with your credentials

# Install dependencies
go mod download

# start server
go run cmd/cmd/main.go
```

## Example Requests

### Create Task (with default start_date)
```bash
curl -X POST http://localhost:3000/tasks \
  -H "Content-Type: application/json" \
  -H "Cookie: taskManager=<your_session_cookie>" \
  -d '{
    "title": "Project Deadline",
    "description": "Finalize project deliverables",
    "due_date": "2025-03-20 17:00"
  }'
```

### Create Task (with custom start_date)
```bash
curl -X POST http://localhost:3000/tasks \
  -H "Content-Type: application/json" \
  -H "Cookie: taskManager=<your_session_cookie>" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly sync meeting",
    "start_date": "2025-03-15 09:00",
    "due_date": "2025-03-15 10:00"
  }'
```

## Automatic Status Updates
The system automatically updates task statuses every 60 seconds (configurable via `CRON_SCHEDULE` environment variable):
- Updates **pending → in-progress** when `start_date` passes
- Updates **in-progress → overdue** when `due_date` passes

## Security
- Password hashing using bcrypt
- Redis session storage
- HTTP-only cookies
- Input validation for all endpoints

## Error Handling
Standard error format:
```json
{
  "status": 400,
  "error": "Validation failed: missing title field"
}
```

## Contact
For questions or support, contact [email](rashisky007@gmail.com).
