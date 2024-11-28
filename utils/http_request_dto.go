package utils

type RedisCommandRequest struct {
	Command string   `json:"command" binding:"required"`
	Args    []string `json:"args"`
}
