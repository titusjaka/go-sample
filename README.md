# go-sample

[![GitHub release](https://img.shields.io/github/v/release/titusjaka/go-sample)](https://github.com/titusjaka/go-sample/releases/latest)
[![CI](https://github.com/titusjaka/go-sample/workflows/CI/badge.svg)](https://github.com/titusjaka/go-sample/actions)
[![codecov](https://codecov.io/gh/titusjaka/go-sample/branch/main/graph/badge.svg?token=UNJY7V5SZL)](https://codecov.io/gh/titusjaka/go-sample)
[![Go Report Card](https://goreportcard.com/badge/github.com/titusjaka/go-sample)](https://goreportcard.com/report/github.com/titusjaka/go-sample)
[![GitHub license](https://img.shields.io/github/license/titusjaka/go-sample)](https://github.com/titusjaka/go-sample/blob/main/LICENSE)

Go Backend Sample. It’s suitable as a starting point for a REST-API Go application.

This example uses:
  - [chi](https://github.com/go-chi/chi) for HTTP router;
  - [cli](https://github.com/urfave/cli) for building neat commands;
  - [embedded-postgres](https://github.com/fergusstrange/embedded-postgres) for database testing;
  - [go-kit](https://github.com/go-kit/kit) for transport layer;
  - [go.rice](https://github.com/GeertJohan/go.rice) for file embedding;
  - [PostgreSQL](https://www.postgresql.org/) as a database and [pq](https://github.com/lib/pq) as a driver;
  - [testify](https://github.com/stretchr/testify) and [mock](https://github.com/golang/mock) for tests;
  - [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) for request validation;
  - [zap](https://go.uber.org/zap) as a logger;
  - [the Clean Architecture](http://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) as the primary approach.

Features:

  - Modular Project Structure.
  - Built-in migration handling.
  - Ready to go example with tests and mocks.
  - Easy-to-go Docker deployment.

## Structure
```text
go-sample
├── cmd
│   └── cli                 // Entry point for the application.
├── commands                // Sub-commands for CLI (stands for Command Line Interface).
│   ├── migrate.go          // Migration sub-command. It creates, applies, and rollbacks DB-migrations.
│   └── server.go           // Server sub-command. It runs the HTTP server.
├── internal                // Internal packages for the application according to Go convention. 
│   ├── business            // Business logic of the application.
│   │   └── snippets        // A specimen business-logic package “snippets” with REST-API for snippets creating, listing, and deleting.
│   └── infrastructure      // Infrastructure code of the application.
│       ├── api             // API-related utilities: middlewares, authentication, error handling for the transport layer.
│       ├── log             // Convenient logger to decouple from the standard log package.
│       ├── migrate         // Reusable migration utilities.
│       └── service         // Service-related reusable code: error handling for the service layer, etc.
└── migrations              // This folder contains *.sql migrations.
```

## Installation

```shell
git clone https://github.com/titusjaka/go-sample
```

## Usage

```shell
docker-compose up --build
```


## Future improvements:
- [ ] Return verbose API errors with exact fields in it:
    ```json
    {
      "errors": {
        "title": "Title must not be empty",
        "expires_at": "Expires at must be within range 1-365 days"
      }
    }
    ```
- [ ] Add user authentication + session storage.
- [ ] Add `/status` handler with service health.
- [ ] Reuse DB connection in tests.
