// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/blah/net"
	"github.com/jteeuwen/blah/proto"
	"log"
	"os"
	"path/filepath"
)

var config *Config

func main() {
	conn, client := setup()

	defer func() {
		if x := recover(); x != nil {
			log.Printf("Error: %x", x)
		}

		shutdown(conn, client)
	}()

	// Perform handshake.
	log.Printf("Performing handshake...")
	client.Login(config.Nickname)
	client.Nick(config.Nickname, config.NickservPassword)

	// Main data loop.
	log.Printf("Entering data loop...")
	for {
		line, err := conn.ReadLine()

		if err != nil {
			break
		}

		client.Read(string(line))
	}
}

// setup initializes the application.
func setup() (*net.Conn, *proto.Client) {
	config = parseArgs()

	// Open connection to server.
	log.Printf("Connecting to %s...", config.Address)
	conn, err := net.Dial(config.Address, config.SSLKey, config.SSLCert)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Dial: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Connection established.")

	// Create protocol handler.
	client := proto.NewClient(func(p []byte) error {
		log.Printf("< %s", p)
		_, err := conn.Write(p)
		return err
	})

	Bind(client)
	return conn, client
}

// shutdown cleans up our mess.
func shutdown(conn *net.Conn, client *proto.Client) {
	log.Printf("Shutting down.")
	client.Quit(config.QuitMessage)
	client.Close()
	conn.Close()
}

// parseArgs reads and verfies commandline arguments.
// It loads and returns a configuration object.
func parseArgs() *Config {
	config := flag.String("c", "", "Path to bot configuraiton file.")
	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if len(*config) == 0 {
		fmt.Fprintf(os.Stderr, "Missing configuraiton file.\n")
		flag.Usage()
		os.Exit(1)
	}

	var c Config
	err := c.Load(filepath.Clean(*config))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Load config: %v\n", err)
		os.Exit(1)
	}

	return &c
}
