package pocManage

import (
	"changeme/services/publicCode"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func DownLoadPocs(pocs []publicCode.Poc) (int, string) {

	var msg string
	var wg sync.WaitGroup
	var result int

	for _, poc := range pocs {
		wg.Add(1)
		go func(poc publicCode.Poc) {
			defer wg.Done()
			//创建yaml文件，并将数据写入
			yamlData, err := yaml.Marshal(poc)
			if err != nil {
				msg = fmt.Sprintf("Error marshaling YAML:%s", err)
				result = -1
				return
			}

			fileName := poc.Name + ".yaml"
			// 获取当前工作目录
			currentDir, err := os.Getwd()
			if err != nil {
				msg = fmt.Sprintf("error getting current working directory: %s", err)
				result = -1
				return
			}

			// 构建目标目录路径
			targetDir := filepath.Join(currentDir, "/poc")

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
			err = ioutil.WriteFile(filePath, yamlData, 0644)
			if err != nil {
				msg = fmt.Sprintf("error writing file: %s", err)
				result = -1
				return
			}
		}(poc)

	}
	wg.Wait()
	if result == -1 {
		return result, msg
	}
	msg = "保存成功！"
	return 1, msg

}
