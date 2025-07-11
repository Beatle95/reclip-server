package internal

import (
	"github.com/gammazero/deque"
)

type ClipboardData struct {
	Text deque.Deque[string]
}

type ClientData struct {
	Id   uint64
	Name string
	Data ClipboardData
}

func IsEqual(lhs ClientData, rhs ClientData) bool {
	if lhs.Id != rhs.Id || lhs.Name != rhs.Name {
		return false
	}

	if lhs.Data.Text.Len() != rhs.Data.Text.Len() {
		return false
	}

	for i := 0; i < lhs.Data.Text.Len(); i++ {
		if lhs.Data.Text.At(i) != rhs.Data.Text.At(i) {
			return false
		}
	}

	return true
}
