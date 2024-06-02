package fingerprintManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"fmt"
	"os"
)

func RemoveFingerprint(fingerprints []publicCode.Fingerprint) (int, string) {

	var msg string
	var status int

	//数据库中删除poc
	status, msg = mysqldb.DelFingerprintDB(fingerprints)
	if status != 1 {
		return status, msg
	}
	for _, fingerprint := range fingerprints {
		filePath := fmt.Sprintf("./fingerprint/%s.json", fingerprint.Name)
		err := os.Remove(filePath)
		if err != nil {
			return -1, "本地文件删除失败！"
		}
	}

	return 1, "成功删除！"
}

func RemoveFingerprintFile(fingerprints []publicCode.Fingerprint) (int, string) {

	for _, fingerprint := range fingerprints {
		filePath := fmt.Sprintf("./fingerprint/%s.json", fingerprint.Name)
		err := os.Remove(filePath)
		if err != nil {
			return -1, "本地文件删除失败！"
		}
	}

	return 1, "成功删除！"

}
