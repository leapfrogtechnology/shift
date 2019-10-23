package system

import "os/user"

// CurrentUser gives the username of the current user.
func CurrentUser() string {
	user, _ := user.Current()

	return user.Username
}
