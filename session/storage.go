package session

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"net/http"
)

func init() {
	gob.Register(M{})
}

type Storage interface {
	Clean(*Session) error
	Flush(*Session) error
	LoadTo(*http.Request, *Session) error
}

const (
	keySize    = 16
	aesKeySize = 32
)

var (
	defaultKey = genKey(keySize)
)

func SetKey(key []byte) {
	defaultKey = key[:keySize]
}

func GetKey() []byte {
	return defaultKey
}

var errKeyTooShort = errors.New("The key is too short")
var errValueTooShort = errors.New("The block is too short")

func encrypt(key, value []byte) ([]byte, error) {
	if len(key) < aesKeySize-keySize {
		return nil, errKeyTooShort
	}
	key = append(key[:aesKeySize-keySize], defaultKey...)
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

func decrypt(key, value []byte) ([]byte, error) {
	if len(key) < aesKeySize-keySize {
		return nil, errKeyTooShort
	}
	key = append(key[:aesKeySize-keySize], defaultKey...)
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
	return nil, errValueTooShort
}

func decoding(key []byte, src string, dst *M) error {
	// 1. base64 decoding
	buf, err := base64.StdEncoding.DecodeString(src)
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

func encoding(key []byte, src map[string]interface{}) (string, error) {
	// 1. gob encoding
	var buf bytes.Buffer
	g := gob.NewEncoder(&buf)
	if err := g.Encode(src); err != nil {
		return "", err
	}
	// 2. cypto encoding
	ciphertext, err := encrypt(key, buf.Bytes())
	if err != nil {
		return "", err
	}
	// 3. base64 encoding
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
