// +build integration

package dal

import (
	"testing"
	"time"

	"github.com/powerman/check"
	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

func TestExampleSmoke(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	var (
		now  = time.Now().Truncate(time.Second)
		idU1 = dom.UserID(1)
		idU2 = dom.UserID(2)
	)

	res, err := r.Example(ctx, idU1)
	t.Err(err, app.ErrNotFound)
	t.Nil(res)

	err = r.IncExample(ctx, idU1)
	t.Nil(err)
	err = r.IncExample(ctx, idU1)
	t.Nil(err)

	res, err = r.Example(ctx, idU1)
	t.Nil(err)
	t.Equal(res.Counter, 2)
	t.GE(res.Mtime, now)
	res, err = r.Example(ctx, idU2)
	t.Err(err, app.ErrNotFound)
	t.Nil(res)
}
