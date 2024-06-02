package publicCode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func GetFingerprintAllFileNames() ([]string, error) {

	var fileNames []string

	files, err := ioutil.ReadDir("./fingerprint")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil

}

func FingerprintFileRead(fileName string) Fingerprint {

	filePath := fmt.Sprintf("./fingerprint/%s", fileName)

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件时发生错误:", err)
		return Fingerprint{}
	}

	var fingerprint Fingerprint
	err = json.Unmarshal(data, &fingerprint)
	if err != nil {
		fmt.Println("json反序列化时出错：", fileName, err)
		return Fingerprint{}
	}

	return fingerprint

}

func FingerprintFileList(key string) []Fingerprint {

	var fileNames []string

	AllFileNames, err := GetFingerprintAllFileNames()
	if err != nil {
		return nil
	}

	for _, name := range AllFileNames {
		if strings.Contains(strings.ToLower(name), strings.ToLower(key)) {
			fileNames = append(fileNames, name)
		}
	}

	num := len(fileNames)
	fingerprints := make([]Fingerprint, 0, num)
	var wg sync.WaitGroup
	for _, fileName := range fileNames {
		wg.Add(1)
		go func(fileName string) {
			fingerprint := FingerprintFileRead(fileName)
			fingerprints = append(fingerprints, fingerprint)
			wg.Done()
		}(fileName)
	}
	wg.Wait()
	sort.Sort(FingerprintSlice(fingerprints))
	return fingerprints
}

func RemoveFingerprintFileByName(name string) {

	filePath := fmt.Sprintf("./fingerprint/%s.json", name)
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("本地文件删除失败！")
	}
}

func SaveFingerprintToFile(fingerprint Fingerprint) {

	//创建json文件，并将数据写入
	jsonData, err := json.Marshal(fingerprint)
	if err != nil {
		fmt.Println("Error marshaling YAML: ", err)
		return
	}

	fileName := fingerprint.Name + ".json"
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current working directory: ", err)
		return
	}

	// 构建目标目录路径
	targetDir := filepath.Join(currentDir, "/fingerprint")

	// 创建目标目录（如果不存在）
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Println("error creating target directory: ", err)
		return
	}

	// 构建文件路径
	filePath := filepath.Join(targetDir, fileName)

	// 写入文件内容
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		fmt.Println("error writing file: ", err)
		return
	}

}
