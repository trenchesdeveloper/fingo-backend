package main

import "github.com/trenchesdeveloper/fingo-backend/api"

func main() {
	server := api.NewServer(".")

	server.Start(3000)

}
