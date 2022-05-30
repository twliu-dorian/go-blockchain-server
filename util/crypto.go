package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"math/big"
)

func LeftPadding(data []byte, length int) (padded []byte) {
	padded = make([]byte, length)
	copy(padded[length-len(data):], data)
	return
}

func Sha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func SerializePublicKey(pub *ecdsa.PublicKey) []byte {
	return append(
		LeftPadding(pub.X.Bytes(), 32),
		LeftPadding(pub.Y.Bytes(), 32)...,
	)
}

func SerializePrivateKey(priv *ecdsa.PrivateKey) []byte {
	return LeftPadding(priv.D.Bytes(), 32)
}

func ParsePrivateKey(bytes []byte, curves ...elliptic.Curve) (privateKey *ecdsa.PrivateKey, err error) {
	if len(bytes) != 32 {
		err = errors.New("invalid private key length")
		return
	}

	curve := elliptic.P256()
	if len(curves) == 1 {
		curve = curves[0]
	}

	privateKey = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
		D: new(big.Int).SetBytes(bytes),
	}
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(bytes)

	return
}

func ParsePublicKey(bytes []byte, curves ...elliptic.Curve) (publicKey ecdsa.PublicKey, err error) {
	if len(bytes) != 32 {
		err = errors.New("invalid public key length")
		return
	}

	publicKey.Curve = elliptic.P256()
	if len(curves) == 1 {
		publicKey.Curve = curves[0]
	}
	publicKey.X = new(big.Int).SetBytes(bytes[:32])
	publicKey.Y = new(big.Int).SetBytes(bytes[32:])

	return
}