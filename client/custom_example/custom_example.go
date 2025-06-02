package custom_example

// CaesarEncrypt performs Caesar cipher encryption on a byte array
// Parameters:
//   - data: the byte array to be encrypted
//   - shift: the number of positions to shift each byte (mod 256)
// Returns:
//   - encrypted byte array
func CaesarEncrypt(data []byte, shift int) []byte {
	encrypted := make([]byte, len(data))
	for i, b := range data {
		encrypted[i] = byte((int(b) + shift) % 256)
	}
	return encrypted
}
