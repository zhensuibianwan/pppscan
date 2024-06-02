package publicCode

type AllSetting struct {
	DB    DBSetting    `json:"db" yaml:"db"`
	Proxy ProxySetting `json:"proxy" yaml:"proxy"`
	Scan  ScanSetting  `json:"scan" yaml:"scan"`
	Other OtherSetting `json:"other" yaml:"other"`
}

type DBSetting struct {
	Mode     bool   `json:"mode" yaml:"mode"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	UserName string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	//HostPublic     string `json:"hostPublic" yaml:"hostPublic"`
	//PortPublic     string `json:"portPublic" yaml:"portPublic"`
	//UsernamePublic string `json:"usernamePublic" yaml:"usernamePublic"`
	//PasswordPublic string `json:"passwordPublic" yaml:"passwordPublic"`
}

type ProxySetting struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Enable   bool   `json:"enable" yaml:"enable"`
	Mode     string `json:"mode" yaml:"mode"`
}

type ScanSetting struct {
	TimeOut   int64 `json:"timeout" yaml:"timeout"`
	ThreadNum int   `json:"threadNum" yaml:"threadNum"`
}

type OtherSetting struct {
	LocalLoading bool `json:"localLoading" yaml:"localLoading" `
	//OnlyLocalLoading bool `json:"onlyLocalLoading" yaml:"onlyLocalLoading" `
}

type Background struct {
	Image string `json:"image" yaml:"image"`
}

var ModeDB bool
var UserNameDB string //MySql的配置
var PassWordDB string
var HostDB string
var PortDB string

var UserNameProxy string //Proxy 的配置
var PassWordProxy string
var ModeProxy string
var HostProxy string
var PortProxy string
var EnableProxy bool

var TimeOut int64 //Scan的配置
var ThreadNum int

var LocalLoading bool
var OnlyLocalLoading bool

var BackgroundImage string
