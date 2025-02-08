# gotgbot-template

This project provides a template for building Telegram bots, leveraging the gotgbot library, dependency injection with fx, and other helpful tools.

## Getting Started

Clone the repository or use as template on github

```bash
git clone https://github.com/fullpipe/gotgbot-template.git
```

Add your env vars

```bash
cp .env.dist .env
```

Make sure you have [task](https://taskfile.dev/installation/) and [gowatch](https://github.com/silenceper/gowatch) installed

Start bot with

```bash
task dev
```

## Key Features and Concepts

### Gotgbot

This project uses the gotgbot library for interacting with the Telegram Bot API. Refer to the gotgbot documentation for details on using the library: [gotgbot](https://github.com/PaulSonOfLars/gotgbot).

### Dependency Injection (DI) with fx

The project utilizes the `fx` framework for dependency injection, promoting modularity and testability.

New repositories can be registered in `di/repository.go`

```go
fx.Provide(repository.NewUserRepo)
```

Controllers are registered in `cmd/bot/bot.go` with `di.AsController` wrapper

```go
di.AsController(controller.NewMessageController),
```

### task

Simplifies common development tasks like running the application, generating mocks, and generating GraphQL code. See the [Taskfile.yml](Taskfile.yml) for available commands.

### mockery

Mocks for testing are generated using [mockery](https://vektra.github.io/mockery/latest/).

```bash
task mockery
```

### GraphQL

For MiniApps api we use [gqlgen](https://github.com/99designs/gqlgen).

Update schema in (api)[api], then run `task gen` to update/generate handlers.

To run GraphQL server

```sh
go run . graph
```

### Metrics

Prometheus metrics are exposed at `:9090/metrics`, allowing monitoring of the application's performance.
