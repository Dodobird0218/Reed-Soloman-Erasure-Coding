package rs

import (
	"fmt"
	"rs-encoder/gf"
)

// RSEncoder Reed-Solomon encoder
type RSEncoder struct {
	field             *gf.GF
	dataShards        int
	parityShards      int
	totalShards       int
	alphaPoints       []byte
	vandermondeMatrix [][]byte
}

// NewRSEncoder creates a new Reed-Solomon encoder
func NewRSEncoder(field *gf.GF, dataShards, parityShards int) *RSEncoder {
	encoder := &RSEncoder{
		field:        field,
		dataShards:   dataShards,
		parityShards: parityShards,
		totalShards:  dataShards + parityShards,
	}
	encoder.generateAlphaPoints()
	encoder.generateVandermondeMatrix()
	return encoder
}

// Encode encodes the message using Reed-Solomon encoding with Vandermonde matrix
func (enc *RSEncoder) Encode(message []byte) []byte {
	if len(message) != enc.dataShards {
		panic("Message length must equal the number of data shards")
	}

	// Create encoded result array, first dataShards items are same as original message (systematic encoding)
	encoded := make([]byte, enc.totalShards)
	copy(encoded, message)

	// Use the correctly implemented Vandermonde encoding method
	enc.vandermondeEncode(message, encoded)

	return encoded
}

// printVandermondeMatrix prints the Vandermonde matrix for debugging
func (enc *RSEncoder) printVandermondeMatrix() {
	fmt.Println("\nVandermonde Matrix:")
	for i := 0; i < enc.parityShards; i++ {
		fmt.Printf("Row %d: [", i)
		for j := 0; j < enc.dataShards; j++ {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("0x%02x", enc.vandermondeMatrix[i][j])
		}
		fmt.Println("]")
	}
}

// generateAlphaPoints generates alpha evaluation points
func (enc *RSEncoder) generateAlphaPoints() {
	enc.alphaPoints = make([]byte, enc.totalShards)

	// Use consecutive integers as evaluation points (1, 2, 3, ...)
	for i := 0; i < enc.totalShards; i++ {
		enc.alphaPoints[i] = byte(i + 1)
	}

	// Print evaluation points
	fmt.Println("Vandermonde evaluation points (alpha points):")
	for i, point := range enc.alphaPoints {
		fmt.Printf("  Point[%d] = 0x%02x\n", i, point)
	}
}

// generateVandermondeMatrix generates the Vandermonde matrix for encoding
func (enc *RSEncoder) generateVandermondeMatrix() {
	// Create Vandermonde matrix (including parity rows)
	enc.vandermondeMatrix = make([][]byte, enc.parityShards)

	// Only need to calculate the matrix for parity data rows
	for i := 0; i < enc.parityShards; i++ {
		// Get evaluation point
		x := enc.alphaPoints[i+enc.dataShards]

		enc.vandermondeMatrix[i] = make([]byte, enc.dataShards)

		// First element is x^0 = 1
		enc.vandermondeMatrix[i][0] = 1

		// For each subsequent column, calculate x^j
		for j := 1; j < enc.dataShards; j++ {
			// Calculate x^j
			enc.vandermondeMatrix[i][j] = enc.field.Pow(x, j)
		}
	}
}

// vandermondeEncode uses Vandermonde matrix to calculate parity data
func (enc *RSEncoder) vandermondeEncode(message []byte, encoded []byte) {
	// For each parity position
	for i := enc.dataShards; i < enc.totalShards; i++ {
		// Calculate polynomial value at this point using Lagrange interpolation
		result := byte(0)

		// Construct Lagrange interpolation polynomial
		for j := 0; j < enc.dataShards; j++ {
			// Get the message value
			y_j := message[j]

			// Skip if the value is 0 (optimization)
			if y_j == 0 {
				continue
			}

			// Calculate Lagrange basis L_j(x)
			basis := byte(1)

			for k := 0; k < enc.dataShards; k++ {
				if j != k {
					// Calculate (x - x_k)
					numerator := enc.field.Sub(enc.alphaPoints[i], enc.alphaPoints[k])
					// Calculate (x_j - x_k)
					denominator := enc.field.Sub(enc.alphaPoints[j], enc.alphaPoints[k])
					// Division
					factor := enc.field.Div(numerator, denominator)
					// Multiply by the current basis
					basis = enc.field.Mul(basis, factor)
				}
			}

			// Calculate this term's contribution: y_j * L_j(x)
			term := enc.field.Mul(y_j, basis)

			// Add to the result
			result = enc.field.Add(result, term)
		}

		encoded[i] = result
		fmt.Printf("Vandermonde encoding at position %d: 0x%02x\n", i, result)
	}
}
