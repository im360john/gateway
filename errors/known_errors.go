package errors

import "golang.org/x/xerrors"

var (
	ErrNotAuthorized = xerrors.New("not authorized")
)
