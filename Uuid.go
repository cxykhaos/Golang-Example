package tool

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func Generate_uuid() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")

}
