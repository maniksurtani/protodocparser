package impl

type CommentType int

const (
	ServiceComment CommentType = iota
	RpcComment
	OtherComment
)

type Example struct {
	Language string `json:"language,omitempty"`
	Code     string `json:"code,omitempty"`
}

type Rpc struct {
	Name     string     `json:"name,omitempty"`
	Request  string     `json:"request,omitempty"`
	Response string     `json:"response,omitempty"`
	Options  []string   `json:"options,omitempty"`
	Doc      string     `json:"doc,omitempty"`
	Examples []*Example `json:"examples,omitempty"`
}

type Service struct {
	Url      string     `json:"url,omitempty"`
	File     string     `json:"file,omitempty"`
	Sha      string     `json:"sha,omitempty"`
	Package  string     `json:"package,omitempty"`
	Name     string     `json:"name,omitempty"`
	Org      string     `json:"org,omitempty"`
	Design   string     `json:"design,omitempty"`
	Doc      string     `json:"doc,omitempty"`
	Api      bool       `json:"api,omitempty"`
	Examples []*Example `json:"examples,omitempty"`
	Rpcs     []*Rpc     `json:"rpcs,omitempty"`
}

type CommentBlock struct {
	Start int
	End   int
	Type  CommentType
}

// NewRpc creates a new instance of Rpc, initialized with sensible defaults.
func NewRpc() *Rpc {
	return &Rpc{
		Name:     "",
		Request:  "",
		Response: "",
		Options:  make([]string, 0),
		Doc:      "",
		Examples: make([]*Example, 0)}
}

// NewService creates a new instance of Service, initialized with sensible defaults.
func NewService() *Service {
	return &Service{
		Name:     "",
		Url:      "",
		File:     "",
		Package:  "",
		Org:      "",
		Design:   "",
		Doc:      "",
		Api:      false,
		Examples: make([]*Example, 0),
		Rpcs:     make([]*Rpc, 0)}
}
