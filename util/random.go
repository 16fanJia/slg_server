package util

import "math/rand"

/*
=======用于随机生成一串随机数字======= 作为链接的 secretKey
*/

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	//指定n位
	res := make([]rune, n)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}
