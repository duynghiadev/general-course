# Software Architecture and Development Learning Path

## Table of Contents

- [Load Balancing in Go](#load-balancing-in-go)
- [Clean Architecture](#clean-architecture)
  - [Go Implementation](#go-implementation)
  - [Express.js Implementation](#expressjs-implementation)
- [Architecture Patterns](#architecture-patterns)
- [Dependency Injection](#dependency-injection)
- [Design Patterns in Go](#design-patterns-in-go)

## Load Balancing in Go

### Key Concepts

- Round-robin load balancing
- Weighted round-robin
- Least connections
- IP hash-based distribution

### Implementation Examples

- HTTP load balancer
- TCP/UDP load balancer
- Service discovery integration
- Health checks and circuit breakers

## Clean Architecture

### Core Principles

- Independence of frameworks
- Testability
- Independence of UI
- Independence of Database
- Independence of external agencies

### Go Implementation

- Domain layer (Entities)
- Use cases (Application layer)
- Interface adapters
- Infrastructure layer
- Repository pattern implementation
- Dependency rule

### Express.js Implementation

- Routes layer
- Controllers
- Services
- Models
- Middleware
- Database adapters

## Architecture Patterns

### Common Patterns

- MVC (Model-View-Controller)
- MVVM (Model-View-ViewModel)
- Hexagonal Architecture
- Event-Driven Architecture
- Microservices Architecture
- Domain-Driven Design (DDD)

### Best Practices

- SOLID Principles
- Separation of Concerns
- Single Responsibility
- Interface Segregation
- Dependency Inversion

## Dependency Injection

### Go Implementation

- Constructor Injection
- Method Injection
- Interface-based Design
- Popular DI containers:
  - Wire
  - Dig
  - Container

### Benefits

- Testability
- Modularity
- Flexibility
- Maintainability

## Design Patterns in Go

### Creational Patterns

- Factory Method
- Abstract Factory
- Builder
- Singleton
- Prototype

### Structural Patterns

- Adapter
- Bridge
- Composite
- Decorator
- Facade
- Proxy

### Behavioral Patterns

- Observer
- Strategy
- Command
- State
- Template Method
- Chain of Responsibility

### Best Practices

- When to use each pattern
- Common implementation pitfalls
- Performance considerations
- Testing strategies

## Learning Resources

### Books

- Clean Architecture by Robert C. Martin
- Design Patterns by Gang of Four
- Domain-Driven Design by Eric Evans

### Online Resources

- Go documentation
- Express.js documentation
- Architecture blogs and articles
- Community forums and discussions

## Project Examples

- Sample implementations
- Code repositories
- Best practices demonstrations
- Testing strategies

## Contributing

Feel free to contribute to this learning path by:

- Adding new resources
- Improving existing content
- Sharing your experiences
- Suggesting better practices

## 📁 Folder Structure Suggestions for Projects

```bash
my-app/
├── cmd/               # Entry point (main.go)
├── internal/
│   ├── domain/        # Entity definitions
│   ├── usecase/       # Business logic
│   ├── handler/       # HTTP handlers or controllers
│   ├── repository/    # DB interfaces
│   └── service/       # External services (email, cache, etc.)
├── pkg/               # Shared libraries/utilities
├── configs/           # Config files (YAML, ENV)
├── scripts/           # Dev scripts
└── docs/              # Documentation
```

## 📖 Resources

> This learning path is compiled from top blog posts, open-source projects, and documentation across the internet. Some helpful sources include:

- [Uber’s Go Style Guide](https://github.com/uber-go/guide)
- [Go Clean Architecture Blog by Uncle Bob adaptation]()
- [Mario Carrion&#39;s Go Microservice Blog &amp; Videos](https://mariocarrion.com/)
- [Golang Load Balancer Tutorial]()
- [Go Patterns](https://github.com/tmrts/go-patterns)
