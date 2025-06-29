[🇯🇵 **日本語**](./README_JA.md)

[**Wire**](https://github.com/google/wire) is a compile-time dependency injection (DI) tool for Go. It analyzes your constructors (functions that return initialized types) and automatically generates the "glue" code that wires them together. In other words, instead of manually calling each `NewFoo(...)` in a builder function, you declare all constructors in a `wire.Build(...)` statement, and Wire writes the exact sequence of calls for you.

---

## What Wire Is For

1. **Centralizing and validating wiring logic**

   * You list every constructor only once. Wire inspects their signatures and builds a dependency graph at compile time.
   * If any constructor is missing a required input—or two constructors become incompatible—Wire fails immediately with a clear error, pointing out exactly which provider is missing or mismatched.

2. **Generating boilerplate code**

   * Rather than hand-writing a long "BuildServer" function that calls `LoadConfig()`, `NewLogger()`, `NewDBClient()`, etc., Wire writes that code for you in a generated file (`wire_gen.go`).
   * You never edit or maintain that generated file; you only maintain your constructors and the `wire.go` declarations.

3. **Keeping DI explicit but automated**

   * All your constructors remain ordinary Go functions. Wire does not add runtime reflection or magic. At build time, you run `wire`, it produces plain Go code, and your application compiles exactly as if you had written that code yourself.

---

## Pros of Using Wire

1. **Compile-time safety**

   * If you change a constructor's signature (e.g. add a new parameter), Wire immediately reports "unbound provider" or "cannot assign type X to Y," so you fix the wiring before running the program.
   * You catch missing or incompatible dependencies as part of the build process, not at runtime.

2. **Scalability for large dependency graphs**

   * Once you have 10+ constructors, manual wiring becomes verbose and error-prone. Wire keeps you from forgetting an argument or misordering calls: it builds the entire graph by inspecting each constructor's inputs and outputs.
   * Adding or removing a service means updating one list in `wire.go`; Wire updates the generated code everywhere it's used.

3. **Easier to swap implementations (for testing or alternative modes)**

   * You can use `wire.Bind` or `wire.Value` to tell Wire how to satisfy an interface or to inject a stub/mock.
   * Test injectors (e.g. `InitializeTestServer`) can be defined separately—Wire will generate boilerplate that passes fakes into the same constructors, automatically ensuring your test setup stays in sync with production constructors.

4. **Supports multiple entry points with shared core dependencies**

   * If you have an HTTP server, a CLI tool, a background worker, or scheduled jobs, each entry point often shares most of the same components (Config → Logger → DB → Repository → ServiceA, etc.). With Wire, you simply declare three different injectors (`InitializeServer`, `InitializeWorker`, `InitializeCLI`) that all refer to the same constructors.
   * You avoid copy-pasting 80% of your wiring logic three times—Wire generates each injector's code based on a shared list of providers.

5. **Clean, readable code separation**

   * Your "business logic" constructors stay in their own files (e.g. `db.go`, `repo.go`, `service.go`). The only place you see a big list of constructors is in `wire.go`.
   * `main.go` remains minimal: it just calls `InitializeServer()` (or `InitializeWorker()`), checks for an error, and starts the process. All the "how to wire things" is encapsulated in Wire's generated code.

6. **CI/CD safety net**

   * You can add `wire` to your CI pipeline. If someone forgets to update a `wire.Build(...)` after changing a constructor, `wire` fails and your build breaks. You never accidentally ship broken wiring to production.
   * Because Wire's errors are explicit, you immediately know which constructor or provider is missing.

7. **No runtime reflection or performance overhead**

   * Generated code is plain Go: there's no reflection, no hidden container at runtime, no lookup tables—everything is resolved at compile time. Startup performance is identical to manually written wiring.

---

See [Wire's Github](https://github.com/google/wire) for further information.
