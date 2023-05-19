# FROM golang:latest as build
# WORKDIR /app
# COPY main.go input processor ./
# RUN ["export", "PATH=$PATH:/usr/local/go/bin"]
# RUN ["export", "GOROOT=/usr/local/go"]
# RUN ["go", "mod", "init", "lib"]
# RUN ["go", "build","main.go"]

FROM alpine:latest as build
WORKDIR /app
COPY main.go input processor ./
RUN ["apt install curl"]
RUN curl https://go.dev/dl/go1.20.4.linux-amd64.tar.gz --output /usr/local/go1.20.4.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin
RUN go version
RUN export GOROOT=/usr/local/go
RUN ["go", "mod", "init", "lib"]
RUN ["go", "build","main.go"]

FROM alpine:latest as exec
WORKDIR /app
COPY --from=build /app/main ./
CMD ["./main"]