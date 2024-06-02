package pocManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"sort"
)

func PocObjects(key string) []publicCode.Poc {
	pocs := mysqldb.SearchPocByName(key)
	sort.Sort(publicCode.PocSlice(pocs))
	return pocs
}
