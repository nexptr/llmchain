package prompts

import (
	"testing"
)

func TestTemplate_Render(t *testing.T) {

	msg, err := QAPrompt.Render(H{`context`: `hello`, "question": "world"})

	if err != nil {
		t.Fatal(err.Error())
	}

	println(msg)

}
