package slice

import (
	"fmt"
	"os"
)

const logo = ` ______   __       __   ______   ______   
/\  ___\ /\ \     /\ \ /\  ___\ /\  ___\  
\ \___  \\ \ \____\ \ \\ \ \____\ \  __\  
 \/\_____\\ \_____\\ \_\\ \_____\\ \_____\
  \/_____/ \/_____/ \/_/ \/_____/ \/_____/`

var exitError = defaultExitError

func defaultExitError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
