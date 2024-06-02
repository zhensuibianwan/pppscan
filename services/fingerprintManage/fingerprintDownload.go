package fingerprintManage

import (
	"changeme/services/publicCode"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func DownLoadFingerprintS(fingerprints []publicCode.Fingerprint) (int, string) {

	var msg string
	var wg sync.WaitGroup
	var result int

	for _, fingerprint := range fingerprints {
		wg.Add(1)
		go func(fingerprint publicCode.Fingerprint) {
			defer wg.Done()
			//创建yaml文件，并将数据写入
			jsonData, err := json.Marshal(fingerprint)
			if err != nil {
				msg = fmt.Sprintf("Error marshaling YAML:%s", err)
				result = -1
				return
			}

			fileName := fingerprint.Name + ".json"
			// 获取当前工作目录
			currentDir, err := os.Getwd()
			if err != nil {
				msg = fmt.Sprintf("error getting current working directory: %s", err)
				result = -1
				return
			}

			// 构建目标目录路径
			targetDir := filepath.Join(currentDir, "/fingerprint")

			// 创建目标目录（如果不存在）
			err = os.MkdirAll(targetDir, 0755)
			if err != nil {
				msg = fmt.Sprintf("error creating target directory: %s", err)
				result = -1
				return
			}

			// 构建文件路径
			filePath := filepath.Join(targetDir, fileName)

			// 写入文件内容
			err = ioutil.WriteFile(filePath, jsonData, 0644)
			if err != nil {
				msg = fmt.Sprintf("error writing file: %s", err)
				result = -1
				return
			}
		}(fingerprint)

	}
	wg.Wait()
	if result == -1 {
		return result, msg
	}
	msg = "保存成功！"
	return 1, msg

}
