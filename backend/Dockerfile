FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /member-api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache postgresql-client && adduser -D appuser
USER appuser
COPY --from=build /member-api /member-api
EXPOSE 8080
ENTRYPOINT ["/member-api"]
