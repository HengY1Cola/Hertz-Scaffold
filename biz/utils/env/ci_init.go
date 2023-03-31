package env

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	if strings.HasSuffix(os.Args[0], ".test") {
		fmt.Println("########  init Ci Env")
		InitTestModules([]Module{UserMapJwtToken})
		fmt.Println("########  init Ci Over")
	}
}
