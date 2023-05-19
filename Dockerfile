FROM golang:latest as build
WORKDIR /app
COPY main.go ./
RUN go build main.go

FROM alpine:latest as exec
WORKDIR /app
COPY --from=build /app/main ./
CMD ["./main"]