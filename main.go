// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/ircb/net"
	"github.com/jteeuwen/ircb/plugins/url"
	"github.com/jteeuwen/ircb/proto"
	"log"
	"os"
	"path/filepath"
)

func main() {
	conn, client := setup()
	defer shutdown(conn, client)

	// Bind protocol handlers and commands.
	bind(client)

	// Initialize plugins.
	url.Init(config.Profile, client)

	// Perform handshake.
	log.Printf("Performing handshake...")
	client.User(config.Nickname)
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
	// parse commandline arguments and create configuration.
	config = parseArgs()

	log.Printf("Connecting to %s...", config.Address)

	// Open connection to server.
	conn, err := net.Dial(config.Address, config.SSLCert, config.SSLKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dial: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connection established.")

	// Create client protocol.
	client := proto.NewClient(func(p []byte) error {
		_, err := conn.Write(p)
		return err
	})

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
	profile := flag.String("p", "", "Path to bot profile directory.")
	version := flag.Bool("v", false, "Display version information.")

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if len(*profile) == 0 {
		fmt.Fprintf(os.Stderr, "Missing profile directory.\n")
		flag.Usage()
		os.Exit(1)
	}

	var c Config
	c.Profile = filepath.Clean(*profile)

	err := c.Load(filepath.Join(c.Profile, "config.ini"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load config: %v\n", err)
		os.Exit(1)
	}

	return &c
}
