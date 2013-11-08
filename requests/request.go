package requests

import (
	"net/http"
)

type Request struct {
	RawRequest *http.Request
	Params     map[string]string
}
