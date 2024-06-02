package config

import (
	"changeme/services/mysqldb"
	"changeme/services/pocScan"
	"changeme/services/publicCode"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SettingSave(setting publicCode.AllSetting, dbChange int) (int, string) {

	var msg string
	//创建yaml文件，并将数据写入
	yamlData, err := yaml.Marshal(setting)
	if err != nil {
		msg = fmt.Sprintf("Error marshaling YAML:%s", err)
		return -1, msg
	}

	fileName := "config.yaml"
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		msg = fmt.Sprintf("error getting current working directory: %s", err)
		return -1, msg
	}

	// 构建目标目录路径
	targetDir := filepath.Join(currentDir, "/")

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

	Setting()
	pocScan.CreatPocScanClient()
	if dbChange == 1 {
		mysqldb.ConnectDB()
	}
	return 1, "保存成功！"

}

func SettingShow() publicCode.AllSetting {

	config := viper.New()
	config.AddConfigPath(".")      // 直接运行用的 文件所在目录
	config.SetConfigName("config") // 文件名
	config.SetConfigType("yaml")   // 文件类型
	settingAll := readConfigAll(config)
	return settingAll

}

//func BackgroundImageCreat(backgroundImageURl string) (int, string) {
//
//	images := strings.SplitN(backgroundImageURl, ";base64,", 2)
//
//	image := images[1]
//	imageBety, err := base64.StdEncoding.DecodeString(image)
//	if err != nil {
//		fmt.Println("解码错误", err)
//	}
//
//	var msg string
//
//	fileName := "/frontend/src/assets/images/background-image.png"
//	// 获取当前工作目录
//	currentDir, err := os.Getwd()
//	if err != nil {
//		msg = fmt.Sprintf("error getting current working directory: %s", err)
//		return -1, msg
//	}
//
//	// 构建目标目录路径
//	targetDir := filepath.Join(currentDir, "/")
//
//	// 创建目标目录（如果不存在）
//	err = os.MkdirAll(targetDir, 0755)
//	if err != nil {
//		msg = fmt.Sprintf("error creating target directory: %s", err)
//		return -1, msg
//	}
//
//	// 构建文件路径
//	filePath := filepath.Join(targetDir, fileName)
//
//	// 写入文件内容
//	err = ioutil.WriteFile(filePath, imageBety, 0644)
//	if err != nil {
//		msg = fmt.Sprintf("error writing file: %s", err)
//		return -1, msg
//	}
//
//	//// 写入文件内容
//	//err = ioutil.WriteFile(filePath, []byte(image), 0644)
//	//if err != nil {
//	//	msg = fmt.Sprintf("error writing file: %s", err)
//	//	return -1, msg
//	//}
//
//	return 1, "设置成功！"
//
//}

func BackgroundSave(image publicCode.Background) (int, string) {

	var msg string
	//创建yaml文件，并将数据写入
	yamlData, err := yaml.Marshal(image)
	if err != nil {
		msg = fmt.Sprintf("Error marshaling YAML:%s", err)
		return -1, msg
	}

	fileName := "background.yaml"
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		msg = fmt.Sprintf("error getting current working directory: %s", err)
		return -1, msg
	}

	// 构建目标目录路径
	targetDir := filepath.Join(currentDir, "/")

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

	BackgroundSetting()
	return 1, "保存成功！"

}
