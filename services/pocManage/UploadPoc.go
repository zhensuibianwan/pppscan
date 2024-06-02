package pocManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"sync"
)

func UploadPocs(pocs []publicCode.Poc) (int, string) {

	var msg string
	var wg sync.WaitGroup
	for _, poc := range pocs {
		wg.Add(1)
		go func(poc publicCode.Poc) {
			defer wg.Done()
			p := mysqldb.CheckPocByUUID(poc.UUID)
			if p == 0 {
				errInt, errStr := mysqldb.SavePocDB(poc)
				if errInt == -1 {
					msg = errStr
				}
			}

		}(poc)

	}
	wg.Wait()
	if msg != "" {
		return -1, "部分poc导入数据库失败！"
	}

	return 1, "导入数据库成功！"

}
