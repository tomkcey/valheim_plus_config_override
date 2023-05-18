FROM golang:latest as build
WORKDIR /app
COPY main.go migrations dbconfig.yml ./
RUN ["go", "build","main.go"]
RUN ["go", "install", "github.com/rubenv/sql-migrate/...@latest"]
RUN ["sql-migrate", "up"]

FROM alpine:latest as exec
WORKDIR /app
COPY --from=build /app/main ./
CMD ["./main"]