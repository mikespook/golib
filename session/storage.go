package session

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
)

type Storage interface {
	Clean(*Session)
	Flush(*Session)
	LoadTo(*http.Request, *Session)
}

const (
	keySize = 16
)

var (
	defaultKey = genKey(keySize)
)

func SetKey(key []byte) {
	defaultKey = key[:keySize]
}

func encrypt(key, value []byte) ([]byte, error) {
	key = append(key[:32-keySize], defaultKey...)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, block.BlockSize())
	rand.Read(iv)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(value, value)
	return append(iv, value...), nil
}

var errTooShort = errors.New("The cipher text is too short.")

func decrypt(key, value []byte) ([]byte, error) {
	key = append(key[:32-keySize], defaultKey...)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(value) > block.BlockSize() {
		iv := value[:block.BlockSize()]
		value = value[block.BlockSize():]
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(value, value)
		return value, nil
	}
	return nil, errTooShort
}

func decoding(key, src []byte, dst *M) error {
	// 1. base64 decoding
	n := base64.StdEncoding.DecodedLen(len(src))
	buf := make([]byte, n)
	_, err := base64.StdEncoding.Decode(buf, src)
	if err != nil {
		return err
	}
	// 2. cypto decoding
	buf, err = decrypt(key, buf)
	if err != nil {
		return err
	}
	// 3. gob decoding
	g := gob.NewDecoder(bytes.NewBuffer(buf))
	if err = g.Decode(&dst); err != nil {
		return err
	}
	return nil
}

func encoding(key []byte, src map[string]interface{}) string {
	// 1. gob encoding
	var buf bytes.Buffer
	g := gob.NewEncoder(&buf)
	if err := g.Encode(src); err != nil {
		return ""
	}
	// 2. cypto encoding
	ciphertext, err := encrypt(key, buf.Bytes())
	if err != nil {
		return ""
	}
	// 3. base64 encoding
	return base64.StdEncoding.EncodeToString(ciphertext)
}
