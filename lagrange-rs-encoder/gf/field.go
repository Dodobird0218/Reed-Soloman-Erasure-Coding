package gf

// GF represents the GF(2^8) finite field
type GF struct {
	// Exponent and logarithm tables for accelerating operations
	expTable [256]byte
	logTable [256]byte
	// Primitive polynomial
	primitivePoly byte
}

// NewGF creates and initializes a GF(2^8) finite field
func NewGF(primitivePoly byte) *GF {
	field := &GF{
		primitivePoly: primitivePoly,
	}
	field.generateTables()
	return field
}

// Add performs addition operation in GF(2^8) (XOR)
func (f *GF) Add(a, b byte) byte {
	return a ^ b
}

// Sub in GF(2^8), addition and subtraction are the same
func (f *GF) Sub(a, b byte) byte {
	return a ^ b
}

// Mul performs multiplication operation in GF(2^8)
func (f *GF) Mul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	// Fast multiplication using logarithm tables
	sum := int(f.logTable[a]) + int(f.logTable[b])
	if sum >= 255 {
		sum -= 255
	}
	return f.expTable[sum]
}

// Div performs division operation in GF(2^8)
func (f *GF) Div(a, b byte) byte {
	if a == 0 {
		return 0
	}
	if b == 0 {
		panic("Division by zero")
	}
	// Fast division using logarithm tables
	diff := int(f.logTable[a]) - int(f.logTable[b])
	if diff < 0 {
		diff += 255
	}
	return f.expTable[diff]
}

// Pow calculates power operation in GF(2^8)
func (f *GF) Pow(a byte, power int) byte {
	if a == 0 {
		return 0
	}
	if power == 0 {
		return 1
	}

	log := int(f.logTable[a])
	result := (log * power) % 255
	if result < 0 {
		result += 255
	}
	return f.expTable[result]
}

// Inv calculates the multiplicative inverse in GF(2^8)
func (f *GF) Inv(a byte) byte {
	if a == 0 {
		panic("0 has no multiplicative inverse")
	}
	// In GF(2^8), the inverse of a is a^254
	return f.expTable[255-f.logTable[a]]
}

// generateTables generates exponent and logarithm tables
func (f *GF) generateTables() {
	x := byte(1)
	for i := 0; i < 255; i++ {
		f.expTable[i] = x
		if i < 254 {
			// Calculate the next exponent value: x = x * 2
			// If the result would overflow, apply the primitive polynomial
			if x&0x80 != 0 {
				x = (x << 1) ^ f.primitivePoly
			} else {
				x = x << 1
			}
		}
	}
	f.expTable[255] = f.expTable[0] // Cyclic property

	// Generate logarithm table
	for i := 0; i < 256; i++ {
		if i == 0 {
			f.logTable[0] = 0 // log(0) is undefined, set to 0 to indicate a special case
		} else {
			for j := 0; j < 256; j++ {
				if f.expTable[j] == byte(i) {
					f.logTable[i] = byte(j)
					break
				}
			}
		}
	}
}
