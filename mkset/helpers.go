package main

import "strings"

func receiver(v string) string {
	return strings.ToLower(v[0:1])
}

func isPrivate(v string) bool {
	a := v[0:1]
	b := strings.ToLower(a)
	return a == b
}

func makePublic(v string) string {
	a := strings.ToUpper(v[0:1])
	return a + v[1:]
}
