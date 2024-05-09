package main

import "fmt"

func main() {
	server := NewHTTPServer(":8080")
	fmt.Println("Server started")
	server.ListenAndServe()
}
