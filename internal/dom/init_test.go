package dom_test

import (
	"testing"

	"github.com/powerman/check"
	_ "github.com/smartystreets/goconvey/convey"

	"github.com/powerman/go-monolith-example/pkg/def"
)

func TestMain(m *testing.M) {
	def.Init()
	check.TestMain(m)
}
