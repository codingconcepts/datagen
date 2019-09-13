package runner

import (
	"fmt"
	"strings"
)

type errBuilder struct {
	b   strings.Builder
	err error
}

func (ew *errBuilder) write(i interface{}) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.b.WriteString(fmt.Sprintf("%v", i))
}
