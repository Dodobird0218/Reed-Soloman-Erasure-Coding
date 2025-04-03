package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// MessageData represents the JSON structure of the input message
type MessageData struct {
	Message []string `json:"message"`
}

// EncodedData represents the JSON structure of the output encoding result
type EncodedData struct {
	Message []string `json:"message"`
	Encoded []string `json:"encoded"`
}

// ReadMessageFromJSON reads message data from a JSON file
func ReadMessageFromJSON(filePath string) ([]byte, error) {
	// Read JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	// Parse JSON
	var messageData MessageData
	if err := json.Unmarshal(data, &messageData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Convert hex strings to byte array
	message := make([]byte, len(messageData.Message))
	for i, hexStr := range messageData.Message {
		// Remove possible "0x" prefix
		hexStr = strings.TrimPrefix(hexStr, "0x")
		// Parse hex string
		value, err := strconv.ParseUint(hexStr, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to parse hex value: %v", err)
		}
		message[i] = byte(value)
	}

	return message, nil
}

// WriteEncodedToJSONWithOriginal writes the encoding result and original message to a JSON file
func WriteEncodedToJSONWithOriginal(filePath string, encoded []byte, originalMessage []byte) error {
	// Convert byte arrays to hex string arrays
	encodedStrings := make([]string, len(encoded))
	messageStrings := make([]string, len(originalMessage))

	for i, b := range encoded {
		encodedStrings[i] = fmt.Sprintf("0x%02x", b)
	}

	for i, b := range originalMessage {
		messageStrings[i] = fmt.Sprintf("0x%02x", b)
	}

	// Create output data structure
	outputData := EncodedData{
		Message: messageStrings,
		Encoded: encodedStrings,
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON encoding failed: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// WriteEncodedToJSON writes the encoding result to a JSON file, using the first 6 bytes of the encoding result as the original message
func WriteEncodedToJSON(filePath string, encoded []byte) error {
	// Assume the original message is the first 6 bytes of the encoded result
	originalMessage := encoded[:min(6, len(encoded))]
	return WriteEncodedToJSONWithOriginal(filePath, encoded, originalMessage)
}

// min returns the minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
