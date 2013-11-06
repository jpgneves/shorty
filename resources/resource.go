package resources

import (
	"github.com/jpgneves/shorty/requests"
)

type Resource interface {
	Get(*requests.Request) *requests.Response
	Post(*requests.Request) *requests.Response
}
