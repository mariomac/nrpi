package iot

import "io"

type Connector func() io.Reader

