package a

import "a/b" // want `.*go1\.100.*`

func f() {
	b.F()
}
