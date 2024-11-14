package ui

import "github.com/lemoony/snipkit/internal/ui/uimsg"

type MessagePrinter interface {
	Print(uimsg.Printable)
}
