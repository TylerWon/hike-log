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

## Testing

There are three kinds of tests for the frontend: unit, component, and visual.
- Unit tests: Test the functionality of a block of code.
- Component tests: Test the functionality of a React component.
- Visual tests: Check for visual regressions in a React component by comparing a new screenshot of the component with an 
existing screenshot.

Tests can either be run once or in watch mode. Also, all test types can be run or just a specific type. To run tests,
use a variation of this command: `npm run test:<test type>:<execution mode>`. Here are some examples:
- Run all tests once: `npm run test:all:once`
- Run unit tests in watch mode: `npm run test:unit:watch`
- Run component tests once: `npm run test:component:once`
- Run visual tests in watch mode: `npm run test:visual:watch`

After updating components, visual tests may fail due to a change in the appearance of the components. If the changes are 
desired, run `npm run test:visual:update` to update the screenshots for the components so that the tests pass again.