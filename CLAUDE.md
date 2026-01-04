Read
- [README.md](mdc:README.md) for general information.
- [WHY.md](mdc:WHY.md) for tech stack, library choice decisions.

# General
- Always explain your reason of your decision, please.
- When I tell you about problems, issues, please analyze them step by step and propose me a simple and sophisticated generic solution. Your solution should be generic and MUST NOT be case-specific.

# Coding guidance
Keep things minumum for maintainability and simplicity.
Follow the programming principles such as

- SOLID
  - S - Single-responsibility Principle
  - O - Open-closed Principle
  - L - Liskov Substitution Principle
  - I - Interface Segregation Principle
  - D - Dependency Inversion Principle
- KISS (Keep It Simple, Stupid)
- DRY (Do not Repeat Yourself)
- Composition Over Inheritance
- Separation of Concerns

# Test
Run `make test` to test. This repo uses [Dockerfile](mdc:../../Dockerfile) and [docker-compose.test.yaml](mdc:docker-compose.test.yaml) for testing so normal `go run test` wouldn't succeed as it won't connect to database.

# Documentation
1. Update `README.md`, `README_JA.md`, `WHY.md`, `WHY_JA.md` when you update code.
2. When you update [README.md](mdc:README.md) or [WHY.md](mdc:WHY.md), always update [README_JA.md](mdc:README_JA.md) or [WHY_JA.md](mdc:WHY_JA.md) accordingly with natural translation, not like machine translation.
3. Don't use `we` in the documents as this is a template that's considered to be forked or derived. I'd like you to use `it` to refer to the project rather than persons.
