package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Use a constant cost of 10 which is the default
	cost := 10
	
	// Make sure password is provided as argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_hash.go <password>")
		os.Exit(1)
	}
	
	password := os.Args[1]
	
	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Password: %s\nHash: %s\n", password, string(hash))
}
