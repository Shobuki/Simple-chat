FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
ENV PORT=8080
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
