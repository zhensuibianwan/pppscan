package publicCode

type Poc struct {
	UUID        string        `json:"uuid" yaml:"uuid"`
	Name        string        `json:"name" yaml:"name"`
	Hunter      string        `json:"hunter" yaml:"hunter"`
	Fofa        string        `json:"fofa" yaml:"fofa"`
	CMS         string        `json:"cms" yaml:"cms"`
	Description string        `json:"description" yaml:"description"`
	OptionValue int           `json:"optionValue" yaml:"optionValue"`
	NeedData    []Need        `json:"needData" yaml:"needData"`
	Request     []RequestData `json:"request" yaml:"request"`
}

type PocSlice []Poc

func (s PocSlice) Len() int           { return len(s) }
func (s PocSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PocSlice) Less(i, j int) bool { return s[i].Name < s[j].Name }

type RequestData struct {
	PocString string `json:"pocString" yaml:"pocString"`
	Status    string `json:"status" yaml:"status"`
	Check     string `json:"check" yaml:"check"`
	Print     string `json:"print" yaml:"print"`
}

type Need struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Rule struct {
	Method  string            `json:"method" yaml:"method"`
	Path    string            `json:"path" yaml:"path"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Body    string            `json:"body" yaml:"body"`
	Status  string            `json:"status" yaml:"status"`
	Check   string            `json:"check" yaml:"check"`
	Print   string            `json:"print" yaml:"print"`
}

type RespondData struct {
	Status  string
	Headers map[string]string
	Body    string
	IsCheck bool
}
