package main

import (
	"errors"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("notes.txt", os.O_APPEND, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist){
			f, _ = os.Create("notes.txt")
			f.Write([]byte("header"))
		}
	}

	
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}