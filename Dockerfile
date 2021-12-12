##### Stage 1 #####

### Use golang:1.15 as base image for building the application
FROM tetafro/golang-gcc:1.17-alpine as builder

### Create new directly and set it as working directory
RUN mkdir -p /app
WORKDIR /app

### Copy Go application dependency files
COPY go.mod .
COPY go.sum .

### Download Go application module dependencies
RUN go mod download

### Copy actual source code for building the application
COPY . .

### Build the Go app for a linux OS
### 'scratch' and 'alpine' both are Linux distributions
RUN go build -o app Nimie_alpha/main

### Define the running image
#FROM scratch

### Alternatively to 'FROM scratch', use 'alpine':
FROM alpine:3.14

### Set working directory
WORKDIR /app

### Copy built binary application from 'builder' image
COPY --from=builder /app .

EXPOSE 8080
### Run the binary application
CMD ["./app"]