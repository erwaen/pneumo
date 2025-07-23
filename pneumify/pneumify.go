package pneumify


import (
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	"os"
	"math/big"
	"bytes"
)

func sha256Pneumo(input string) []byte {
	algo := sha256.New()
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err:= encoding.Encode(bytes)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

func PneumifyURL(url string) string {
	hash := sha256Pneumo(url)
	generatedNumber:= new(big.Int).SetBytes(hash).Uint64()

	buf := make([]byte, 0, 20) 
    buf = fmt.Appendf(buf, "%d", generatedNumber)

    finalString := base58Encoded(buf)
	return finalString[:8] 
}
