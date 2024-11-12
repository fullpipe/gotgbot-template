# gotgbot-template

## start

```sh
docker-compose up -d
go run . migrate
go run . bot

# or
gowatch -args bot
```

## Dependency Injection

For DI we use [fx](https://github.com/uber-go/fx).

So to register new repository add `fx.Provide(repository.NewUserRepo)`
and to register new controller

```go
		fx.Provide(di.AsController(controller.NewStartController)),
```

where `NewStartController` function to initialize controller

