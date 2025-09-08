package template

const Dockerfile = `# Use the official Golang image
FROM golang:1.25-alpine as builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

RUN mkdir "gen"

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

COPY . .

RUN oapi-codegen --config ./config/config.yaml ./api/swagger.yaml

# Download dependencies
RUN go mod tidy

# Copy the entire source code

# Build the Go app
RUN go build -o app cmd/server/main.go

FROM alpine:latest

LABEL mainers=Manteiner

RUN apk add --no-cache curl

RUN mkdir /opt/resources/

COPY --from=builder app/app /bin/app
COPY --from=builder app/api/swagger.yaml /opt/resources/swagger.yaml

ENV PORT 8080
ENV CONTEXT_PATH /api/${PROJECT}/v1/
ENV SWAGGER_FILE /opt/resources/swagger.yaml

EXPOSE $PORT

CMD ["app"]
`

const DockerCompose = `services:
  ${PROJECT}:
    build:
      dockerfile: Dockerfile
      context: .
      args:
        PORT: 8080
    image: ${PROJECT}
    container_name: ${PROJECT}
    restart: always
    deploy:
      resources:
        limits:
          cpus: "1"
          memory: "256M"
    healthcheck:
      test: [ "CMD", "curl", "-f", "http://172.17.0.1:8080/api/${PROJECT}/v1/health" ]
      interval: 30s
      timeout: 15s
      retries: 5
      start_period: 10s
    ports:
      - "8080:8080"
    volumes:
      - ./.log/container/:/.log
`
