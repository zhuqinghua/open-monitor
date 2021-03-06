package models

import (
	"fmt"
	"math/big"
	"strings"
	"encoding/base64"
	"io/ioutil"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"log"
)

func DecryptRsa(inputString string) string {
	if !strings.HasPrefix(strings.ToLower(inputString), "rsa@") {
		return inputString
	}
	inputString = inputString[4:]
	result := inputString
	inputBytes,err := base64.StdEncoding.DecodeString(inputString)
	if err != nil {
		log.Printf("Input string format to base64 fail,%s \n", err.Error())
		return inputString
	}
	pemPath := "/data/certs/rsa_key"
	fileContent,err := ioutil.ReadFile(pemPath)
	if err != nil {
		log.Printf("Read file %s fail,%s \n", pemPath, err.Error())
		return result
	}
	block,_ := pem.Decode(fileContent)
	privateKeyInterface,err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("Parse private key fail,%s \n", err.Error())
		return result
	}
	privateKey := privateKeyInterface.(*rsa.PrivateKey)
	decodeBytes,err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, inputBytes)
	if err != nil {
		log.Printf("Decode fail,%s \n", err.Error())
		return result
	}
	result = string(decodeBytes)
	return result
}

func RSAEncryptByPrivate(orgidata []byte,privatekey string) ([]byte, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(privatekey)
	if err != nil {
		return nil, fmt.Errorf("RSASign private key is bad")
	}

	privInterface, err := x509.ParsePKCS8PrivateKey(decodeBytes)
	if err != nil {
		return nil, err
	}

	priv := privInterface.(*rsa.PrivateKey)

	k := (priv.N.BitLen() + 7) / 8
	tLen := len(orgidata)
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], orgidata)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(priv.N) > 0 {
		return nil, nil
	}
	var m *big.Int
	var ir *big.Int
	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}
	return m.Bytes(), nil
}