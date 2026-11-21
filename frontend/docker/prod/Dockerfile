# Build stage
# ---------------------
FROM node:26.13-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# App runtime paramenters must be supplied as build args since they are baked into the build
ARG VITE_API_URL
RUN if [ -z "$VITE_API_URL" ]; then \
      echo "ERROR: VITE_API_URL build-arg is required" >&2; \
      exit 1; \
    fi
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

# Run stage
# ---------------------
FROM nginx:1.27-alpine

COPY docker/prod/nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=builder /app/dist /usr/share/nginx/html

# Modify permissions so non-root user (nginx) can be used
RUN mkdir -p /tmp/nginx \
    && chown -R nginx:nginx /usr/share/nginx/html /var/cache/nginx /var/log/nginx /tmp/nginx \
    && sed -i 's|pid\s\+/run/nginx.pid;|pid /tmp/nginx/nginx.pid;|' /etc/nginx/nginx.conf \
    && sed -i 's|^user\s\+nginx;|# user directive removed for non-root|' /etc/nginx/nginx.conf

USER nginx

EXPOSE 8080

# Start nginx to serve app files
CMD ["nginx", "-g", "daemon off;"]
