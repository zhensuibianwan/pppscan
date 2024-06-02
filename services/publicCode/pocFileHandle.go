package publicCode

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// GetAllFileNames 获取指定目录下的所有文件名
func GetAllFileNames() ([]string, error) {

	var fileNames []string

	files, err := ioutil.ReadDir("./poc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		fileName := strings.TrimRight(file.Name(), ".yaml")
		fileNames = append(fileNames, fileName)
	}

	return fileNames, nil

}

func CheckFileName(key string) int {
	var wg sync.WaitGroup
	result := -1
	filenames, err := GetAllFileNames()
	if err != nil {
		return 0
	}
	for _, filename := range filenames {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			if filename == key {
				result = 1
			}
		}(filename)
	}
	wg.Wait()
	return result
}

func readContent(config *viper.Viper) Poc {
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..", config)
			return Poc{}
		} else {
			fmt.Printf("解析配置文件出错: %v\n", err)
			return Poc{}
		}
	}
	var poc Poc

	poc.UUID = config.GetString("uuid")
	poc.Hunter = config.GetString("hunter")
	poc.Fofa = config.GetString("fofa")
	poc.CMS = config.GetString("cms")
	poc.Description = config.GetString("description")
	poc.OptionValue = config.GetInt("optionValue")
	needData := config.Get("needData").([]interface{})
	for _, need := range needData {
		if itemMap, ok := need.(map[string]interface{}); ok {
			var newItem Need
			if label, ok := itemMap["label"].(string); ok {
				newItem.Label = label
			}
			if value, ok := itemMap["value"].(string); ok {
				newItem.Value = value
			}
			poc.NeedData = append(poc.NeedData, newItem)
		}
	}

	requestData := config.Get("request").([]interface{})

	for _, item := range requestData {
		if itemMap, ok := item.(map[string]interface{}); ok {
			var newItem RequestData
			if pocString, ok := itemMap["pocstring"].(string); ok {
				newItem.PocString = pocString
			}
			if status, ok := itemMap["status"].(string); ok {
				newItem.Status = status
			}

			if check, ok := itemMap["check"].(string); ok {
				newItem.Check = check
			}
			if print1, ok := itemMap["print"].(string); ok {
				newItem.Print = print1
			}
			poc.Request = append(poc.Request, newItem)
		}
	}
	return poc

}

func PocFileRead(fileName string) Poc {

	config := viper.New()
	config.AddConfigPath("./poc")  // 直接运行用的 文件所在目录
	config.SetConfigName(fileName) // 文件名
	config.SetConfigType("yaml")   // 文件类型
	poc := readContent(config)
	poc.Name = fileName
	return poc
}

func PocFileList(key string) []Poc {

	var fileNames []string

	AllFileNames, err := GetAllFileNames()
	if err != nil {
		return nil
	}

	for _, name := range AllFileNames {
		if strings.Contains(strings.ToLower(name), strings.ToLower(key)) {
			fileNames = append(fileNames, name)
		}
	}

	num := len(fileNames)
	pocs := make([]Poc, 0, num)
	var wg sync.WaitGroup
	for _, fileName := range fileNames {
		wg.Add(1)
		go func(fileName string) {
			poc := PocFileRead(fileName)
			pocs = append(pocs, poc)
			wg.Done()
		}(fileName)
	}
	wg.Wait()
	sort.Sort(PocSlice(pocs))
	return pocs
}

func RemovePocByName(name string) {

	filePath := fmt.Sprintf("./poc/%s.yaml", name)
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("本地文件删除失败！")
	}
}

func SavePocToFile(poc Poc) {
	//创建yaml文件，并将数据写入
	yamlData, err := yaml.Marshal(poc)
	if err != nil {
		fmt.Println("Error marshaling YAML: ", err)
		return
	}

	fileName := poc.Name + ".yaml"
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current working directory: ", err)
		return
	}

	// 构建目标目录路径
	targetDir := filepath.Join(currentDir, "/poc")

	// 创建目标目录（如果不存在）
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Println("error creating target directory: ", err)
		return
	}

	// 构建文件路径
	filePath := filepath.Join(targetDir, fileName)

	// 写入文件内容
	err = ioutil.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		fmt.Println("error writing file: ", err)
		return
	}

}
