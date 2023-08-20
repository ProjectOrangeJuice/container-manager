package main

import (
	"bufio"
	"container-manager/client/storage"
	"container-manager/shared"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	// Define the server address and port.
	addr := "localhost:8080"

	// Create a TCP connection to the server.
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Error dialing: %s", err)
		return
	}

	fmt.Fprintf(conn, "Test client\n")
	log.Print("Connected to server")

	// Create a buffered reader
	reader := bufio.NewReader(conn)

	for {
		// Read a line of data
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		// Print the line
		fmt.Println(line)
		readLine(line, conn)
	}

	storages, err := storage.GetFreeStorageSpace()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("%+v", storages)
}

func readLine(line string, conn net.Conn) {
	switch strings.TrimSpace(line) {
	case "STORAGE_INFO":
		sendBackStorage(conn)
	}
}

func sendBackStorage(conn net.Conn) {
	log.Print("Sending storage info")
	storages, err := storage.GetFreeStorageSpace()
	if err != nil {
		log.Printf("Error getting storage info, %s", err)
		return
	}

	out, err := createEvent("STORAGE", storages)
	if err != nil {
		log.Printf("Error creating event, %s", err)
		return
	}
	fmt.Fprintf(conn, "%s\n", out)
	log.Printf("Sent storage info %s", out)
}

// A generic function that creates an event.
func createEvent[R any](request string, result R) ([]byte, error) {
	resultByte, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("could not marshal result, %s", err)
	}

	evt := shared.EventData{
		Request: request,
		Result:  resultByte,
	}
	eventOut, err := json.Marshal(evt)
	if err != nil {
		return nil, fmt.Errorf("could not marshal event, %s", err)
	}
	return eventOut, nil
}
