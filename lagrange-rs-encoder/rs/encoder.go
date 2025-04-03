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

	// 使用连续整数作为评估点 (从1开始)
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

	// 使用拉格朗日插值计算冗余数据
	enc.lagrangeInterpolation(message, encoded)

	return encoded
}

// lagrangeInterpolation calculates redundant shards using Lagrange interpolation
func (enc *RSEncoder2) lagrangeInterpolation(message []byte, encoded []byte) {
	// 对每个冗余分片位置
	for i := enc.dataShards; i < enc.totalShards; i++ {
		// 计算多项式在此点的值
		result := byte(0)

		// 构建拉格朗日插值多项式
		for j := 0; j < enc.dataShards; j++ {
			term := message[j]

			// 计算拉格朗日基函数
			for k := 0; k < enc.dataShards; k++ {
				if j != k {
					// 计算 (x - x_k)
					numerator := enc.field.Sub(enc.evalPoints[i], enc.evalPoints[k])
					// 计算 (x_j - x_k)
					denominator := enc.field.Sub(enc.evalPoints[j], enc.evalPoints[k])
					// 除法
					factor := enc.field.Div(numerator, denominator)
					// 乘以当前项
					term = enc.field.Mul(term, factor)
				}
			}

			// 将此项添加到结果中
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

	// 对每个冗余分片位置
	for i := enc.dataShards; i < enc.totalShards; i++ {
		result := byte(0)

		// 获取评估点
		x := enc.evalPoints[i]

		// 使用霍纳方法计算多项式值
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

	// 只需要数据分片数量的分片来重建
	shards := availableShards[:enc.dataShards]
	indices := availableIndices[:enc.dataShards]

	// 创建原始数据数组
	originalData := make([]byte, enc.dataShards)

	// 对每个数据位置
	for i := 0; i < enc.dataShards; i++ {
		// 计算插值结果
		result := byte(0)

		// 使用拉格朗日插值恢复原始数据
		for j := 0; j < enc.dataShards; j++ {
			term := shards[j]

			// 计算拉格朗日基函数
			for k := 0; k < enc.dataShards; k++ {
				if j != k {
					// 计算 (x_i - x_k)
					numerator := enc.field.Sub(enc.evalPoints[i], enc.evalPoints[indices[k]])
					// 计算 (x_j - x_k)
					denominator := enc.field.Sub(enc.evalPoints[indices[j]], enc.evalPoints[indices[k]])
					// 除法
					factor := enc.field.Div(numerator, denominator)
					// 乘以当前项
					term = enc.field.Mul(term, factor)
				}
			}

			// 将此项添加到结果中
			result = enc.field.Add(result, term)
		}

		originalData[i] = result
	}

	return originalData
}
