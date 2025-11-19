package uuidx

import (
	"strings"

	"github.com/gofrs/uuid"
)

func GenCert() string {
	u2, _ := uuid.NewV4()
	res := u2.String()
	res = strings.Replace(res, "-", "", -1)
	return res
}

func GenAtkId() string {
	u2, _ := uuid.NewV4()
	res := u2.String()
	res = strings.Replace(res, "-", "", -1)
	return res
}
