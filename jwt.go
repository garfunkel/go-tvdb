package tvdb

import (
	"strings"
	"errors"
	"encoding/base64"
	"encoding/json"
)

type jwt struct {
	Header struct {
		Algorithm string `json:"alg"`
	}
	Claims struct {
		IssuedAt unixTime `json:"orig_iat"`
		Expires unixTime `json:"exp"`
		ID string `json:"id"`
	}
	Signature string
	JWT string
}

func DecodeJWT(jwtStr string) (j jwt, err error) {
	fields := strings.Split(jwtStr, ".")

	if len(fields) != 3 {
		err = errors.New("Invalid JWT string")

		return
	}

	header, err := base64.StdEncoding.DecodeString(fields[0])

	if err != nil {
		return
	}

	claims, err := base64.StdEncoding.DecodeString(fields[1])

	if err != nil {
		return
	}

	signature, err := base64.RawURLEncoding.DecodeString(fields[2])

	if err != nil {
		return
	}

	j = jwt{
		JWT: jwtStr,
		Signature: string(signature),
	}

	err = json.Unmarshal(header, &j.Header)

	if err != nil {
		return
	}

	err = json.Unmarshal(claims, &j.Claims)

	return
}
