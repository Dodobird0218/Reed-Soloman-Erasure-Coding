package rs

import (
	"rs-encoder/gf"
)

// RSEncoder2 Reed-Solomon encoder with consecutive integers as evaluation points
type RSEncoder2 struct {
	field        *gf.GF
	dataShards   int
	parityShards int
	totalShards  int
	evalPoints   []byte
}

// NewRSEncoder2 creates a new Reed-Solomon encoder using consecutive integer evaluation points
func NewRSEncoder2(field *gf.GF, dataShards, parityShards int) *RSEncoder2 {
	encoder := &RSEncoder2{
		field:        field,
		dataShards:   dataShards,
		parityShards: parityShards,
		totalShards:  dataShards + parityShards,
	}
	encoder.generateEvalPoints()
	return encoder
}

// generateEvalPoints generates the evaluation points using consecutive integers
func (enc *RSEncoder2) generateEvalPoints() {
	enc.evalPoints = make([]byte, enc.totalShards)

	// Use consecutive integers as evaluation points (starting from 1)
	for i := 0; i < enc.totalShards; i++ {
		enc.evalPoints[i] = byte(i + 1)
	}
}

// Encode encodes the message using Reed-Solomon encoding
func (enc *RSEncoder2) Encode(message []byte) []byte {
	if len(message) != enc.dataShards {
		panic("Message length must equal the number of data shards")
	}

	// Create encoded result array, first dataShards items are same as original message
	encoded := make([]byte, enc.totalShards)
	copy(encoded, message)

	// Use Lagrange interpolation to calculate redundant data
	enc.lagrangeInterpolation(message, encoded)

	return encoded
}

// lagrangeInterpolation calculates redundant shards using Lagrange interpolation
func (enc *RSEncoder2) lagrangeInterpolation(message []byte, encoded []byte) {
	// For each redundant shard position
	for i := enc.dataShards; i < enc.totalShards; i++ {
		// Calculate the value of the polynomial at this point
		result := byte(0)

		// Build the Lagrange interpolation polynomial
		for j := 0; j < enc.dataShards; j++ {
			term := message[j]

			// Calculate the Lagrange basis function
			for k := 0; k < enc.dataShards; k++ {
				if j != k {
					// Calculate (x - x_k)
					numerator := enc.field.Sub(enc.evalPoints[i], enc.evalPoints[k])
					// Calculate (x_j - x_k)
					denominator := enc.field.Sub(enc.evalPoints[j], enc.evalPoints[k])
					// Division
					factor := enc.field.Div(numerator, denominator)
					// Multiply by the current term
					term = enc.field.Mul(term, factor)
				}
			}

			// Add this term to the result
			result = enc.field.Add(result, term)
		}

		encoded[i] = result
	}
}

// EncodeEfficient is a more efficient implementation using Horner's method
func (enc *RSEncoder2) EncodeEfficient(message []byte) []byte {
	if len(message) != enc.dataShards {
		panic("Message length must equal the number of data shards")
	}

	// Create encoded result array, first dataShards items are same as original message
	encoded := make([]byte, enc.totalShards)
	copy(encoded, message)

	// For each redundant shard position
	for i := enc.dataShards; i < enc.totalShards; i++ {
		result := byte(0)

		// Get the evaluation point
		x := enc.evalPoints[i]

		// Use Horner's method to calculate the polynomial value
		// p(x) = message[0] + message[1]*x + message[2]*x^2 + ... + message[dataShards-1]*x^(dataShards-1)
		for j := enc.dataShards - 1; j >= 0; j-- {
			result = enc.field.Add(enc.field.Mul(result, x), message[j])
		}

		encoded[i] = result
	}

	return encoded
}

// ReconstructData reconstructs the original data from any combination of data and parity shards
// This is an additional method to demonstrate the full capability of Reed-Solomon codes
func (enc *RSEncoder2) ReconstructData(availableShards []byte, availableIndices []int) []byte {
	if len(availableShards) < enc.dataShards || len(availableShards) != len(availableIndices) {
		panic("Not enough shards to reconstruct data")
	}

	// Only need the number of data shards to reconstruct
	shards := availableShards[:enc.dataShards]
	indices := availableIndices[:enc.dataShards]

	// Create original data array
	originalData := make([]byte, enc.dataShards)

	// For each data position
	for i := 0; i < enc.dataShards; i++ {
		// Calculate interpolation result
		result := byte(0)

		// Use Lagrange interpolation to recover original data
		for j := 0; j < enc.dataShards; j++ {
			term := shards[j]

			// Calculate the Lagrange basis function
			for k := 0; k < enc.dataShards; k++ {
				if j != k {
					// Calculate (x_i - x_k)
					numerator := enc.field.Sub(enc.evalPoints[i], enc.evalPoints[indices[k]])
					// Calculate (x_j - x_k)
					denominator := enc.field.Sub(enc.evalPoints[indices[j]], enc.evalPoints[indices[k]])
					// Division
					factor := enc.field.Div(numerator, denominator)
					// Multiply by the current term
					term = enc.field.Mul(term, factor)
				}
			}

			// Add this term to the result
			result = enc.field.Add(result, term)
		}

		originalData[i] = result
	}

	return originalData
}
