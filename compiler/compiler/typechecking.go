package zoe

const (
	CONTEXT_NAMESPACE = iota + 1
	CONTEXT_FUNCTION
	CONTEXT_IMPLEMENT
)

type Context struct {
}

func (c *Context) ReportError(_ string) {

}
