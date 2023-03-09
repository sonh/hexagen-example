package field

import (
	"math/rand"
	"time"
)

type String struct {
	value string
}

func (field String) String() string {
	return field.value
}

func NewNullString() String {
	return String{}
}

func NewString(value string) String {
	return String{
		value: value,
	}
}

func NewRandomString(n int) String {
	return String{
		value: RandStringRunes(n),
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
