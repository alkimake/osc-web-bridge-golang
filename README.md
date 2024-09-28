# Multicast Data Monitor

This application is a WebSocket-based multicast data monitor built with Go. It allows you to receive and display multicast messages in real-time through a web interface.

## Features

- Real-time data display using WebSocket.
- Simple and clean user interface.
- Built with Go and Nix for development environment management.

## Prerequisites

- Go 1.22 or higher
- **OPTIONAL** Nix package manager

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/alkimake/osc-web-bridge-golang.git
cd osc-web-bridge-golang
```

### Using Nix Flakes

To set up the development environment using Nix flakes, run the following command:

```bash
nix develop
```

This will create a development shell with all necessary dependencies.

### Running the Application

1. Start the Go server:

```bash
go run ./cmd
```

2. Open your web browser with index.htm under web folder

3. Send sample message by using osc_client.go:

```bash
go run ./cmd/osc_client.go
```

### Using Just

If you prefer using Just for task management, you can run the application with:

```
just run
```

sending message will be;

```bash
just osc_client
```


## License

This project is licensed under the MIT License.
