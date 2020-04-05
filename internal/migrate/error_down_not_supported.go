package migrate

import "errors"

// ErrDownNotSupported must be returned from goose Down function in case
// this migration does not support downgrade.
var ErrDownNotSupported = errors.New("downgrade is not supported, restore from backup instead")
