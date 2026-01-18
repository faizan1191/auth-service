# Auth Service

A production-style authentication service built in **Go** using **Gin**, **PostgreSQL**, **Redis**, and **JWT**.  
This project demonstrates real-world backend patterns such as secure password handling, token-based authentication, middleware-driven authorization, and clean service structure.

---

## Features

- User signup and login
- Password hashing with bcrypt
- JWT access tokens
- Refresh tokens stored and managed in Redis
- Token refresh and logout
- Protected routes using middleware
- Rate Limiting
- Environment-based configuration

---

## Tech Stack

- Go
- Gin
- PostgreSQL
- Redis
- JWT
- bcrypt

---

## Project Structure

cmd/server/
└── main.go # Application entrypoint

internal/
├── auth/ # Auth domain (handlers, repository, models, JWT)
├── config/ # Application configuration
├── middleware/ # Auth, rate-limit, and other middleware
├── redis/ # Redis client and helpers
└── router/ # HTTP route definitions



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

## 2. Set environment variables

```bash
export PORT=8080
export JWT_SECRET=your_secret_key
export POSTGRESQL_DATABASE_URL="postgres://<user>@localhost:5432/auth_service?sslmode=disable"

## 3. Create database and table

```sql
CREATE DATABASE auth_service;

CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


## 4. Start Redis

```bash
redis-server

## 5. Run the application

```bash
go run cmd/server/main.go

## API Overview
Public Endpoints

POST /auth/signup

POST /auth/login

POST /auth/refresh

POST /auth/logout

GET /health

## Protected Endpoints

GET /api/me
Requires header: Authorization: Bearer <access_token>