# Use Debian-based image for Playwright
FROM node:26-bookworm-slim

WORKDIR /app

COPY package*.json ./
RUN npm ci

# Download chromium browser for Playwright. Playwright is used by Vitest Browser Mode for testing
RUN npx playwright install --with-deps chromium

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev"]
