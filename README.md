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
   docker compose up --build
   ```

3. Access the API endpoints through `http://localhost:8080`.

# Current Tasks
- [x] Implement server start functionality
- [x] Add server stop functionality
- [ ] Add file upload/download capabilities
- [ ] Add server configuration management
- [ ] Add file management features
- [ ] Implement server logs viewing
- [ ] Add server status monitoring
- [ ] Implement server monitoring features
- [ ] Add user authentication and authorization
- [ ] Implement server backup functionality
- [ ] Add a minimal "lobby" server that users are redirected to when the server gets shutdown / restarted
 
# Docs
API documentation is served via Swagger UI at `http://localhost:8080/swagger/index.html`.

Docs are generated from comment annotations on handlers using [swaggo](https://github.com/swaggo/swag). To regenerate after editing annotations:
```bash
swag init
```
