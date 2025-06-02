# Build stage
FROM node:22-alpine AS build

# Install pnpm
RUN npm install -g pnpm@10.11.0

# Set working directory
WORKDIR /app

# Copy package.json files
COPY package.json pnpm-lock.yaml ./

# Install dependencies - don't need --frozen-lockfile as it is assumed in the install command when a lock file is present
RUN pnpm install 

# Copy all files
COPY . .

# Build the application
RUN pnpm build

# Production stage
FROM nginx:alpine AS production

# Copy built assets from the build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Add runtime configuration script
RUN echo '#!/bin/sh' > /docker-entrypoint.d/40-config-runtime.sh && \
    echo 'envsubst < /usr/share/nginx/html/assets/env-config.template.js > /usr/share/nginx/html/assets/env-config.js' >> /docker-entrypoint.d/40-config-runtime.sh && \
    chmod +x /docker-entrypoint.d/40-config-runtime.sh

# Create a template file for runtime environment variables
RUN mkdir -p /usr/share/nginx/html/assets
RUN echo "window.ENV = { \
  API_URL: '${API_URL:-/api}', \
  AUTH_URL: '${AUTH_URL:-/auth}' \
};" > /usr/share/nginx/html/assets/env-config.template.js && \
    echo "" > /usr/share/nginx/html/assets/env-config.js

# Expose port 80
EXPOSE 80

# No need to specify CMD as it's inherited from the nginx base image
