---
name: go-dev
description: Golang development skill for building, testing, and deploying Go applications with Fiber framework and PostgreSQL database.
---

## Workflow
1. Analyze the requirements and design the architecture of the Go application from <user requirements>.
2. Implement the application using Go programming language, Fiber framework, and PostgreSQL database.
3. Write unit tests and integration tests to ensure the correctness of the application.
   3.1 If fail, fix the code and re-run the tests until they pass.
4. Create or update to file `MEMORY.md` to record the knowledge and lessons learned from the implementation.

## Technologies stack
* Go version 1.26.5
* Fiber framework version 3.5.0
* PostgreSQL version 16.0 with [pgx package](https://github.com/jackc/pgx) + connection pooling
* Unit Testing and Integration testing with [Testify](https://github.com/stretchr/testify) and [Testcontainers-go](https://golang.testcontainers.org/quickstart/) to spin up a PostgreSQL container for integration tests

## Project structure with Clean Architecture + feature-based approach 
* Flow
```
main.go -> config.go -> user_handler.go -> user_usecase.go -> user_repository.go -> user.go
```

```
api
├── cmd
│   └── main.go
├── config
│   └── config.go
├── internal
│   ├── app
│   │   ├── user
│   │   │   ├── delivery
│   │   │   │   └── http
│   │   │   │       └── user_handler.go
│   │   │   ├── repository
│   │   │   │   └── user_repository.go  # Use interface for repository implementation
│   │   │   └── usecase
│   │   │       └── user_usecase.go
│   │   │       └── user_usecase_test.go
│   │   └── ...
│   ├── domain
│   │   └── user.go
│   └── ...
├── pkg
│   └── ...
├── go.mod
└── go.sum
├── integration_tests
│   ├── user
│       └── user_success_integration_test.go
│       └── user_failure_integration_test.go
```

## Best practices
* Manage configuration using environment variables and a configuration package to load them.
* Use Clean Architecture principles to separate concerns and maintain a modular codebase.
* Follow a feature-based approach to organize code by domain features, making it easier to navigate and maintain.
* Use dependency injection to manage dependencies and improve testability.
* Error handling: Use proper error handling techniques, including custom error types and error wrapping with `%w`, to provide meaningful error messages and facilitate debugging. And check with `errors.Is` and `errors.As` to handle specific error types.