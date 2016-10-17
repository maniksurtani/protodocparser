package impl

type CommentType int

const (
	ServiceComment CommentType = iota
	RpcComment
	OtherComment
)

type Example struct {
	Language string
	Code string
}

type Rpc struct {
	Name string
	Request string
	Response string
	Options []string
	Doc string
	Examples []*Example
}

type Service struct {
	Url string
	File string
	Package string
	Name string
	Org string
	Design string
	Doc string
	Examples []*Example
	Rpcs []*Rpc
}

type CommentBlock struct {
	Start int
	End int
	Type CommentType
}
