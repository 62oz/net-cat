package main

import (
	"fmt"
	"net-cat/server"
	"os"
	"strconv"
)

func main() {
	port := ""
	if len(os.Args) < 2 {
		fmt.Println("Default port :8989")
		port = "8989"
	} else if len(os.Args) == 2 {
		port = os.Args[1]
		if _, err := strconv.Atoi(port); err != nil {
			fmt.Println("Invalid port.")
			return
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	server.Server(port)
}
