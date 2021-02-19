// +build integration

package dal_test

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
	r := newTestRepo(t)

	var (
		now    = time.Now().Truncate(time.Second)
		nameU1 = dom.NewUserName("1")
		nameU2 = dom.NewUserName("2")
	)

	res, err := r.Example(ctx, nameU1)
	t.Err(err, app.ErrNotFound)
	t.Nil(res)

	err = r.IncExample(ctx, nameU1)
	t.Nil(err)
	err = r.IncExample(ctx, nameU1)
	t.Nil(err)

	res, err = r.Example(ctx, nameU1)
	t.Nil(err)
	t.Equal(res.Counter, 2)
	t.GE(res.Mtime, now)
	res, err = r.Example(ctx, nameU2)
	t.Err(err, app.ErrNotFound)
	t.Nil(res)
}
