package dal

import (
	"time"
)

const (
	sqlExampleInc = `
 INSERT INTO example (user_id, counter)
 VALUES (:user_id, 1)
     ON DUPLICATE KEY
 UPDATE counter = example.counter + 1, mtime = NOW()
	`
	sqlExampleGet = `
 SELECT counter, mtime
   FROM example
  WHERE user_id = :user_id
	`
)

type (
	argExampleInc struct {
		UserID string
	}

	argExampleGet struct {
		UserID string
	}
	rowExampleGet struct {
		Counter int
		Mtime   time.Time
	}
)
