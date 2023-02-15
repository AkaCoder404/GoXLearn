package main

import "main/app/service"

var (
	Version string = "" // exposed globally (uppercase)
)

// Initialize the Program
func main() {
	// Initialize the application
	service.Start()
	defer service.Shutdown()
}
