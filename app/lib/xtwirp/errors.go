package xtwirp

import (
	"github.com/twitchtv/twirp"
)

func FromError(err error) (twirp.Error, bool) {
	status, match := err.(twirp.Error)
	if match {
		return status, true
	}
	return twirp.InternalError(err.Error()), false
}
