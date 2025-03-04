# Project API Documentation

## Introduction

This project provides an API for user authentication and task management. It allows users to register, login, and manage tasks, including creating, retrieving, updating, and deleting tasks.

## Routes

### Public Routes

These routes do not require authentication.

#### **POST /login**

**Description:** Log in a user.

- **Request Body:**

  ```json
  {
      "email": "user@example.com",
      "password": "password123"
  }
  ```

- **Response:**
  - `200 OK`: User successfully logged in.
  - `401 Unauthorized`: Invalid email or password.

#### **POST /register**

**Description:** Register a new user.

- **Request Body:**

  ```json
  {
      "email": "user@example.com",
      "password": "password123",
      "fullname": "John Doe"
  }
  ```

- **Response:**
  - `201 Created`: User successfully registered.
  - `400 Bad Request`: Validation error.
  - `409 Conflict`: Email already in use.

---

### Protected Routes

These routes require authentication. Pass the token in the `Authorization` header:

```text
Authorization: Bearer <token>
```

#### **GET /task**

**Description:** Retrieve all tasks for the authenticated user.

- **Response:**

  ```json
  [
      {
          "taskID": "12345",
          "description": "Task description",
          "title": "Task title",
          "startDate": "2025-02-25 12:34",
          "dueDate": "2025-02-27 14:00",
          "status": "in-progress"
      }
  ]
  ```

#### **GET /task/:taskID**

**Description:** Retrieve a specific task by its ID.

- **Response:**

  ```json
  {
      "taskID": "12345",
      "description": "Task description",
      "title": "Task title",
      "startDate": "2025-02-25 12:34",
      "dueDate": "2025-02-27 14:00",
      "status": "in-progress"
  }
  ```

#### **POST /task**

**Description:** Create a new task.

- **Request Body:**

  ```json
  {
      "description": "Task description",
      "title": "Task title",
      "startDate": "2025-02-26 08:00",
      "dueDate": "2 days"
  }
  ```

  - `startDate`: Optional (defaults to current time if not provided).
  - `dueDate`: Required. Acceptable formats:
    - `YYYY-MM-DD HH:MM` (e.g., `2025-02-27 14:00`)
    - `x day(s)` (e.g., `2 days`)
    - `x hour(s)` (e.g., `5 hours`)
    - `x minute(s)` (e.g., `30 minutes`)

- **Response:**
  - `201 Created`: Task successfully created.
  - `400 Bad Request`: Validation error.

#### **PUT /task/:taskID**

**Description:** Update an existing task.

- **Request Body:**

  ```json
  {
      "description": "Updated task description",
      "title": "Updated task title",
      "startDate": "2025-03-01 09:00",
      "dueDate": "3 days"
  }
  ```

- **Response:**
  - `200 OK`: Task successfully updated.
  - `400 Bad Request`: Validation error.

#### **DELETE /task/:taskID**

**Description:** Delete a specific task by its ID.

- **Response:**
  - `200 OK`: Task successfully deleted.
  - `404 Not Found`: Task not found.

#### **DELETE /user**

**Description:** Delete the authenticated user.

- **Response:**
  - `200 OK`: User successfully deleted.
  - `401 Unauthorized`: Authentication required.

---

## Date Format and Validation

- **startDate:** Must be in the format `YYYY-MM-DD HH:MM`. Relative formats (e.g., `2 days`) are **not** allowed.
- **dueDate:** Can be in one of the following formats:
  - Absolute: `YYYY-MM-DD HH:MM`
  - Relative: `x day(s)`, `x hour(s)`, `x minute(s)`

---

## Error Handling

- **Validation Errors:**

  ```json
  {
      "errors": [
          "Field 'description' is required",
          "Field 'dueDate' must be after 'startDate'"
      ]
  }
  ```

- **Authentication Errors:**

  ```json
  {
      "error": "Unauthorized"
  }
  ```

---

## Technologies Used

- **Backend Framework:** Gin
- **Database:** GORM with a relational database (e.g., PostgreSQL, MySQL)
- **Authentication:** Gin sessions with Redis

---

## Setup and Installation

1. Clone the repository:

   ```bash
   git clone <repository_url>
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   go run main.go
   ```

---

## Future Enhancements

- Add pagination to the `GET /task` endpoint.
- Implement role-based access control for tasks.
- Add filtering and sorting options for tasks.

---

## Contact

For questions or support, contact [email](rashisky007@gmail.com).
