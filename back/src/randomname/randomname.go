package randomname

import (
	"math/rand"
	"strconv"
)

var syllables = []string{"mon", "fay", "shi", "zag", "blarg", "rash", "izen"}

func GenerateNickname() string {
	res := ""
	for i := 0; i < 3; i++ {
		res += syllables[rand.Intn(len(syllables))]
	}
	res += strconv.Itoa(rand.Intn(10))
	res += strconv.Itoa(rand.Intn(10))
	return res
}
