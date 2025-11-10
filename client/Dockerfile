# Stage 1: Build the Vue application
FROM node:20-alpine AS builder

# Set working directory for client code
WORKDIR /app/client

# Copy package files and install dependencies
COPY ./package.json ./package-lock.json* ./
RUN npm install

# Copy the rest of the client source code
COPY . .

# Copy the .env file if your build process needs it (adjust path if needed)
# COPY .env .env

# Build the application
# Ensure your build script outputs to './dist' inside the WORKDIR
RUN npm run build

# Stage 2: Serve the built app with Nginx
FROM nginx:alpine

# Remove default nginx website conf
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx conf (expects nginx.conf in the build context root)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder stage to nginx html directory
COPY --from=builder /app/client/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
