# Stage 1: Build the React application
FROM node:22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package.json and pnpm-lock.yaml to leverage Docker cache
COPY package.json ./
COPY pnpm-lock.yaml ./

# Install pnpm and then frontend dependencies
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Copy the rest of the frontend source code
COPY . .

# Build the React application
RUN pnpm build

# Stage 2: Serve the application with Nginx
FROM nginx:alpine

# Copy the built React app from the builder stage into Nginx's public directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Remove default Nginx configuration
RUN rm /etc/nginx/conf.d/default.conf

# Copy a custom Nginx configuration for our React app
COPY nginx.conf /etc/nginx/conf.d/iotgo.conf

# Expose port 80 for the web server
EXPOSE 3000

# Command to start Nginx
CMD ["nginx", "-g", "daemon off;"]
