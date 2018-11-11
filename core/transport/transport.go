package transport

import "io"

type Receiver func() io.Reader

