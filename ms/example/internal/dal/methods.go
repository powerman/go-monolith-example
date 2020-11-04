package dal

import (
	"database/sql"
	"errors"

	"github.com/powerman/sqlxx"

	"github.com/powerman/go-monolith-example/internal/dom"
	"github.com/powerman/go-monolith-example/ms/example/internal/app"
)

func (r *Repo) Example(ctx app.Ctx, userName dom.UserName) (res *app.Example, err error) {
	err = r.Tx(ctx, &sql.TxOptions{ReadOnly: true}, func(tx *sqlxx.Tx) error {
		var resExampleGet rowExampleGet
		err := tx.NamedGetContext(ctx, &resExampleGet, sqlExampleGet, argExampleGet{
			UserID: userName.ID(),
		})
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return app.ErrNotFound
		case err != nil:
			return err
		}
		res = &app.Example{
			Counter: resExampleGet.Counter,
			Mtime:   resExampleGet.Mtime,
		}
		return nil
	})
	return
}

func (r *Repo) IncExample(ctx app.Ctx, userName dom.UserName) error {
	return r.Tx(ctx, nil, func(tx *sqlxx.Tx) error {
		_, err := tx.NamedExecContext(ctx, sqlExampleInc, argExampleInc{
			UserID: userName.ID(),
		})
		return err
	})
}
