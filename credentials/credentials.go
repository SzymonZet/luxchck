package credentials

import (
	"SzymonZet/LuxmedCheck/erroring"
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func GetSecureString(name string) string {
	fmt.Printf("Type %v and press enter: \t", name)
	val, err := term.ReadPassword(int(syscall.Stdin))
	erroring.QuitIfError(err, fmt.Sprintf("error when trying to obtain %v", name))
	fmt.Println()

	return string(val)
}
