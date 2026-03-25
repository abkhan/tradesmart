# Gemini Context: mongotest

This file provides foundational mandates and context for Gemini CLI in this workspace.

## Project Overview
`mongotest` is a Go-based project designed to demonstrate and test connections to MongoDB using the official Go driver.

## Tech Stack
- **Language:** Go (Golang)
- **Database:** MongoDB (using `go.mongodb.org/mongo-driver`)

## Foundational Mandates
- **Security:** Never commit or log the MongoDB connection string, specifically the password. Use environment variables for sensitive credentials in production-ready code.
- **Error Handling:** Follow idiomatic Go error handling. Avoid using `panic()` in library code; return errors to the caller instead.
- **Context Management:** Always use `context.Context` for MongoDB operations to ensure proper timeout and cancellation support.

## Development Standards
- **Naming:** Use camelCase for internal variables and PascalCase for exported symbols.
- **Formatting:** Run `go fmt ./...` before finalizing changes.
- **Testing:** Add unit tests for new logic using the standard `testing` package.

## Workspace Conventions
- Entry points are located in `cmd/`.
- Keep `main.go` focused on orchestration; move business logic or database abstractions to specialized packages as the project grows.

## Other information for coding
- DB to be used:   tradesmart
- for sales data, the collection is: sales