package assistant

const markdownScriptParts = 3

type Client interface {
	Query(string) (string, error)
}
