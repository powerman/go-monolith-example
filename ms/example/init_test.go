package example

import (
	"testing"

	"github.com/powerman/go-monolith-example/internal/def"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
	"github.com/powerman/go-monolith-example/proto"
	"github.com/powerman/gotest/testinit"
)

func TestMain(m *testing.M) { testinit.Main(m) }

const (
	serialMain = iota
	serialIntegration
)

func init() { testinit.Setup(serialMain, setupMain) }

func setupMain() {
	def.Init()
}

var (
	ctx         = def.NewContext(app.ServiceName)
	tokenAdmin  = proto.AccessToken("admin")
	tokenUser   = proto.AccessToken("user")
	userIDAdmin = dom.UserID(1)
	userIDUser  = dom.UserID(2)
)
