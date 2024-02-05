package main

import (
	"Cloudflare-DDNS/src"
	"fmt"
	"os"
)

func main() {
	src.Logger, src.LogFile, _ = src.CreateLogger()
	// Read the config

	_, err := src.ReadConfig()
	fmt.Println(err)
	if err != nil {
		src.Logger.Fatalln("Error reading the config file:", err)
		return
	}
	src.Logger.Println("Config file read successfully")
	// Start the server
	src.Logger.Println("Starting the server")
	err = src.StartServer()
	if err != nil {
		src.Logger.Fatalln("Error running the server:", err)
		return
	}

	defer func(LogFile *os.File) {
		err := LogFile.Close()
		if err != nil {
			return
		}
	}(src.LogFile)
}
