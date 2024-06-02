package fingerprintManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"sync"
)

func UploadFingerprints(fingerprints []publicCode.Fingerprint) (int, string) {

	var msg string
	var wg sync.WaitGroup
	worker := make(chan int, 100)
	for _, fingerprint := range fingerprints {
		wg.Add(1)
		go func(fingerprint publicCode.Fingerprint) {
			defer wg.Done()
			worker <- 1
			p := mysqldb.CheckFingerprintByUUID(fingerprint.UUID)
			<-worker
			if p == 0 {
				worker <- 1
				errInt, errStr := mysqldb.SaveFingerprintDB(fingerprint)
				if errInt == -1 {
					msg = errStr
				}
				<-worker
			}

		}(fingerprint)

	}
	wg.Wait()
	if msg != "" {
		return -1, "部分指纹导入数据库失败！"
	}

	return 1, "导入数据库成功！"

}

func UploadFingerprintsList(key string) []publicCode.Fingerprint {

	fileList := publicCode.FingerprintFileList(key)
	dbList := mysqldb.SearchFingerprintByName(key)
	for _, db := range dbList {
		for i, file := range fileList {
			if db.UUID == file.UUID {

				fileList = append(fileList[:i], fileList[i+1:]...)
			}
		}
	}
	return fileList
}
