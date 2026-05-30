package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

func rnd() []byte {
	ret := make([]byte, 0, 16)
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

func gen() {
	rrr := rnd()
	fmt.Printf("%s\n", base64.URLEncoding.EncodeToString(rrr))

	tm := strconv.FormatInt(time.Now().Unix(), 10)

	rrr2 := append(rrr, []byte(tm)...)
	fmt.Printf("%s\n", base64.URLEncoding.EncodeToString(rrr2))

	//tm := fmt.Sprint(time.Now().Unix())
	fmt.Printf("%s|%s\n", tm, base64.URLEncoding.EncodeToString([]byte(tm)))
}

func main() {
	for i := 0; i < 10; i++ {
		gen()
		time.Sleep(1 * time.Second)
	}
}
