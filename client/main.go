package main

import (
	"container-manager/client/storage"
	"log"
)

func main() {

	storages, err := storage.GetFreeStorageSpace()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("%+v", storages)
}
