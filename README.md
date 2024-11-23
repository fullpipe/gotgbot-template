# gotgbot-template

## start

```sh
docker-compose up -d
go run . migrate
go run . bot

# or
gowatch -args bot
```

## GraphQL

```sh
task gen
go run . graph
```

## Metrics

Prometheus metrics exposed on port `:9090/metrics`

## Dependency Injection

For DI we use [fx](https://github.com/uber-go/fx).

So to register new repository add `fx.Provide(repository.NewUserRepo)`
and to register new controller

```go
		fx.Provide(di.AsController(controller.NewStartController)),
```

where `NewStartController` function to initialize controller

