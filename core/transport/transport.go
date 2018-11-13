package transport

import "io"

// Receiver is any function that returns a io.Reader, which will allow reading metrics that are
// received from remote sources.
type Receiver func() io.Reader
