## Initialize the project:
```
make install-tools
make generate
make tidy
```
## Set up development environment:

```
chmod +x scripts/dev.sh
./scripts/dev.sh
```

## Run the application:
```
make run
```

# Common development workflow:

## Create a new migration
```
make migrate-create
```

## Run migrations
```
make migrate-up
```

## Generate code after changes
```
make generate
```

## Run tests and linting
```
make test
make lint
```

## Build for production
```
make build
```

This setup provides a complete development environment with all the tooling you need for a professional Go project with gRPC, SQLC, and clean architecture.