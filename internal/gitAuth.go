package internal

import (
	"errors"
	"regexp"
	"strings"

	"github.com/cli/go-gh/v2"
)

type Auth struct {
	username string
}

func GetAuthState() (login Auth, err error) {
	auth, _, err := gh.Exec("auth", "status")
	if err != nil {
		return
	}

	a := auth.String()
	if !strings.Contains(a, "Logged in") {
		err = errors.New("user is not currently logged into gh")
	}

	r, _ := regexp.Compile(`account.*?\(`)

	f := r.Find(auth.Bytes())
	var str string
	str = strings.Replace(string(f), "account", "", 1)
	str = strings.Replace(str, "(", "", 1)
	login = Auth{strings.Trim(str, " ")}

	return
}
