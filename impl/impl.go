package impl

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
