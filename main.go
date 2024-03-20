package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	z85 "github.com/braheezy/z85/pkg"
)

func main() {
	// Define flags
	mode := flag.String("mode", "encode", "Mode of operation, 'encode' or 'decode'.")
	flag.Parse()

	var inputData []byte
	var err error
	if flag.NArg() > 0 { // Data is provided through command-line arguments
		argData := strings.Join(flag.Args(), " ")
		if *mode == "encode" {
			// Convert hex string to bytes for encoding
			inputData, err = hex.DecodeString(argData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding hex input: %v\n", err)
				os.Exit(1)
			}
		} else {
			// For decoding, the input is treated directly as a string
			inputData = []byte(argData)
		}
	} else if isInputFromPipe() { // Data is provided through a pipe
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		inputStr := scanner.Text()
		if *mode == "encode" {
			// Convert hex string to bytes for encoding
			inputData, err = hex.DecodeString(inputStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error decoding hex input: %v\n", err)
				os.Exit(1)
			}
		} else {
			inputData = []byte(inputStr)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: z85 -mode=[encode|decode] \"data\"")
		fmt.Println("Or: echo \"data\" | z85 -mode=[encode|decode]")
		os.Exit(1)
	}

	// Perform the operation based on the mode
	switch *mode {
	case "encode":
		encoded, err := z85.Encode(inputData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(encoded)
	case "decode":
		decoded, err := z85.Decode(string(inputData))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding data: %v\n", err)
			os.Exit(1)
		}
		if _, err := os.Stdout.Write(decoded); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing decoded data: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode specified: %s\n", *mode)
		os.Exit(1)
	}
}

// isInputFromPipe checks if there is data being piped into the stdin.
func isInputFromPipe() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking stdin: %v\n", err)
		os.Exit(1)
	}
	return fileInfo.Mode()&os.ModeCharDevice == 0
}
