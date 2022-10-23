package tool

import (
	"math/rand"
	"time"
)

func RandaomString(n int) string {
	var letters = []byte("eKEzxcBvWCUfF9ilOopL51kXYG7gSrtyhNj4mVRnb2suJD3AdQvIvTcbZ8aM0PqH")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func RandomName() string {
	name := make([]byte, 20)
	i := 0
	for i < 20 {
		//利用rand.Intn(x)来伪随机生成一个[0,x)的数
		a := rand.Intn(26)
		b := rand.Intn(2)
		if b == 1 {
			name[i] = byte('A' + a)
		} else {
			name[i] = byte('a' + a)
		}
		i++
	}
	return string(name)
}
