package config

import (
	"changeme/services/publicCode"
	"fmt"
	"github.com/spf13/viper"
)

func readConfig(config *viper.Viper) {
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Printf("解析配置文件出错: %v\n", err)
		}
	}

	//mysql
	publicCode.ModeDB = config.GetBool("db.mode")

	//if publicCode.ModeDB {
	//	publicCode.HostDB = config.GetString("db.hostPublic")
	//	publicCode.PortDB = config.GetString("db.portPublic")
	//	publicCode.UserNameDB = config.GetString("db.usernamePublic")
	//	publicCode.PassWordDB = config.GetString("db.passwordPublic")
	//} else {
	//	publicCode.HostDB = config.GetString("db.host")
	//	publicCode.PortDB = config.GetString("db.port")
	//	publicCode.UserNameDB = config.GetString("db.username")
	//	publicCode.PassWordDB = config.GetString("db.password")
	//}

	publicCode.HostDB = config.GetString("db.host")
	publicCode.PortDB = config.GetString("db.port")
	publicCode.UserNameDB = config.GetString("db.username")
	publicCode.PassWordDB = config.GetString("db.password")

	//proxy
	publicCode.ModeProxy = config.GetString("proxy.mode")
	publicCode.HostProxy = config.GetString("proxy.host")
	publicCode.PortProxy = config.GetString("proxy.port")
	publicCode.UserNameProxy = config.GetString("proxy.username")
	publicCode.PassWordProxy = config.GetString("proxy.password")
	publicCode.EnableProxy = config.GetBool("proxy.enable")

	//scan
	publicCode.TimeOut = config.GetInt64("scan.timeout")
	publicCode.ThreadNum = config.GetInt("scan.threadNum")

	//other
	publicCode.LocalLoading = config.GetBool("other.localLoading")
	publicCode.OnlyLocalLoading = config.GetBool("other.onlyLocalLoading")

}

func Setting() {

	config := viper.New()
	config.AddConfigPath(".")      // 直接运行用的 文件所在目录
	config.SetConfigName("config") // 文件名
	config.SetConfigType("yaml")   // 文件类型
	readConfig(config)

}

func readConfigAll(config *viper.Viper) publicCode.AllSetting {

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Printf("解析配置文件出错: %v\n", err)
		}
	}

	var settingAll publicCode.AllSetting

	//mysql
	settingAll.DB.Mode = config.GetBool("db.mode")
	//settingAll.DB.HostPublic = config.GetString("db.hostPublic")
	//settingAll.DB.PortPublic = config.GetString("db.portPublic")
	//settingAll.DB.UsernamePublic = config.GetString("db.usernamePublic")
	//settingAll.DB.PasswordPublic = config.GetString("db.passwordPublic")
	settingAll.DB.Host = config.GetString("db.host")
	settingAll.DB.Port = config.GetString("db.port")
	settingAll.DB.UserName = config.GetString("db.username")
	settingAll.DB.Password = config.GetString("db.password")

	//proxy
	settingAll.Proxy.Mode = config.GetString("proxy.mode")
	settingAll.Proxy.Host = config.GetString("proxy.host")
	settingAll.Proxy.Port = config.GetString("proxy.port")
	settingAll.Proxy.Username = config.GetString("proxy.username")
	settingAll.Proxy.Password = config.GetString("proxy.password")
	settingAll.Proxy.Enable = config.GetBool("proxy.enable")

	//scan
	settingAll.Scan.TimeOut = config.GetInt64("scan.timeout")
	settingAll.Scan.ThreadNum = config.GetInt("scan.threadNum")

	// other
	settingAll.Other.LocalLoading = config.GetBool("other.localLoading")
	//settingAll.Other.OnlyLocalLoading = config.GetBool("other.onlyLocalLoading")

	return settingAll

}

func backgroundConfig(config *viper.Viper) {
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Printf("解析配置文件出错: %v\n", err)
		}
	}
	publicCode.BackgroundImage = config.GetString("image")
}

func BackgroundSetting() {

	config := viper.New()
	config.AddConfigPath(".")               // 直接运行用的 文件所在目录
	config.SetConfigName("background.yaml") // 文件名
	config.SetConfigType("yaml")            // 文件类型
	backgroundConfig(config)
}
