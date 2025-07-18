package hash

import (
	"fmt"
	"os"
)

func Check(should string) error {
	is, err := Calc(os.Args[0])
	if err != nil {
		return fmt.Errorf("cannot calculate own hash: %w", err)
	}
	if is != should {
		// fmt.Printf("%s (is)\n%s (should)\n", is, should)
		return fmt.Errorf("local: %s remote: %s", should, is)
	}
	return nil
}
