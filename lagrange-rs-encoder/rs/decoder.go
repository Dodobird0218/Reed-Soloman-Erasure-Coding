package rs

import (
	"fmt"
	"rs-encoder/gf"
)

// RSDecoder Reed-Solomon decoder
type RSDecoder struct {
	field       *gf.GF // GF(2^8) finite field
	dataShards  int    // Number of original data shards
	totalShards int    // Total number of shards
	evalPoints  []byte // Evaluation points, same as used by the encoder
}

// NewRSDecoder Create a new Reed-Solomon decoder
func NewRSDecoder(field *gf.GF, dataShards, totalShards int) *RSDecoder {
	decoder := &RSDecoder{
		field:       field,
		dataShards:  dataShards,
		totalShards: totalShards,
	}
	decoder.generateEvalPoints()
	return decoder
}

// generateEvalPoints Generate evaluation points, same as the encoder
func (dec *RSDecoder) generateEvalPoints() {
	dec.evalPoints = make([]byte, dec.totalShards)

	// Use consecutive integers as evaluation points (starting from 1), same as the encoder
	for i := 0; i < dec.totalShards; i++ {
		dec.evalPoints[i] = byte(i + 1)
	}

	// Output evaluation points information
	fmt.Println("Reed-Solomon decoder evaluation points:")
	for i, point := range dec.evalPoints {
		fmt.Printf("  Point[%d] = 0x%02x\n", i, point)
	}
}

// Decode Recover the original message from any dataShards shards
// availableShards: Available shard data
// availableIndices: Corresponding shard indices (0-based)
func (dec *RSDecoder) Decode(availableShards []byte, availableIndices []int) []byte {
	if len(availableShards) < dec.dataShards || len(availableShards) != len(availableIndices) {
		panic("Not enough shards to reconstruct data")
	}

	// Only dataShards shards are needed to recover the original data
	shards := availableShards[:dec.dataShards]
	indices := availableIndices[:dec.dataShards]

	// Create an array for the recovered original data
	decodedData := make([]byte, dec.dataShards)

	// Recover each original data position
	for i := 0; i < dec.dataShards; i++ {
		// Use Lagrange interpolation to calculate the value at the i-th original data position
		result := byte(0)

		// Construct polynomial interpolation
		for j := 0; j < dec.dataShards; j++ {
			// Get the value of the known shard
			y_j := shards[j]

			// Skip the case where the value is 0 (optimization)
			if y_j == 0 {
				continue
			}

			// Calculate the Lagrange basis function
			basis := byte(1)

			// Calculate L_j(x) for the Lagrange basis function
			for k := 0; k < dec.dataShards; k++ {
				if j != k {
					// Calculate (x - x_k)
					numerator := dec.field.Sub(dec.evalPoints[i], dec.evalPoints[indices[k]])
					// Calculate (x_j - x_k)
					denominator := dec.field.Sub(dec.evalPoints[indices[j]], dec.evalPoints[indices[k]])
					// Division
					factor := dec.field.Div(numerator, denominator)
					// Multiply by the current basis function value
					basis = dec.field.Mul(basis, factor)
				}
			}

			// Calculate the contribution of this term: y_j * L_j(x)
			term := dec.field.Mul(y_j, basis)

			// Add to the result
			result = dec.field.Add(result, term)
		}

		decodedData[i] = result
		fmt.Printf("Decoded data at position %d: 0x%02x\n", i, result)
	}

	return decodedData
}

// DecodeLastShards Recover the original message from the last dataShards shards of the encoded result
func (dec *RSDecoder) DecodeLastShards(encodedData []byte) []byte {
	if len(encodedData) < dec.totalShards {
		panic("Not enough shards in encoded data")
	}

	// Get the last dataShards shards
	lastShards := encodedData[dec.totalShards-dec.dataShards:]

	// Construct the index array
	indices := make([]int, dec.dataShards)
	for i := 0; i < dec.dataShards; i++ {
		indices[i] = dec.totalShards - dec.dataShards + i
	}

	// Call the main decode function
	return dec.Decode(lastShards, indices)
}
