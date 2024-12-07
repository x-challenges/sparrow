package main

import (
	_ "github.com/joho/godotenv/autoload"

	"sparrow/cmd"
)

func main() {
	cmd.Execute()
}
