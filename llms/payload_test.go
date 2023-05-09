package llms

import (
	"testing"
)

func TestNewPayload(t *testing.T) {
	a := Payload{Name: `a`}
	b := a

	b.Name = `b`

	println(` a.name: `, a.Name, ` b.name: `, b.Name)
}
