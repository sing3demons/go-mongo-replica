FROM golang:1.20-alpine3.16 as step_1
WORKDIR /app
COPY go.* .
RUN go mod tidy
COPY . .
RUN go build -o main

FROM alpine:3.16.0
WORKDIR /app
RUN apk add ca-certificates
COPY --from=step_1 /app/main .
EXPOSE 8080
CMD ["/app/main"]