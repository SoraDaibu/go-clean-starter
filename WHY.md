# Technology Decisions for `go-clean-starter`

[🇯🇵 **日本語**](./WHY_JA.md)

This document explains the key technology decisions behind `go-clean-starter`. It exists to help developers quickly understand the purpose of each major tool or library in the stack and decision history of them.

## 🧱 Context

`go-clean-starter` is a backend template built with **Go** and designed around **Clean Architecture** principles. It aims to provide a solid foundation for scalable, testable, and maintainable applications—suitable for small to medium teams, indie products, or internal tools.

## 📚 Key Libraries

- **[Echo](https://github.com/labstack/echo)** - Fast and lightweight HTTP web framework for Go
- **[sqlc](https://github.com/kyleconroy/sqlc)** - Generate type-safe Go code from SQL queries
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migration tool with version control
- **[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)** - Generate Go types from OpenAPI 3.x specifications
- **[air](https://github.com/air-verse/air)** - Hot reload tool for Go development
- **[pgx](https://github.com/jackc/pgx)** - High-performance PostgreSQL driver with rich PosgreSQL features and toolkit
- **[Zerolog](https://github.com/rs/zerolog)** - Fast and simple JSON logger
- **[CLI v3](https://github.com/urfave/cli)** - Command line interface framework
- **[Testify](https://github.com/stretchr/testify)** - Testing toolkit with assertions and mocks


## 🧠 Core Decisions

### 🔹 Language: **Go**

Go is used for its:
- Simplicity and clarity
- High performance and low resource usage
- First-class concurrency support
- Strong ecosystem for backend and DevOps

---

### 🔹 Web Framework: **Echo**

[Echo](https://github.com/labstack/echo) is used as the HTTP framework because:
- It's fast and lightweight
- It offers a clean, intuitive router
- Middleware support is rich and extensible
- It integrates easily with `net/http` and works well in REST API scenarios

---

### 🔹 Database: **PostgreSQL**

PostgreSQL was chosen as the primary database because:
- It's reliable, open-source, and widely adopted
- It supports rich features like JSONB, full-text search, and strong consistency
- It works well with both ORMs and raw SQL workflows

---

- [`golang-migrate`](https://github.com/golang-migrate/migrate) is used for schema versioning, which is choosen because `golang-mirate` supports variety of databases and its simplicity.

---

### 🔹 SQL Tooling: **sqlc + golang-migrate**

* **[`sqlc`](https://github.com/kyleconroy/sqlc)** generates **type-safe Go code** from raw SQL queries.

  * Keeps business queries close to the database for better transparency and debugging.
  * No runtime reflection — efficient and safe at compile time.
  * Makes it easy for developers familiar with SQL (but not Go) to contribute effectively.
  * Also improves debuggability by allowing raw SQL inspection.

* **[`golang-migrate`](https://github.com/golang-migrate/migrate)** handles schema migrations with raw SQL using `.up.sql` and `.down.sql` files.
  * Supports version control, locking, and rollback.
  * Has built-in mechanisms to prevent conflicts when multiple migration processes run concurrently, which works well in CI/CD and multi-environment deployments.

---

### 🔹 OpenAPI Code Generation: **[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)**

**oapi-codegen** is used to generate Go structs and types from OpenAPI 3.x specifications because:

* **Type safety between API spec and implementation**
  * Automatically generates Go types that match the OpenAPI schema definitions, ensuring consistency between documentation and code.
  * Prevents runtime errors caused by mismatched request/response structures.

* **Single source of truth**
  * The OpenAPI specification (`doc/api.yaml`) serves as the authoritative definition of the API contract.
  * Changes to the API schema automatically propagate to the Go code through regeneration.

* **Reduced boilerplate and maintenance**
  * Eliminates the need to manually write and maintain request/response structs.
  * Automatically handles complex nested structures, validation tags, and type conversions.

oapi-codegen ensures that the API implementation stays synchronized with its specification, reducing bugs and improving developer productivity when working with REST APIs.

---

### 🔹 Dependency Injection: **Manual**

**Manual dependency injection** is used to keep object initialization clean, explicit, and maintainable because:
* Simple and idiomatic—follows standard Go practices without additional tooling
* Explicit dependency graphs make code easy to read and debug
* No build-time code generation or extra dependencies required
* Constructors are centralized in the `builder` package
* Encourages a modular, testable architecture (especially useful in Clean Architecture setups)

Manual DI keeps the project simple and maintainable without sacrificing clarity as services, repositories, and external clients increase in number.

### 🔹 Hot reload: **[air](https://github.com/air-verse/air)**

**Air** is used for hot reload during development because:
* Provides fast, automatic rebuilds when Go source files change, significantly improving developer productivity.
* Lightweight and focused—designed specifically for Go applications without unnecessary complexity.
* Configurable via `.air.toml` files, allowing customization of watched directories, build commands, and exclusion patterns.
* Works seamlessly with Docker and containerized development environments.

Air eliminates the tedious cycle of manually stopping, rebuilding, and restarting the application during development, making the feedback loop much faster for iterative coding.

---

### 🔹 PostgreSQL Driver: **[pgx](https://github.com/jackc/pgx)**

**pgx** is used as the PostgreSQL driver (native API, not stdlib mode) because:

* **High performance and efficiency**
  * Supports both text and binary wire protocols for optimal performance
  * 2-3x faster than traditional database/sql drivers for many operations
  * More efficient memory usage with fewer allocations
  * Native connection pooling (`pgxpool`) with prepared statement caching

* **Rich PostgreSQL feature support**
  * Native support for PostgreSQL-specific types (arrays, JSONB, hstore, geometric types, etc.)
  * Built-in COPY protocol support for high-speed bulk operations
  * LISTEN/NOTIFY support for real-time event notifications
  * Better handling of custom and composite types

* **Active development and modern design**
  * Actively maintained with regular updates and performance improvements
  * Built for modern Go patterns and PostgreSQL versions
  * Better observability with connection pool metrics and health checks
  * Strong type safety when used with sqlc's `sql_package: "pgx/v5"` option

* **Production-ready connection management**
  * Sophisticated connection pooling with automatic health checks
  * Better retry logic and error handling
  * Support for connection lifecycle hooks and query logging

While pgx can be used as a drop-in `database/sql` driver (stdlib mode), this project uses the **native pgx API** to take full advantage of its performance and PostgreSQL-specific features. This decision aligns with the goal of providing a high-performance, feature-rich template for production use.

---
## ✏️ Alternatives Considered

### **Database: MySQL**

`MySQL` was skipped in favor of `PostgreSQL` for a few reasons:

1. **Advanced feature set**

   * PostgreSQL supports row-level security (RLS), JSONB with native indexing, richer indexing options (GIN, GiST, partial indexes), and built-in Full Text Search.
   * MySQL lacks JSONB and these advanced indexing capabilities, making it harder to handle semi-structured data or complex text queries without workarounds.

2. **Extension ecosystem**

   * PostgreSQL has mature extensions like PostGIS (for geospatial queries) and pg\_trgm (for similarity searches).
   * MySQL's ecosystem does not provide equivalent built-in geospatial or high-performance text search support.

3. **Flexibility for evolving requirements**

   * JSONB in PostgreSQL allows storing and querying dynamic structures without needing a separate NoSQL store.
   * MySQL's more limited feature set can become a bottleneck when future use cases demand richer data types or indexing.

---

### **Migration tool: [Goose](https://github.com/pressly/goose)**

`Goose` was considered, but `golang-migrate` was chosen instead:

1. **Multi-database support**

   * `golang-migrate` natively handles PostgreSQL, MySQL, SQLite, MongoDB, Cassandra, etc...
   * Goose supports fewer dialects, which could limit flexibility. I wanted to have a starter template that serves for variety of tech stacks.

2. **Locking and concurrency**

   * `golang-migrate` uses advisory locks (e.g., PostgreSQL advisory locks) to ensure only one migration runs at a time in CI/CD or multi-developer environments.
   * Goose lacks built-in advisory locking, requiring external coordination to avoid race conditions.

3. **Extensible CLI**

   * `golang-migrate` includes commands like `version`, `status`, and `force`, and integrates smoothly with Docker-based pipelines for automated workflows.
   * Goose's CLI (`up`, `down`, `fix`, `status`, `version`) is simpler but less flexible for complex CI/CD scenarios.

---

### **ORM: [GORM](https://gorm.io/)**

`GORM` was skipped in favor of `sqlc` for several reasons:

1. **Compile-time safety**
   * `sqlc` generates Go code from SQL queries, ensuring any schema-query mismatch is caught at compile time.
   * GORM relies on runtime reflection, which can hide errors until execution.

2. **Explicit SQL and debuggability**
   * `sqlc` produces idiomatic, explicit SQL calls, making it easier to identify and optimize performance bottlenecks without magic joins or unintended eager loads.
   * GORM's abstractions can lead to opaque behavior that's harder to trace when debugging complex queries.

3. **Performance**
   * `sqlc` generates static code with no reflection overhead, resulting in faster query execution under load.
   * GORM's reflection-based query building incurs additional runtime cost.

---

### **OpenAPI Code Generation: [OpenAPI Generator](https://github.com/OpenAPITools/openapi-generator)**

`OpenAPI Generator` was considered, but `oapi-codegen` was chosen instead:

1. **Go-idiomatic code generation**

   * `oapi-codegen` is designed specifically for Go and generates clean, idiomatic Go code that follows Go conventions.
   * OpenAPI Generator supports many languages but often produces verbose, less Go-idiomatic code that requires additional customization.

2. **Focused scope and simplicity**

   * `oapi-codegen` focuses on generating models and client/server interfaces, which aligns perfectly with our "models only" use case. While it currently only generates models, `oapi-codegen` was chosen for its extensibility and customizability.
   * OpenAPI Generator is a comprehensive tool designed for full code generation across multiple languages, making it overkill for simple model generation.

3. **Integration with Go ecosystem**

   * `oapi-codegen` integrates seamlessly with popular Go frameworks like Echo, Chi, and net/http without additional configuration.
   * OpenAPI Generator requires template customization and additional adaptation work to integrate cleanly with Go web frameworks.

4. **Maintenance and configuration overhead**

   * `oapi-codegen` requires minimal configuration - a simple command-line flag specifies what to generate.
   * OpenAPI Generator requires complex template management and configuration files, adding unnecessary complexity for straightforward model generation.

---

### **PostgreSQL Driver: [lib/pq](https://github.com/lib/pq)**

`lib/pq` was considered as it was the de facto standard PostgreSQL driver, but `pgx` was chosen instead:

1. **Active maintenance and future-proofing**
   * `lib/pq` is in maintenance mode since 2018 with no new feature development
   * `pgx` is actively developed with regular performance improvements and bug fixes
   * `pgx` supports the latest PostgreSQL features and Go versions

2. **Performance advantages**
   * `pgx` is 2-3x faster than `lib/pq` for most operations
   * Binary protocol support in `pgx` reduces serialization overhead
   * More efficient memory allocation and connection pooling in `pgx`

3. **PostgreSQL feature support**
   * `pgx` provides native support for PostgreSQL arrays, JSONB, and custom types
   * `lib/pq` requires manual marshaling for many PostgreSQL-specific features
   * `pgx` supports advanced features like COPY protocol and LISTEN/NOTIFY

4. **Better tooling integration**
   * `sqlc` has first-class support for pgx with `sql_package: "pgx/v5"`
   * Type-safe code generation works seamlessly with pgx types
   * Better observability and debugging capabilities

While `lib/pq` remains stable and suitable for simple use cases, `pgx` offers significantly better performance, richer features, and active development, making it the better choice for a modern, production-ready template.

---

### **Dependency Injection Tools: [Wire](https://github.com/google/wire) / [Fx](https://github.com/uber-go/fx)**

Specialized DI frameworks were skipped in favor of manual dependency injection:

1. **Simplicity and clarity**
   * Manual DI uses plain Go code without code generation or framework abstractions
   * The dependency graph is explicit and easy to trace through the codebase
   * No build-time tooling or special commands required

2. **Standard Go idioms**
   * Most Go projects use manual DI for small-to-medium dependency graphs
   * Follows the principle of "clear is better than clever"
   * Reduces onboarding time for new developers familiar with Go

3. **No overhead or magic**
   * No code generation step (Wire) or runtime reflection (Fx)
   * Stack traces and IDE navigation work seamlessly
   * Easier debugging with explicit constructor calls

For starting a Go backend project, manual DI normally provides the best balance of simplicity, maintainability, and clarity. Also it's not much common to have a dependency injection tool for a project based in Go.

## ✅ Why This Stack?

- **Minimal, but production-friendly**
- **Highly performant**
- **Easy to understand and extend**
- Works well for **small to big teams**
- Encourages **clear separation of concerns** via Clean Architecture

## 🔄 Future Considerations

Depending on project needs, the following may be added:
- Authentication and session management (project-specific requirements)
- gRPC support for service-to-service communication
- Background job processing enhancements
- Monitoring/observability integrations (Prometheus, OpenTelemetry)
- Caching layer (Redis) for high-traffic scenarios and user session management

## 📝 Contributing

This project is open for contributions! If you suggest alternative libraries or patterns, please explain the tradeoffs in your PR or issue.
