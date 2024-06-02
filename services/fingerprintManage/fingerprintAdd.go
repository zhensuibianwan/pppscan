package fingerprintManage

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func SaveFingerprint(fingerprint publicCode.Fingerprint) (int, string) {

	checkFingerprint := mysqldb.SearchFingerprintByName(fingerprint.Name)
	if (len(checkFingerprint) == 1) && checkFingerprint[0].UUID != fingerprint.UUID && checkFingerprint[0].Name == fingerprint.Name {
		return -1, "存在相同的名称"
	}

	//将数据保存到数据库
	var msg string
	var status int
	if fingerprint.UUID == "" {
		currentTime := time.Now().Format("20060102150405")
		time.Sleep(time.Microsecond * 200)
		randNum := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(900) + 100
		id := currentTime + strconv.Itoa(int(randNum))
		fingerprint.UUID = id
		status, msg = mysqldb.SaveFingerprintDB(fingerprint)
		if status != 1 {
			return status, msg
		}
	} else {
		checkFingerprint = mysqldb.SearchFingerprintByUUID(fingerprint.UUID)
		status, msg = mysqldb.UpdateFingerprintDB(fingerprint)
		if status != 1 {
			return status, msg
		}
		if checkFingerprint[0].Name != fingerprint.Name && checkFingerprint[0].UUID == fingerprint.UUID {
			publicCode.RemoveFingerprintFileByName(checkFingerprint[0].Name)
			if len(fingerprint.PocsInfo) > 0 {
				for _, pocInfo := range fingerprint.PocsInfo {
					poc := mysqldb.SearchPocByUUID(pocInfo.UUID)
					poc[0].CMS = fingerprint.Name
					status, msg = mysqldb.UpdatePocDB(poc[0])
					publicCode.SavePocToFile(poc[0])
				}
			}
		}
	}

	//创建json文件，并将数据写入
	jsonData, err := json.Marshal(fingerprint)
	if err != nil {
		msg = fmt.Sprintf("Error marshaling YAML:%s", err)
		return -1, msg
	}

	fileName := fingerprint.Name + ".json"
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		msg = fmt.Sprintf("error getting current working directory: %s", err)
		return -1, msg
	}

	// 构建目标目录路径
	targetDir := filepath.Join(currentDir, "/fingerprint")

	// 创建目标目录（如果不存在）
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		msg = fmt.Sprintf("error creating target directory: %s", err)
		return -1, msg
	}

	// 构建文件路径
	filePath := filepath.Join(targetDir, fileName)

	// 写入文件内容
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		msg := fmt.Sprintf("error writing file: %s", err)
		return -1, msg
	}

	return 1, "保存成功！"
}
