# Auth Service

A production-style authentication service built in **Go** using **Gin**, **PostgreSQL**, **Redis**, and **JWT**.  
This project demonstrates real-world backend patterns such as secure password handling, token-based authentication, middleware-driven authorization, and clean service structure.

---

## Features

* User signup and login
* Password hashing with bcrypt
* JWT access tokens
* Refresh tokens stored and managed in Redis
* Token refresh and logout
* Protected routes using middleware
* Rate Limiting
* Reset Password Emails
* Environment-based configuration

---

## Tech Stack

* **Go**
* **Gin**
* **PostgreSQL**
* **Redis**
* **JWT**
* **bcrypt**

---

## Project Structure

```text
cmd/
└── server/
    └── main.go         # Application entrypoint
internal/
├── auth/               # Auth domain (handlers, repository, models, JWT)
├── config/             # Application configuration
├── middleware/         # Auth, rate-limit, and other middleware
├── redis/              # Redis client and helpers
└── router/             # HTTP route definitions
└── email/              # Email Package (Client, Sender, Templates)
```



**Design approach:**  
Each domain owns its logic (handlers + repository + models), keeping the codebase modular, testable, and easy to extend.

---

## Prerequisites

- Go ≥ 1.20
- PostgreSQL
- Redis

---

## Setup & Run

### 1. Clone the repository

```bash
git clone https://github.com/faizan1191/auth-service.git
cd auth-service
```

## 2. Set environment variables

```bash
# .env
PORT=8080
JWT_SECRET=your_secret_key
BREVO_API_KEY=your_brevo_api_key
POSTGRESQL_DATABASE_URL=postgres://user@host:5432/database?sslmode=disable
REDIS_ADDR=localhost:6379
SENDER_EMAIL=brevo_sender_email
MAILTRAP_HOST=smtp.mailtrap.io
MAILTRAP_PORT=587
MAILTRAP_USER=user
MAILTRAP_PASS=password
MAILTRAP_FROM="no-reply@auth-service.com"
```

## 3. Create database and table

```sql
CREATE DATABASE auth_service;

CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## 4. Start Redis

```bash
redis-server
```

## 5. Run the application

```bash
go run cmd/server/main.go
```

## 📡 API Endpoints Overview

This service exposes authentication and user-related APIs.  
Endpoints are grouped into **public** and **protected** categories for clarity.

---

## 🔓 Public Endpoints

These endpoints **do not require authentication**.

| Method | Endpoint | Description |
|------|---------|-------------|
| `POST` | `/auth/signup` | Register a new user |
| `POST` | `/auth/login` | Authenticate user and issue access & refresh tokens |
| `POST` | `/auth/refresh` | Generate a new access token using a refresh token |
| `POST` | `/auth/logout` | Invalidate refresh token |
| `POST` | `/auth/forgot-password` | Send password reset link to user email |
| `POST` | `/auth/reset-password` | Reset password using reset token |
| `GET` | `/health` | Health check endpoint |
---

## 🔐 Protected Endpoints

These endpoints **require a valid access token**.

### Authorization Header

```http
Authorization: Bearer <access_token>
```

| Method | Endpoint | Description |
|------|---------|-------------|
| `GET` | `/api/me` | Get details of the currently authenticated user |
---
