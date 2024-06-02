package mysqldb

import (
	"changeme/services/publicCode"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
)

var DB *sql.DB

func ConnectDB() {

	var err error
	// 准备连接MySQL数据库的信息
	if publicCode.ModeDB {
		DB, err = sql.Open("mysql", publicCode.UserNameDB+":"+publicCode.PassWordDB+"@tcp("+publicCode.HostDB+":"+publicCode.PortDB+")/pppscan?charset=utf8mb4")
		DB.SetConnMaxLifetime(time.Minute * 60)
		DB.SetMaxOpenConns(200)   //  最大连接数
		DB.SetConnMaxIdleTime(50) //   最大空闲连接数
		if err != nil {
			fmt.Println("连接数据库失败：", err)
		}

		err = DB.Ping()
		if err != nil {
			fmt.Println("连接数据库失败：", err)
		}
	} else {
		// 准备连接SQLite数据库的信息
		DB, err = sql.Open("sqlite3", "./pppscan.db?charset=utf8mb4")
		DB.SetConnMaxLifetime(time.Minute * 60)
		DB.SetMaxOpenConns(200)   //  最大连接数
		DB.SetConnMaxIdleTime(50) //   最大空闲连接数
		if err != nil {
			fmt.Println("连接数据库失败：", err)
		}

		err = DB.Ping()
		if err != nil {
			fmt.Println("连接数据库失败：", err)
		}
	}

}

func InitializeDB() (int, string) {

	var errStr string
	// 准备连接MySQL数据库的信息

	if publicCode.ModeDB {
		db1, err := sql.Open("mysql", publicCode.UserNameDB+":"+publicCode.PassWordDB+"@tcp("+publicCode.HostDB+":"+publicCode.PortDB+")/?charset=utf8mb4")
		db1.SetConnMaxLifetime(time.Second * 10)
		if err != nil {
			errStr = "连接数据库失败!"
			fmt.Println("连接数据库失败：", err)
			return -1, errStr
		}

		defer db1.Close()
		sqlStr := "CREATE DATABASE `pppscan` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;"
		_, err = db1.Exec(sqlStr)
		if err != nil {
			errStr := fmt.Sprintf("初始化出错！%s", err)
			if !strings.Contains(errStr, "database exists") {
				return -1, errStr
			}
		}

		// 准备连接MySQL数据库的信息
		db2, err := sql.Open("mysql", publicCode.UserNameDB+":"+publicCode.PassWordDB+"@tcp("+publicCode.HostDB+":"+publicCode.PortDB+")/pppscan?charset=utf8mb4")
		db2.SetConnMaxLifetime(time.Second * 10)
		if err != nil {
			errStr = "连接数据库失败!"
			fmt.Println("连接数据库失败：", err)
			return -1, errStr
		}

		defer db2.Close()
		sqlStr1 := "CREATE TABLE `poc` (\n  `uuid` varchar(100) NOT NULL,\n  `name` varchar(100) NOT NULL,\n  `hunter` varchar(100) DEFAULT NULL,\n  `fofa` varchar(100) DEFAULT NULL,\n  `cms` varchar(100) NOT NULL,\n  `description` varchar(300) DEFAULT NULL,\n  `optionValue` int(11) NOT NULL,\n  `needData` varchar(2000) DEFAULT NULL,\n  `request` text CHARACTER SET utf8 NOT NULL,\n  UNIQUE KEY `poc_uuid_IDX` (`uuid`) USING BTREE\n) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COMMENT='存放poc';"
		sqlStr2 := "-- pppscan.fingerprint definition\n\nCREATE TABLE `fingerprint` (\n  `uuid` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `fingerprintScan` varchar(5000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,\n  `pocsInfo` varchar(5000) COLLATE utf8mb4_unicode_ci DEFAULT NULL\n) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

		_, err = db2.Exec(sqlStr2)
		if err != nil {
			errStr := fmt.Sprintf("初始化出错！%s", err)
			return -1, errStr
		}

		_, err = db2.Exec(sqlStr1)
		if err != nil {
			errStr := fmt.Sprintf("初始化出错！%s", err)
			return -1, errStr
		}
	} else {
		// 准备连接MySQL数据库的信息
		db2, err := sql.Open("sqlite3", "./pppscan.db?charset=utf8mb4")
		db2.SetConnMaxLifetime(time.Second * 10)
		if err != nil {
			errStr = "连接数据库失败!"
			fmt.Println("连接数据库失败：", err)
			return -1, errStr
		}

		defer db2.Close()
		sqlStr1 := "-- poc definition\n\nCREATE TABLE poc (\n\tuuid TEXT NOT NULL,\n\tname TEXT NOT NULL,\n\thunter TEXT,\n\tfofa TEXT,\n\tcms TEXT NOT NULL,\n\tdescription TEXT,\n\toptionValue INTEGER NOT NULL,\n\tneedData TEXT,\n\trequest INTEGER\n);"
		sqlStr2 := "-- fingerprint definition\n\nCREATE TABLE fingerprint (\n\tuuid TEXT NOT NULL,\n\tname TEXT NOT NULL,\n\tfingerprintScan TEXT,\n\tpocsInfo TEXT\n);"

		_, err = db2.Exec(sqlStr2)
		if err != nil {
			errStr := fmt.Sprintf("初始化出错！%s", err)
			return -1, errStr
		}

		_, err = db2.Exec(sqlStr1)
		if err != nil {
			errStr := fmt.Sprintf("初始化出错！%s", err)
			return -1, errStr
		}
	}

	return 1, "初始化成功！"
}
