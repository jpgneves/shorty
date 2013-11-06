package resources

type Resource interface {
	Get(url string) string
	Post(url string, data interface{}) string
}