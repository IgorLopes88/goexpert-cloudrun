FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/cloud-run main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin/cloud-run .
ENTRYPOINT [ "./cloud-run" ]
