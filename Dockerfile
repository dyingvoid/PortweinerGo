FROM golang:1.23

WORKDIR /usr/src/app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application code
COPY . .

# Create the build directory if it doesn't exist
RUN mkdir -p ./build

# Copy the necessary files into the build directory
COPY ./build/static ./build/static
COPY ./build/config ./build/config

# Build the Go application
RUN go build -v -o ./build/app ./cmd/main.go

# Change the working directory to the build directory
WORKDIR /usr/src/app/build

# Set the default command to run the application
CMD ["./app", "-config-path", "./config/local.yaml"]