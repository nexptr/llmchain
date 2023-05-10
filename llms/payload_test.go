package llms

import (
	"testing"
)

func TestNewPayload(t *testing.T) {
	a := ModelOptions{Name: `a`}
	b := a

	b.Name = `b`

	println(` a.name: `, a.Name, ` b.name: `, b.Name)
}
