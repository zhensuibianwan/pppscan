package mysqldb

import (
	"changeme/services/publicCode"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func SearchPocByName(nameKey string) []publicCode.Poc {

	sqlStr := "select uuid,name,hunter,fofa,cms,description,optionValue,needData,request from poc where name like ?"

	if DB == nil {
		ConnectDB()
	}

	db := DB

	var count int
	countQuery := "SELECT COUNT(*) FROM poc WHERE name LIKE ?"
	err := db.QueryRow(countQuery, "%"+nameKey+"%").Scan(&count)
	if err != nil {
		fmt.Printf("countQuery failed, err:%v\n", err)
		return nil
	}

	rows, err := db.Query(sqlStr, "%"+nameKey+"%")
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	pocs := make([]publicCode.Poc, count)
	i := 0
	// 循环读取结果集中的数据
	for rows.Next() {
		var needData string
		var request string
		var needByte []byte
		var requestByte []byte

		err = rows.Scan(&pocs[i].UUID, &pocs[i].Name, &pocs[i].Hunter, &pocs[i].Fofa, &pocs[i].CMS, &pocs[i].Description, &pocs[i].OptionValue, &needData, &request)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil
		}

		needByte, err = base64.StdEncoding.DecodeString(needData)
		if err != nil {
			fmt.Printf("Base64解码错误！err is %v", err)
		}

		err = json.Unmarshal(needByte, &pocs[i].NeedData)
		if err != nil {
			fmt.Printf("反序列化错误 err is %v	Need name:%s", err, pocs[i].Name)
			poc := publicCode.PocFileRead(pocs[i].Name)
			if len(poc.NeedData) > 0 {
				pocs[i].NeedData = poc.NeedData
				UpdatePocDB(pocs[i])
			}
		}

		requestByte, err = base64.StdEncoding.DecodeString(request)
		if err != nil {
			fmt.Printf("Base64解码错误！err is %v", err)
		}
		err = json.Unmarshal(requestByte, &pocs[i].Request)
		if err != nil {
			fmt.Printf("反序列化错误 err is %v	Request name:%s", err, pocs[i].Name)
			poc := publicCode.PocFileRead(pocs[i].Name)
			if len(poc.Request) > 0 {
				pocs[i].Request = poc.Request
				UpdatePocDB(pocs[i])
			}

		}
		i++
	}

	return pocs
}

func SearchPocByUUID(UUID string) []publicCode.Poc {

	sqlStr := "select uuid,name,hunter,fofa,cms,description,optionValue,needData,request from poc where UUID = ?"
	db := DB

	var count int
	countQuery := "SELECT COUNT(*) FROM poc WHERE UUID = ?"
	err := db.QueryRow(countQuery, UUID).Scan(&count)
	if err != nil {
		fmt.Printf("countQuery failed, err:%v\n", err)
		return nil
	}

	rows, err := db.Query(sqlStr, UUID)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()
	pocs := make([]publicCode.Poc, count)
	if count != 0 {
		i := 0
		// 循环读取结果集中的数据
		for rows.Next() {
			var needData string
			var request string
			var needByte []byte
			var requestByte []byte

			err = rows.Scan(&pocs[i].UUID, &pocs[i].Name, &pocs[i].Hunter, &pocs[i].Fofa, &pocs[i].CMS, &pocs[i].Description, &pocs[i].OptionValue, &needData, &request)
			if err != nil {
				fmt.Printf("scan failed, err:%v\n", err)
				return nil
			}
			needByte, err = base64.StdEncoding.DecodeString(needData)
			if err != nil {
				fmt.Printf("Base64解码错误！err is %v", err)
			}

			err = json.Unmarshal(needByte, &pocs[i].NeedData)
			if err != nil {
				fmt.Printf("反序列化错误 err is %v", err)
			}

			requestByte, err = base64.StdEncoding.DecodeString(request)
			if err != nil {
				fmt.Printf("Base64解码错误！err is %v", err)
			}
			err = json.Unmarshal(requestByte, &pocs[i].Request)
			if err != nil {
				fmt.Printf("反序列化错误 err is %v", err)
			}
			i++
		}
	}

	return pocs
}

func CheckPocByUUID(UUID string) int {

	db := DB

	count := 1
	countQuery := "SELECT COUNT(*) FROM poc WHERE UUID = ?"
	err := db.QueryRow(countQuery, UUID).Scan(&count)
	if err != nil {
		fmt.Printf("countQuery failed, err:%v\n", err)
	}

	return count
}

// SavePocDB InsertPoc 将poc保存到数据库
func SavePocDB(poc publicCode.Poc) (int, string) {

	var errStr string
	requestByte, err := json.Marshal(poc.Request)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	requestData := base64.StdEncoding.EncodeToString(requestByte)

	NeedByte, err := json.Marshal(poc.NeedData)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	NeedData := base64.StdEncoding.EncodeToString(NeedByte)

	sqlStr := "INSERT INTO poc (uuid,name,hunter,fofa,cms,description,optionValue,needData,request) VALUES (?,?,?,?,?,?,?,?,?)"
	db := DB
	stmt, err := db.Prepare(sqlStr)

	if err != nil {
		errStr = "数据库执行错误！"
		fmt.Printf("db.Prepare执行错误！ err is %v", err)
		return -1, errStr
	}
	defer stmt.Close()
	_, err = stmt.Exec(poc.UUID, poc.Name, poc.Hunter, poc.Fofa, poc.CMS, poc.Description, poc.OptionValue, string(NeedData), string(requestData))
	if err != nil {
		errStr = "添加数据失败！"
		fmt.Printf("db.Exec执行错误！ err is %v", err)
		return -1, errStr
	}

	return 1, "数据库写入成功！"

}

// UpdatePocDB InsertPoc 将poc保存到数据库
func UpdatePocDB(poc publicCode.Poc) (int, string) {

	var errStr string
	requestByte, err := json.Marshal(poc.Request)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	requestData := base64.StdEncoding.EncodeToString(requestByte)

	NeedByte, err := json.Marshal(poc.NeedData)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	NeedData := base64.StdEncoding.EncodeToString(NeedByte)

	sqlStr := "UPDATE poc set name=?,hunter=?,fofa=?,cms=?,description=?,optionValue=?,needData=?,request=? WHERE uuid=?"
	db := DB
	stmt, err := db.Prepare(sqlStr)

	if err != nil {
		errStr = "数据库执行错误！"
		fmt.Printf("db.Prepare执行错误！ err is %v", err)
		return -1, errStr
	}
	defer stmt.Close()
	_, err = stmt.Exec(poc.Name, poc.Hunter, poc.Fofa, poc.CMS, poc.Description, poc.OptionValue, string(NeedData), string(requestData), poc.UUID)
	if err != nil {
		errStr = "修改数据失败！"
		fmt.Printf("db.Exec执行错误！ err is %v", err)
		return -1, errStr
	}

	return 1, "数据库修改成功！"

}

// DelPocDB 将poc从数据库中删除
func DelPocDB(pocs []publicCode.Poc) (int, string) {
	var errStr string

	sqlStr := "delete from poc where uuid = ?"
	db := DB
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		errStr = "数据库执行错误！"
		fmt.Printf("db.Prepare执行错误！ err is %v", err)
		return -1, errStr
	}

	defer stmt.Close()
	for _, poc := range pocs {
		_, err := stmt.Exec(poc.UUID)
		if err != nil {
			errStr = "数据库执行错误！"
			fmt.Printf("db.Exec执行错误！ err is %v", err)
			return -1, errStr
		}
	}

	return 1, "删除成功！"
}
