package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// List of API URLs to ping
var apiURLs = []string{
	// "http://127.0.0.1:21886/ping",
	"https://dnd-ke.onrender.com/ping",
	"https://mpesaapi.onrender.com",
	"https://residential-tracker-and-booking-images.onrender.com",
	"https://finder-site.onrender.com",
	"https://uploadapi-8244.onrender.com",
	"https://service-2p6f.onrender.com",

}

// Check the status of a single URL
func checkStatus(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server down, status: %d", resp.StatusCode)
	}

	var ping string
	if err := json.NewDecoder(resp.Body).Decode(&ping); err != nil {
		return "", err
	}

	return ping, nil
}

// Ping all URLs in parallel
func startHealthCheck() {
	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(len(apiURLs))

			for _, url := range apiURLs {
				go func(u string) {
					defer wg.Done()
					ping, err := checkStatus(u)
					if err != nil {
						log.Printf("[HealthCheck] %s Error: %v", u, err)
					} else {
						log.Printf("[HealthCheck] %s Success: %v", u, ping)
					}
				}(url)
			}

			wg.Wait()
			log.Println("[HealthCheck] All checks done. Next run in 30 minutes...")
			time.Sleep(25 * time.Minute)
		}
	}()
}

func main() {
	app := fiber.New()

	startHealthCheck()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Health check service is running...")
	})


	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON("Pong")
	})

	log.Fatal(app.Listen(":21887"))
}
