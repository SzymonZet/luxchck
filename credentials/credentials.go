package credentials

import (
	"fmt"
	"syscall"
	"szymonzet/luxchck/erroring"

	"golang.org/x/term"
)

func GetSecureString(name string) string {
	fmt.Printf("Type %v and press enter: \t", name)
	val, err := term.ReadPassword(int(syscall.Stdin))
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to obtain %v", name))
	fmt.Println()

	return string(val)
}
