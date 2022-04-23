package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

var CORRECT = []byte("hello signing world")
var WRONG = []byte("hello world")

func main() {
	TestECDSA()
	TestRSA()

	//openpgp.

	//Generate Keypair
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	//Sign a message
	sig := ed25519.Sign(priv, CORRECT)

	pubkey := priv.Public().(ed25519.PublicKey) //You can always create the PubKey from the PrivKey!
	shouldMatch := ed25519.Verify(pubkey, CORRECT, sig)
	shouldFail := ed25519.Verify(pubkey, WRONG, sig)

	if !shouldMatch || shouldFail {
		panic("Signature Error!")
	}

	fmt.Print("\nExample ED25519\n")
	fmt.Printf("Priv: %s\n", prettyPrint(priv))
	fmt.Printf("Pub: %s\n", prettyPrint(pubkey))
	fmt.Printf("Sig: %s\n", prettyPrint(sig))
}

func TestRSA() {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	correctSUM := sha512.Sum512(CORRECT)

	sig, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA512, correctSUM[:], nil)
	if err != nil {
		panic(err)
	}

	fmt.Print("\nExample RSA\n")
	fmt.Printf("Sig: %s \n", prettyPrint(sig))

}

func TestECDSA() {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	correctHash := sha256.Sum256(CORRECT) //ecdsa signs hashes!
	wrongHash := sha256.Sum256(WRONG)     //ecdsa signs hashes!

	sigASN1, err := ecdsa.SignASN1(rand.Reader, privKey, correctHash[:])

	pubKey := &privKey.PublicKey
	match := ecdsa.VerifyASN1(pubKey, correctHash[:], sigASN1)
	mustfail := ecdsa.VerifyASN1(pubKey, wrongHash[:], sigASN1)

	if !match || mustfail {
		panic("Signature Error!")
	}
	//TODO write a Test that shows a fail with ECDAS truncated hash!!
	fmt.Print("\nExample ECDSA\n")

	pkASN1, _ := x509.MarshalPKCS8PrivateKey(privKey) //PKCS #8, ASN.1 DER
	pubASN1, _ := x509.MarshalPKIXPublicKey(pubKey)   //PKIX, ASN.1 DER form

	fmt.Printf("Priv: %s\n", prettyPrint(pkASN1))
	fmt.Printf("Pub: %s\n", prettyPrint(pubASN1))
	fmt.Printf("Sig: %s\n", prettyPrint(sigASN1))
}

func prettyPrint(in []byte) string {
	sigB64 := base64.StdEncoding.EncodeToString(in)
	return fmt.Sprintf("%s", sigB64)
}
