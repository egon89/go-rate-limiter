# Stage 1: build application
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# the -C flag is used to specify the directory where the go build command should be executed from the root of the project
# the binary will be built in the cmd directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C cmd -o goratelimiter

# Stage 2: production go image
FROM scratch
WORKDIR /app

COPY --from=builder /app/cmd/goratelimiter .

EXPOSE 8080


ENTRYPOINT [ "./goratelimiter" ]