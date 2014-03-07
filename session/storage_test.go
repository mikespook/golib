package session

import (
	"bytes"
	"testing"
)

var (
	legalKey    = []byte("1234567890123456")
	tooShortKey = []byte("1234567890")
	tooLongKey  = []byte("12345678901234567890")
	originText  = []byte("a brown fox jumps over the lazy dog")
)

func TestCrypt(t *testing.T) {
	// Vital, encrypt will reuse value's allocation.
	src := make([]byte, len(originText))
	copy(src, originText)
	cipherText, err := encrypt(legalKey, src)
	if err != nil {
		t.Error(err)
		return
	}
	text, err := decrypt(legalKey, cipherText)
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(originText, text) != 0 {
		t.Errorf("text[%s] != origin[%s]", text, originText)
		return
	}
	copy(src, originText)
	_, err = encrypt(tooShortKey, src)
	if err != errKeyTooShort {
		t.Errorf("Error %s needed", errKeyTooShort)
		return
	}
	_, err = decrypt(tooShortKey, cipherText)
	if err != errKeyTooShort {
		t.Errorf("Error %s needed", errKeyTooShort)
		return
	}

	copy(src, originText)
	cipherText, err = encrypt(tooLongKey, src)
	if err != nil {
		t.Error(err)
		return
	}
	text, err = decrypt(tooLongKey, cipherText)
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(originText, text) != 0 {
		t.Errorf("[%s] != [%s]", originText, text)
		return
	}
}

func TestCoding(t *testing.T) {
	srcM := make(M)
	srcM["foo"] = 123
	srcM["bar"] = "abc"
	var dstM M
	cipherText, err := encoding(legalKey, srcM)
	if err != nil {
		t.Error(err)
		return
	}
	if err := decoding(legalKey, cipherText, &dstM); err != nil {
		t.Error(err)
		return
	}

	if v, ok := dstM["foo"]; !ok || v != 123 {
		t.Errorf("Decofing issue: %s", dstM["foo"])
		return
	}

	cipherText, err = encoding(tooShortKey, srcM)
	if err != errKeyTooShort {
		t.Errorf("Error %s needed", errKeyTooShort)
		return
	}
}
