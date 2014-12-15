package utils

import "golang.org/x/crypto/openpgp/packet"

// GetAlgorithmName returns algorithm's name depending on its ID
func GetAlgorithmName(id packet.PublicKeyAlgorithm) string {
	switch id {
	case 1, 2, 3:
		return "RSA"
	case 16:
		return "ElGamal"
	case 17:
		return "DSA"
	case 18:
		return "ECDH"
	case 19:
		return "ECDSA"
	default:
		return "unknown"
	}
}
