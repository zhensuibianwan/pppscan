package pocManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"fmt"
	"os"
)

func RemovePoc(pocs []publicCode.Poc) (int, string) {

	var msg string
	var status int

	//数据库中删除poc
	status, msg = mysqldb.DelPocDB(pocs)
	if status != 1 {
		return status, msg
	}
	for _, poc := range pocs {
		filePath := fmt.Sprintf("./poc/%s.yaml", poc.Name)
		err := os.Remove(filePath)
		if err != nil {
			return -1, "本地文件删除失败！"
		}
	}

	return 1, "成功删除！"
}

func RemovePocFile(pocs []publicCode.Poc) (int, string) {

	for _, poc := range pocs {
		filePath := fmt.Sprintf("./poc/%s.yaml", poc.Name)
		err := os.Remove(filePath)
		if err != nil {
			return -1, "本地文件删除失败！"
		}
	}

	return 1, "成功删除！"

}
