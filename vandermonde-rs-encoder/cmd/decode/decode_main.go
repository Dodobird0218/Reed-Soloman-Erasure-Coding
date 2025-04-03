package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"rs-encoder/gf"
	"rs-encoder/rs"
	"strings"
)

// MessageData struct for parsing JSON input
type MessageData struct {
	Message    []string `json:"message"`
	StartIndex int      `json:"start_index,omitempty"` // Optional starting index for shards
}

// DecodedData struct for generating output JSON
type DecodedData struct {
	EncodedShards []string `json:"encoded_shards"`
	DecodedData   []string `json:"decoded_data"`
}

func main() {
	// Check command line arguments
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <input_file> <output_file>\n", os.Args[0])
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Initialize finite field GF(2^8)
	field := gf.NewGF(0x1D) // Using GF(2^8) finite field, polynomial x^4 + x^3 + x^2 + 1

	// Read input from specified JSON file
	messageData, err := readMessageFromJSON(inputFile)
	if err != nil {
		fmt.Printf("Cannot read input file: %v\n", err)
		return
	}

	// Convert hexadecimal strings to byte array
	encodedShards := hexStringsToBytes(messageData.Message)
	dataShards := 6   // Number of original data shards
	totalShards := 18 // Total number of shards (original data + redundancy)

	// Create Vandermonde Reed-Solomon decoder
	decoder := rs.NewVandermondeDecoder(field, dataShards, totalShards)

	// Set shard indices based on specified starting index
	indices := make([]int, len(encodedShards))
	startIndex := messageData.StartIndex // Get starting index from JSON
	for i := 0; i < len(encodedShards); i++ {
		indices[i] = startIndex + i
	}

	// Print input encoded shards and indices
	fmt.Println("Input encoded shards:")
	printArray(encodedShards)
	fmt.Println("Used shard indices:", indices)

	// Decode
	decodedData := decoder.Decode(encodedShards, indices)

	// Print decoding result
	fmt.Println("\nDecoding result (original message):")
	printArray(decodedData)

	// Create output JSON structure
	outputData := DecodedData{
		EncodedShards: messageData.Message,
		DecodedData:   bytesToHexStrings(decodedData),
	}

	// Save to specified output file
	err = saveToJSON(outputFile, outputData)
	if err != nil {
		fmt.Printf("Cannot save output file: %v\n", err)
		return
	}

	fmt.Println("\nDecoding result saved to", outputFile)
}

// Read message from JSON file
func readMessageFromJSON(filename string) (MessageData, error) {
	var data MessageData
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(fileContent, &data)
	return data, err
}

// Save result to JSON file
func saveToJSON(filename string, data DecodedData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, jsonData, 0644)
}

// Convert hexadecimal string array to byte array
func hexStringsToBytes(hexStrings []string) []byte {
	bytes := make([]byte, len(hexStrings))
	for i, hexStr := range hexStrings {
		// Remove "0x" prefix
		cleanHex := strings.TrimPrefix(hexStr, "0x")

		// Parse hexadecimal string
		var value byte
		fmt.Sscanf(cleanHex, "%x", &value)
		bytes[i] = value
	}
	return bytes
}

// Convert byte array to hexadecimal string array
func bytesToHexStrings(bytes []byte) []string {
	hexStrings := make([]string, len(bytes))
	for i, b := range bytes {
		hexStrings[i] = fmt.Sprintf("0x%02x", b)
	}
	return hexStrings
}

// Print array
func printArray(array []byte) {
	fmt.Print("[ ")
	for i, val := range array {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(val)
	}
	fmt.Println(" ]")

	// Print in hexadecimal format
	fmt.Print("Hexadecimal: [ ")
	for i, val := range array {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("0x%02x", val)
	}
	fmt.Println(" ]")
}
