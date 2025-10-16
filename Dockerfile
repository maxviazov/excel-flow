FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/testdata/sample.xlsx ./testdata/sample.xlsx
COPY --from=builder /app/testdata/drivers_summary.xlsx ./testdata/drivers_summary.xlsx
RUN mkdir -p uploads outputs

EXPOSE 8080
CMD ["./server"]
