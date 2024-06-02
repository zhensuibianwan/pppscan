package fingerprintManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"sort"
)

func FingerprintsList(key string) []publicCode.Fingerprint {

	fingerprints := mysqldb.SearchFingerprintByName(key)
	sort.Sort(publicCode.FingerprintSlice(fingerprints))
	return fingerprints
}
