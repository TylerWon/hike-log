# Backend

## Installing/Uninstalling Packages

1. Start the services: `docker compose up`
2. Open a terminal in the backend container: `docker exec -it hike-log-backend-1 sh`
3. Add the dependency: `go get <module>`
4. Clean-up go.mod: `go mod tidy`
5. Restart services if hot reload doesn't pick it up: `docker compose stop && docker compose up`

## Linting and Formating

To lint all code: `go vet ./...`
To format all code: `go fmt ./...`

Install the [Go](https://marketplace.cursorapi.com/items/?itemName=golang.Go) extension to lint and format automatically when a file is saved. These settings can be adjusted in `.vscode/settings.json`.

## Testing

Tests should be run in the development container: `docker compose up` then `docker exec -it hike-log-backend-1 sh`.

To run all tests: `go test ./...`.

To run all tests in a package: `go test <path to package>`

- Example: `go test ./handler` runs all tests in the `handler` package.

To run specific tests in a package: `go test <path to package> -run '<test suite name>/<test name>'

- Example: `go test ./handler -run 'TestHandlerTestSuite'` runs all tests in the `handler` package that are part of the `TestHandlerTestSuite` suite.
- Example: `go test ./handler -run 'TestHandlerTestSuite/TestCheckAuth'` runs all tests in the `handler` package that are part of the `TestHandlerTestSuite` suite and start with the prefix `TestCheckAuth`.
