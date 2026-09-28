package shortener

import (
	"github.com/pedroalbanese/ff1"
)

func GenerateShort(cnt uint64) string {
	encoded := Encode(cnt)

	key := []byte("988303e0afa5822f9f95fb76681eb818")
	tweak := []byte{}

	cipher, err := ff1.NewCipher(62, 8, key, tweak)
	if err != nil {
		panic(err)
	}

	encrypted, err := cipher.Encrypt(encoded)
	if err != nil {
		panic(err)
	}

	return encrypted
}
