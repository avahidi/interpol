package interpol

import (
	"testing"
)

func TestRegisterFactory(t *testing.T) {

	h1 := findDefaultHandlerFactory("dummy")
	if h1 != nil {
		t.Errorf("test internal error")
	}

	addDefaultHandlerFactory("dummy", newCopyHandler)

	h2 := findDefaultHandlerFactory("dummy")
	if h2 == nil {
		t.Errorf("Could not find registred handler")
	}
}
