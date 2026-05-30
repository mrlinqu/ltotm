package message_id

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	a := []byte{105, 233, 168, 251, 92, 102, 222, 139, 49, 55, 56, 48, 49, 51, 48, 50, 49, 50}
	b := base64.URLEncoding.EncodeToString(a)

	//buf := make([]byte, 0, messageIdLen)
	c, err := base64.URLEncoding.DecodeString(b)
	fmt.Println(err)
	fmt.Printf("%v\n%v\n", a, c)
}

func TestGetTtl(t *testing.T) {
	t.Parallel()

	//expected := time.Unix(1780213389, 0)
	//actual, err := GetTtl("rs7TVO9RcUsxNzgwMjEzMzg5")

	actual, err := GetTtl("5uK0QPsHvOcxNzgwMjQ5NTEw")
	expected := time.Unix(1780213389, 0)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestGenerateFromNow(t *testing.T) {
	hash := GenerateFromNow(1)

	expected := time.Now().Add(1 * time.Hour).Truncate(time.Second)

	actual, err := GetTtl(hash)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
