package helper

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/libonomy/node-extract/constants"
)

//StringContains Funtion for Checking if slice of array (a) contains string (x)
func StringContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

//DecryptDatafromClient decryption of data coming from client
func DecryptDatafromClient(encryptedBytes []byte) []byte {

	decryptedBytes, err := constants.PrivateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	fmt.Println("decrypted message: ", string(decryptedBytes))
	return decryptedBytes
}

// ServerGenerateKey server keys
func ServerGenerateKey() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	/* 	modulusBytes := base64.StdEncoding.EncodeToString(privateKey.N.Bytes())
	   	privateExponentBytes := base64.StdEncoding.EncodeToString(privateKey.D.Bytes()) */
	/* fmt.Println("showing modules", modulusBytes)
	fmt.Println("showing private key", privateExponentBytes) */
	constants.PublicKey = privateKey.PublicKey
	constants.PrivateKey = privateKey
	fmt.Println("keys Registered")
}
