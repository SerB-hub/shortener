package hasher

import (
	"crypto/md5"
	"fmt"
	"github.com/deatil/go-encoding/encoding"
)

type Hasher interface {
	Hash(str string) (string, error)
}

type MD5Base62Hasher struct{}

func (h *MD5Base62Hasher) Hash(str string) (string, error) {
	hash := md5.Sum([]byte(str))
	shortHash := hash[0:6]
	//decimalHash, err := strconv.ParseInt(fmt.Sprintf("%x", shortHash), 16, 64)

	//if err != nil {
	//	panic(err)
	//}

	hashBase62 := encoding.FromBytes(
		shortHash, // []byte(strconv.FormatInt(decimalHash, 10)),
	).
		Base62Encode().
		ToString()

	fmt.Printf(
		"created new hash %v from string %v\n",
		hashBase62,
		str,
	)

	return hashBase62, nil

	//decimalBigIntHash := new(big.Int)
	//decimalBigIntHash.SetString(string(shortHash), 16)
	//decimalHash := int(decimalBigIntHash.Int64())
	//hashBase62 := encoding.FromString(strconv.Itoa())
}
