package pile

import "errors"

var (
	errDuplicate    = errors.New("pile number already registered")
	errUnknown      = errors.New("unknown pile")
	errNoPile       = errors.New("no pile available")
	errUnauthorized = errors.New("bus not authorized")
	errNumberUsed   = errors.New("number already in use")
)
