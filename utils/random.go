package utils

import "math/rand"

func RandomString(n int) string {
	var letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomEmail() string {
	return RandomString(6) + "@example.com"
}
