package tvdb

import (
	"strconv"
	"strings"
	"time"
)

// pipeList type representing pipe-separated string values.
type pipeList []string

// unixTime type representing UNIX time in seconds from 1970-01-01.
type unixTime time.Time

// UnmarshalJSON unmarshals an XML element with string value into a pip-separated list of strings.
func (p *pipeList) UnmarshalJSON(data []byte) (err error) {
	*p = strings.Split(strings.Trim(string(data), "|"), "|")

	return
}

// UnmarshalJSON unmarshals a unix type byte slice to a unixTime object.
func (t *unixTime) UnmarshalJSON(data []byte) (err error) {
	unixSeconds, err := strconv.ParseInt(string(data), 10, 64)

	if err != nil {
		return
	}

	*t = unixTime(time.Unix(unixSeconds, 0))

	return
}
