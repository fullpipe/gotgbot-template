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

New repositories can be registered in [di/repository.go](di/repository.go)

```go
fx.Provide(repository.NewUserRepo)
```

#### Controllers

As you can see in [controller/start.go](controller/start.go). Controller is a struct which implements `di.Controller` interface.

Controllers are registered in [cmd/bot/bot.go](cmd/bot/bot.go) with `di.AsController` wrapper

```go
di.AsController(controller.NewMessageController),
```

### task

Simplifies common development tasks like running the application, generating mocks, and generating GraphQL code. See the [Taskfile.yaml](Taskfile.yaml) for available commands.

### mockery

Mocks for testing are generated using [mockery](https://vektra.github.io/mockery/latest/).

```bash
task mockery
```

### GraphQL

For Telegram MiniApps api we use [gqlgen](https://github.com/99designs/gqlgen).

Update schema in (api)[api], then run `task gen` to update/generate handlers.

To run GraphQL server

```sh
go run . graph
```

#### Auth

On your telegram miniapp client you should add `Authorization` header to all requests.

```ts
headers: {
    Authorization: `tma ${window.Telegram.WebApp.initData}`,
},
```

<details>
  <summary>Angular + apollo/client example</summary>

```ts
import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import { provideHttpClient } from '@angular/common/http';
import { provideApollo } from 'apollo-angular';
import { ApolloLink, InMemoryCache } from '@apollo/client/core';
import { environment } from '../environments/environment';
import { setContext } from '@apollo/client/link/context';
import { ErrorResponse, onError } from '@apollo/client/link/error';
import { createUploadLink } from 'apollo-upload-client';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';

export const appConfig: ApplicationConfig = {
  providers: [
    provideAnimationsAsync(),
    provideRouter(routes),
    provideHttpClient(),
    provideApollo(() => {
      const error = onError((e: ErrorResponse) => {
        console.log(e);
      });

      const auth = setContext((operation, context) => {
        if (!window.Telegram.WebApp.initData) {
          return {};
        } else {
          return {
            headers: {
              Authorization: `tma ${window.Telegram.WebApp.initData}`,
            },
          };
        }
      });

      return {
        link: ApolloLink.from([error, auth, createUploadLink({ uri: environment.graphUrl })]),
        cache: new InMemoryCache(),
        defaultOptions: {
          watchQuery: {
            fetchPolicy: 'no-cache',
            errorPolicy: 'ignore',
          },
          query: {
            fetchPolicy: 'no-cache',
            errorPolicy: 'all',
          },
        },
      };
    }),
  ],
};
```

</details>

### Metrics

Prometheus metrics are exposed at `:9090/metrics`, allowing monitoring of the application's performance.
