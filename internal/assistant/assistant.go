package assistant

const markdownScriptParts = 3

type Assistant interface {
	Query(string) string
}
