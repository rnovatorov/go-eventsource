package eventstore

import "errors"

var ErrConcurrentUpdate = errors.New("concurrent update")
