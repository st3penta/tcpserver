package main

func main() {

	port := 5555
	server := NewServer(port)

	server.Start()
}
