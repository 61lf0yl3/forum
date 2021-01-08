FROM golang:1.13.5

# Set the Current Working Directory inside the container
WORKDIR /app/forum

# Populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Build the Go app
RUN go build -v cmd/app/main.go

# Run the binary program produced by `go install`
CMD ["./main"]