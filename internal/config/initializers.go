package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

const (
	privKeyPath = "./internal/keys/csrf.pem"
	pubKeyPath  = "./internal/keys/csrf-pub.pem"
)

var (
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func InitJWT() error {
	signBytes, err := os.ReadFile(privKeyPath)
	if err != nil {
		return err
	}

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

func CreateKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
		return err
	}

	privateFile, err := os.Create(privKeyPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer privateFile.Close()

	// Save the private key to file
	privatePEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateFile, privatePEM); err != nil {
		log.Fatal(err)
		return err
	}

	publicKey := &privateKey.PublicKey
	publicFile, err := os.Create(pubKeyPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer publicFile.Close()

	// Save the public key to file
	publicPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := pem.Encode(publicFile, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicPEM,
	}); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
