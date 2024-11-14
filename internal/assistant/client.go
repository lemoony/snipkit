package assistant

type Client interface {
	Query(string) (string, error)
}
