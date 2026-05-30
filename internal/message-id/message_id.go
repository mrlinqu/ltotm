package message_id

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"strconv"
	"time"
)

const (
	messageIdLen = 18
)

var (
	ErrIncorrectMsgId = errors.New("incorrect message id")
)

func rnd() []byte {
	ret := make([]byte, 0, messageIdLen)
	for i := 8; i > 0; {
		a, err := rand.Int(rand.Reader, big.NewInt(255))
		if err != nil {
			continue
		}

		ret = append(ret, a.Bytes()...)

		i--
	}

	return ret
}

func GenerateFromNow(dur int64) string {
	ttl := time.Now().Add(time.Duration(dur) * time.Hour)

	return Generate(ttl)
}

func Generate(ttl time.Time) string {
	salt := rnd()
	tm := strconv.FormatInt(ttl.Unix(), 10)

	ret := append(salt, []byte(tm)...)

	return base64.URLEncoding.EncodeToString(ret)
}

func GetTtl(msgId string) (time.Time, error) {
	buf, _ := base64.URLEncoding.DecodeString(msgId)
	if len(buf) != messageIdLen {
		return time.Time{}, ErrIncorrectMsgId
	}

	tm, err := strconv.Atoi(string(buf[8:]))
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(tm), 0), nil
}
