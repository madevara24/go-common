package server

import "time"

type Config struct {
	Address        string
	Env            string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	AllowedOrigins string
}
