FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o go-api .

FROM scratch
COPY --from=builder /app/go-api /go-api
USER 1001:1001
EXPOSE 8080
ENTRYPOINT ["/go-api"]
