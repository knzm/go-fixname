package rename

import (
	"fmt"
	"io/ioutil"
)

func plural(n int, singularUnit, pluralUnit string) string {
	var unit string
	if n == 1 {
		unit = singularUnit
	} else {
		unit = pluralUnit
	}
	return fmt.Sprintf("%d %s", n, unit)
}

func WriteFile(filename string, content []byte) error {
	return ioutil.WriteFile(filename, content, 0644)
}
