# MC Manager
MC Manager is a tool for managing Minecraft servers, built with my homelab in mind. It allows you to easily create, manage, and monitor minecraft servers from a web interface or command line.

Web interface not yet implemented. Will be added once the API is in a reasonably ready state.

## Tech Requirements / Stack
- Go 1.23+
- Docker

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/lomokwa/mc-manager.git
   cd mc-manager
   ```

2. Start the application using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Access the API endpoints through `http://localhost:8080`.

# Current Tasks
- [ ] Implement server start functionality
- [ ] Add server stop functionality
- [ ] Add file upload/download capabilities
- [ ] Add server configuration management
- [ ] Add file management features
- [ ] Implement server logs viewing
- [ ] Add server status monitoring
- [ ] Implement server monitoring features
- [ ] Add user authentication and authorization
- [ ] Implement server backup functionality

# API Endpoints
- `POST /start`: Start a Minecraft server.
