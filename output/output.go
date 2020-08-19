package output

import "io"

type Output interface {
	Output(wr io.Writer) error
}
