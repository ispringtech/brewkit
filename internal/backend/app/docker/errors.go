package docker

import (
	"fmt"
)

type RequestError struct {
	Output string
	Code   int
}

func (e RequestError) Error() string {
	msg := fmt.Sprintf("request to docker client failed: code %d\n", e.Code)
	if e.Output != "" {
		msg += fmt.Sprintf("%s\n", e.Output)
	}

	return msg
}
