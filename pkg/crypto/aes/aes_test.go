package aes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESEncrypt(t *testing.T) {
	as := assert.New(t)

	data := []byte("hello world")

	bs64, err := EncryptToBase64(data, SecretKey)
	as.Nil(err)
	as.NotEmpty(bs64)

	t.Log(bs64)

	result, err := DecryptFromBase64(bs64, SecretKey)
	as.Nil(err)
	as.Equal(data, result)
}
