package dal

import (
	"time"

	"github.com/powerman/go-monolith-example/internal/dom"
)

const (
	sqlExampleInc = `
INSERT INTO example (user_id, counter) VALUES (:user_id, 1)
ON DUPLICATE KEY
UPDATE counter = example.counter + 1, mtime = NOW()
	`
	sqlExampleGet = `
SELECT counter, mtime FROM example WHERE user_id = :user_id
	`
)

type (
	argExampleInc struct {
		UserID dom.UserID
	}

	argExampleGet struct {
		UserID dom.UserID
	}
	rowExampleGet struct {
		Counter int
		Mtime   time.Time
	}
)
