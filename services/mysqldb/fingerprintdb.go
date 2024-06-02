package mysqldb

import (
	"changeme/services/publicCode"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func SearchFingerprintByName(nameKey string) []publicCode.Fingerprint {

	sqlStr := "select uuid,name,fingerprintScan,pocsInfo from fingerprint where name like ?"

	if DB == nil {
		ConnectDB()
	}

	db := DB

	var count int
	countQuery := "SELECT COUNT(*) FROM fingerprint WHERE name LIKE ?"
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

	fingerprints := make([]publicCode.Fingerprint, count)
	i := 0
	// 循环读取结果集中的数据
	for rows.Next() {
		var fingerprintScan string
		var pocsInfo string
		var fingerprintScanByte []byte
		var pocsInfoByte []byte

		err = rows.Scan(&fingerprints[i].UUID, &fingerprints[i].Name, &fingerprintScan, &pocsInfo)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil
		}

		fingerprintScanByte, err = base64.StdEncoding.DecodeString(fingerprintScan)
		if err != nil {
			fmt.Printf("Base64解码错误！err is %v", err)
		}

		err = json.Unmarshal(fingerprintScanByte, &fingerprints[i].FingerprintScan)
		if err != nil {
			fmt.Printf("反序列化错误 err is %v	Need name:%s", err, fingerprints[i].Name)
			fingerprint := publicCode.FingerprintFileRead(fingerprints[i].Name + ".json")
			if len(fingerprint.FingerprintScan) > 0 {
				fingerprints[i].FingerprintScan = fingerprint.FingerprintScan
				UpdateFingerprintDB(fingerprints[i])
			}
		}

		pocsInfoByte, err = base64.StdEncoding.DecodeString(pocsInfo)
		if err != nil {
			fmt.Printf("Base64解码错误！err is %v", err)
		}
		err = json.Unmarshal(pocsInfoByte, &fingerprints[i].PocsInfo)
		if err != nil {
			fmt.Printf("反序列化错误 err is %v	Request name:%s", err, fingerprints[i].Name)
			fingerprint := publicCode.FingerprintFileRead(fingerprints[i].Name + ".json")
			if len(fingerprint.PocsInfo) > 0 {
				fingerprints[i].PocsInfo = fingerprint.PocsInfo
				UpdateFingerprintDB(fingerprints[i])
			}

		}
		i++
	}

	return fingerprints
}

func SearchFingerprintByUUID(UUID string) []publicCode.Fingerprint {

	sqlStr := "select uuid,name,fingerprintScan,pocsInfo from fingerprint where uuid = ?"

	if DB == nil {
		ConnectDB()
	}

	db := DB

	var count int
	countQuery := "SELECT COUNT(*) FROM fingerprint WHERE uuid = ?"
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

	fingerprints := make([]publicCode.Fingerprint, count)
	i := 0
	// 循环读取结果集中的数据
	for rows.Next() {
		var fingerprintScan string
		var pocsInfo string
		var fingerprintScanByte []byte
		var pocsInfoByte []byte

		err = rows.Scan(&fingerprints[i].UUID, &fingerprints[i].Name, &fingerprintScan, &pocsInfo)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil
		}

		fingerprintScanByte, err = base64.StdEncoding.DecodeString(fingerprintScan)
		if err != nil {
			fmt.Printf("Base64解码错误！err is %v", err)
		}

		err = json.Unmarshal(fingerprintScanByte, &fingerprints[i].FingerprintScan)
		if err != nil {
			fmt.Printf("反序列化错误 err is %v	Need name:%s", err, fingerprints[i].Name)
			fingerprint := publicCode.FingerprintFileRead(fingerprints[i].Name + ".json")
			if len(fingerprint.FingerprintScan) > 0 {
				fingerprints[i].FingerprintScan = fingerprint.FingerprintScan
				UpdateFingerprintDB(fingerprints[i])
			}
		}

		pocsInfoByte, err = base64.StdEncoding.DecodeString(pocsInfo)
		if err != nil {
			fmt.Printf("Base64解码错误！err is %v", err)
		}
		err = json.Unmarshal(pocsInfoByte, &fingerprints[i].PocsInfo)
		if err != nil {
			fmt.Printf("反序列化错误 err is %v	Request name:%s", err, fingerprints[i].Name)
			fingerprint := publicCode.FingerprintFileRead(fingerprints[i].Name + ".json")
			if len(fingerprint.PocsInfo) > 0 {
				fingerprints[i].PocsInfo = fingerprint.PocsInfo
				UpdateFingerprintDB(fingerprints[i])
			}
		}
		i++
	}

	return fingerprints
}

func SearchFingerprintGetNameByName(nameKey string) []string {

	sqlStr := "select name from fingerprint where name like ?"

	if DB == nil {
		ConnectDB()
	}

	db := DB

	var count int
	countQuery := "SELECT COUNT(*) FROM fingerprint WHERE name LIKE ?"
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

	names := make([]string, count)
	i := 0
	// 循环读取结果集中的数据
	for rows.Next() {

		err = rows.Scan(&names[i])
		if err != nil {
			fmt.Printf("rows failed, err:%v\n", err)
			return nil
		}
		i++

	}

	return names
}

func CheckFingerprintByUUID(UUID string) int {

	db := DB

	count := 1
	countQuery := "SELECT COUNT(*) FROM fingerprint WHERE UUID = ?"
	err := db.QueryRow(countQuery, UUID).Scan(&count)
	if err != nil {
		fmt.Printf("countQuery failed, err:%v\n", err)
	}

	return count
}

// SaveFingerprintDB  将poc保存到数据库
func SaveFingerprintDB(fingerprint publicCode.Fingerprint) (int, string) {

	var errStr string
	fingerprintScanByte, err := json.Marshal(fingerprint.FingerprintScan)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	fingerprintScanData := base64.StdEncoding.EncodeToString(fingerprintScanByte)

	pocsInfoByte, err := json.Marshal(fingerprint.PocsInfo)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	pocsInfoData := base64.StdEncoding.EncodeToString(pocsInfoByte)

	sqlStr := "INSERT INTO fingerprint (uuid,name,fingerprintScan,pocsInfo) VALUES (?,?,?,?)"
	db := DB
	stmt, err := db.Prepare(sqlStr)

	if err != nil {
		errStr = "数据库执行错误！"
		fmt.Printf("db.Prepare执行错误！ err is %v", err)
		return -1, errStr
	}
	defer stmt.Close()
	_, err = stmt.Exec(fingerprint.UUID, fingerprint.Name, fingerprintScanData, pocsInfoData)
	if err != nil {
		errStr = "添加数据失败！"
		fmt.Printf("db.Exec执行错误！ err is %v", err)
		return -1, errStr
	}

	return 1, "数据库写入成功！"

}

func UpdateFingerprintDB(fingerprint publicCode.Fingerprint) (int, string) {

	var errStr string
	fingerprintScanByte, err := json.Marshal(fingerprint.FingerprintScan)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	fingerprintScanData := base64.StdEncoding.EncodeToString(fingerprintScanByte)

	pocsInfoByte, err := json.Marshal(fingerprint.PocsInfo)
	if err != nil {
		errStr = "序列化错误！"
		fmt.Printf("序列化错误 err is %v", err)
		return -1, errStr
	}

	pocsInfoData := base64.StdEncoding.EncodeToString(pocsInfoByte)

	sqlStr := "UPDATE fingerprint set name=?,fingerprintScan=?, pocsInfo=? WHERE uuid=?"
	db := DB
	stmt, err := db.Prepare(sqlStr)

	if err != nil {
		errStr = "数据库执行错误！"
		fmt.Printf("db.Prepare执行错误！ err is %v", err)
		return -1, errStr
	}
	defer stmt.Close()
	_, err = stmt.Exec(fingerprint.Name, fingerprintScanData, pocsInfoData, fingerprint.UUID)
	if err != nil {
		errStr = "修改数据失败！"
		fmt.Printf("db.Exec执行错误！ err is %v", err)
		return -1, errStr
	}

	return 1, "数据库修改成功！"

}

func DelFingerprintDB(fingerprints []publicCode.Fingerprint) (int, string) {
	var errStr string

	sqlStr := "delete from fingerprint where uuid = ?"
	db := DB
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		errStr = "数据库执行错误！"
		fmt.Printf("db.Prepare执行错误！ err is %v", err)
		return -1, errStr
	}

	defer stmt.Close()
	for _, fingerprint := range fingerprints {
		_, err := stmt.Exec(fingerprint.UUID)
		if err != nil {
			errStr = "数据库执行错误！"
			fmt.Printf("db.Exec执行错误！ err is %v", err)
			return -1, errStr
		}
	}

	return 1, "删除成功！"
}
