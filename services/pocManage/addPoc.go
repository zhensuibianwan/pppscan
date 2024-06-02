package pocManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func SavePoc(poc publicCode.Poc) (int, string) {

	checkPoc := mysqldb.SearchPocByName(poc.Name)
	if (len(checkPoc) > 0) && checkPoc[0].UUID != poc.UUID && checkPoc[0].Name == poc.Name {
		return -1, "存在相同的文件名"
	}
	poc.Request = poc.Request[0:poc.OptionValue]

	//将数据保存到数据库
	var msg string
	var status int
	if poc.UUID == "" {
		currentTime := time.Now().Format("20060102150405")
		time.Sleep(time.Microsecond * 200)
		randNum := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(900) + 100
		id := currentTime + strconv.Itoa(int(randNum))
		poc.UUID = id
		status, msg = mysqldb.SavePocDB(poc)
		if status != 1 {
			return status, msg
		}

	} else {
		checkPoc = mysqldb.SearchPocByUUID(poc.UUID)
		status, msg = mysqldb.UpdatePocDB(poc)
		if status != 1 {
			return status, msg
		}
		if checkPoc[0].Name != poc.Name {
			publicCode.RemovePocByName(checkPoc[0].Name)
		}

	}

	fingerprint := mysqldb.SearchFingerprintByName(poc.CMS)
	if len(fingerprint) == 1 {
		checkLinkPoc := 1
		for _, pocsInfo := range fingerprint[0].PocsInfo {
			if pocsInfo.UUID == poc.UUID {
				checkLinkPoc = 0
			}
		}
		if checkLinkPoc == 1 {
			var pocsInfo publicCode.PocsInfoData
			pocsInfo.Name = poc.Name
			pocsInfo.UUID = poc.UUID
			fingerprint[0].PocsInfo = append(fingerprint[0].PocsInfo, pocsInfo)
			status, msg = mysqldb.UpdateFingerprintDB(fingerprint[0])
			if status != 1 {
				return status, msg
			}
			publicCode.SaveFingerprintToFile(fingerprint[0])
		}
	}

	//创建yaml文件，并将数据写入
	yamlData, err := yaml.Marshal(poc)
	if err != nil {
		msg = fmt.Sprintf("Error marshaling YAML:%s", err)
		return -1, msg
	}

	fileName := poc.Name + ".yaml"
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		msg = fmt.Sprintf("error getting current working directory: %s", err)
		return -1, msg
	}

	// 构建目标目录路径
	targetDir := filepath.Join(currentDir, "/poc")

	// 创建目标目录（如果不存在）
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		msg = fmt.Sprintf("error creating target directory: %s", err)
		return -1, msg
	}

	// 构建文件路径
	filePath := filepath.Join(targetDir, fileName)

	// 写入文件内容
	err = ioutil.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		msg := fmt.Sprintf("error writing file: %s", err)
		return -1, msg
	}

	return 1, "保存成功！"
}
