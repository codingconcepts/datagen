package runner

import (
	"io/ioutil"
)

func (r *Runner) mustDumpQuery(stmt []byte) {
	if err := ioutil.WriteFile(r.queryErrFile, stmt, 0644); err != nil {
		panic(err)
	}
}
