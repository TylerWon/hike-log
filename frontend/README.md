# Frontend

## Installing/Uninstalling Packages

Since the application runs in a container, it is easiest to install/uninstall the package from in the container:

1. Start the services: `docker compose up`
2. Open a terminal in the frontend container: `docker exec -it hike-log-frontend-1 sh`
3. Install or uninstall your package
4. Restart all the services: `docker compose stop && docker compose up`

## Linting and Formating

To lint all code: `npm run lint`
To format all code: `npm run format`

It is recommended to install the [ESLint](https://marketplace.cursorapi.com/items/?itemName=dbaeumer.vscode-eslint) and
[Prettier](https://marketplace.cursorapi.com/items/?itemName=esbenp.prettier-vscode) extensions for your code editor to
automate these actions.
