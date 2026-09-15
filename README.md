# hike-log

## Development Set-up

To run the application locally, use docker compose: `docker compose up`. This runs the frontend, backend, database (PostgreSQL), and infra (AWS via Localstack) in containers. The frontend is hosted on [http://localhost:5173](http://localhost:5173), while the backend runs on [http://localhost:8080](http://localhost:8080).

See the [frontend README](frontend/README.md) and [backend README](backend/README.md) for more specific instructions for those environments.

## Localstack Set-up

You must obtain an [auth token](https://docs.localstack.cloud/aws/getting-started/auth-token/) in order for Localstack to work. Add it to a `.env` file in the root of the project:

```
LOCALSTACK_AUTH_TOKEN=<your_token_here>

```
