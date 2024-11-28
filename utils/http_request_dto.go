package utils

type RedisCommandRequest struct {
	Command string        `json:"command"`
	Args    []interface{} `json:"args"`
}
