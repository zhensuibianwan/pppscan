package publicCode

type Fingerprint struct {
	UUID            string                `json:"uuid"`
	Name            string                `json:"name"`
	FingerprintScan []FingerprintScanData `json:"fingerprintScan"`
	PocsInfo        []PocsInfoData        `json:"pocsInfo"`
}

type FingerprintScanData struct {
	Path           string            `json:"path"`
	RequestMethod  string            `json:"request_method"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestData    string            `json:"request_data"`
	StatusCode     int               `json:"status_code"`
	Headers        map[string]string `json:"headers"`
	Keyword        []string          `json:"keyword"`
	FaviconHash    []string          `json:"favicon_hash"`
	Priority       int               `json:"priority"`
	Name           string            `json:"name"`
}

type PocsInfoData struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type FingerprintSlice []Fingerprint

func (s FingerprintSlice) Len() int           { return len(s) }
func (s FingerprintSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s FingerprintSlice) Less(i, j int) bool { return s[i].Name < s[j].Name }
