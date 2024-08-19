# Go Backend API template

A clean architecture Go backend template built with Gin, PostgreSQL, GORM, and Uber FX. This project provides a robust starting point for building scalable and maintainable Go applications.

## Features

- Clean architecture with separation of concerns
- RESTful API using the Gin framework
- PostgreSQL database with GORM ORM
- Dependency injection with Uber FX
- Docker support for easy containerization
- AWS integration for S3,SES,etc...

## Getting Started

### Prerequisites

Before you begin, ensure you have met the following requirements:

- [Go](https://golang.org/dl/) version 1.23 or later
- [Docker](https://www.docker.com/get-started) and Docker Compose
- PostgreSQL (if running outside of Docker)

### Project structure

```shell
├── cmd
│    └── server
│         ├── main.go
│         └── server.go
├── internal
│   ├── aws_client
│   │    └── aws_client.go
│   ├── config
│   │    ├── config.go
│   │    └── default.go
│   ├── features
│   │   ├── auth
│   │   │   ├── auth_handler.go
│   │   │   ├── auth_service.go
│   │   │   ├── dto
│   │   │   │    ├── request.go
│   │   │   │    └── response.go
│   │   │   └── tokens
│   │   │       └── tokens.go
│   │   ├── email
│   │   │    ├── email_service.go
│   │   │    └── entities
│   │   │         └── email.go
│   │   └── user
│   │       ├── dto
│   │       │    ├── request.go
│   │       │    └── response.go
│   │       ├── encoder.go
│   │       ├── entity
│   │       │    └── user.go
│   │       ├── user_handler.go
│   │       ├── user_repository.go
│   │       └── user_service.go
│   └── postgres
│       ├── context.go
│       ├── errors.go
│       ├── helper.go
│       ├── logging.go
│       └── postgres.go
├── pkg
│   ├── errors
│   │    └── errors.go
│   ├── logging
│   │    └── logging.go
│   └── validator.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── README.md
```

### Installation
