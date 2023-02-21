package errorhandler

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
