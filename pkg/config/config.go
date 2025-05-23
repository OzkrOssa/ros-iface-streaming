package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App  *App
		Http *Http
		Ros  *RouterOS
	}
	App struct {
		Name string
		Env  string
	}
	Http struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}
	RouterOS struct {
		Port string
		User string
		Pass string
	}
)

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "prod" {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	http := &Http{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	ros := &RouterOS{
		Port: os.Getenv("ROUTEROS_PORT"),
		User: os.Getenv("ROUTEROS_USER"),
		Pass: os.Getenv("ROUTEROS_PASS"),
	}
	return &Container{
		App:  app,
		Http: http,
		Ros:  ros,
	}, nil
}
