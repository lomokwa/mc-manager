package types

type StartServerRequest struct {
	CreateLaunchScript  bool              `json:"createLaunchScript"`
	ConfigureProperties bool              `json:"configureProperties"`
	Properties          map[string]string `json:"properties"`
}

type UpdateServerPropertiesRequest struct {
	Properties map[string]string `json:"properties"`
}

type RegisterRequest struct {
	Token    string `json:"token" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
