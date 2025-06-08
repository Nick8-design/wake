package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func checkStatus() (string, error) {
	url := "https://dnd-ke.onrender.com/ping"

	
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}



	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server down, status: %d", resp.StatusCode)
	}

	// var result map[string]string
	var ping string
	if err := json.NewDecoder(resp.Body).Decode(&ping); err != nil {
		return "", err
	}

	return ping, nil
}

func startHealthCheck() {
	go func() {
		for {
			ping, err := checkStatus()
			if err != nil {
				log.Printf("[HealthCheck] Error: %v. Retrying in 30 seconds...", err)
				time.Sleep(30 * time.Second)
			} else {
				log.Printf("[HealthCheck] Success: %v. Next check in 30 minutes.", ping)
				time.Sleep(30 * time.Minute)
			}
		}
	}()
}

func main() {
	app := fiber.New()

	
	startHealthCheck()


	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running...")
	})

	log.Fatal(app.Listen(":21886"))
}
