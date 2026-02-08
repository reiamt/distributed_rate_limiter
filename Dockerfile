# use go binary
FROM golang:1.25-alpine AS build

WORKDIR /app

# build stage
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -o server ./cmd/server 

# run stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server .

ENTRYPOINT ["./server"]