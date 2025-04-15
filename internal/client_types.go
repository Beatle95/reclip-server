package internal

import (
	"github.com/gammazero/deque"
)

type ClipboardData struct {
	text deque.Deque[string]
}

type ClientData struct {
	id   string
	name string
	data ClipboardData
}
