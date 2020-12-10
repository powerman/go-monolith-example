package dom_test

import (
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/internal/dom"
)

func TestNewID(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	id1 := dom.NewID()
	id2 := dom.NewID()
	t.Match(id1, `^[a-z0-9]{26}$`)
	t.Match(id2, `^[a-z0-9]{26}$`)
	t.NotEqual(id1, id2)
}
