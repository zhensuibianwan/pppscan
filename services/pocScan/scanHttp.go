package pocScan

import (
	"bytes"
	"changeme/services/publicCode"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func httpRequest(rule publicCode.Rule, url string) (publicCode.RespondData, error) {

	//定义返回包数据结构体
	var res publicCode.RespondData

	//定义http响应
	var resp *http.Response
	var requ *http.Request
	var err error
	var start time.Time

	URL := url + rule.Path

	rule.Body = strings.ReplaceAll(rule.Body, "\n", "\r\n")
	rule.Body = strings.ReplaceAll(rule.Body, "\r\r\n", "\r\n")

	if strings.Contains(rule.Body, "{{hex_decode(\"") {
		regex := regexp.MustCompile(`{{hex_decode\("(.+?)"\)}}`)
		regexAry := regex.FindAllString(rule.Body, -1)
		ruleBody := strings.ReplaceAll(rule.Body, regexAry[0], "~HexDecodeString~")
		hexDecodeString := strings.ReplaceAll(regexAry[0], "{{hex_decode(\"", "")
		hexDecodeString = strings.ReplaceAll(hexDecodeString, "\")}}", "")
		//ruleBody := strings.ReplaceAll(regexAry[0], "{{hex_decode(\"", "")
		//ruleBody = strings.ReplaceAll(ruleBody, "\")}}", "")
		ruleBodyByte, err := hex.DecodeString(hexDecodeString)
		ruleBody = strings.ReplaceAll(ruleBody, "~HexDecodeString~", string(ruleBodyByte))

		if err != nil {
			fmt.Println("解码失败：", err)
			return res, err
		}
		//requ, err = http.NewRequest(rule.Method, URL, bytes.NewReader(ruleBodyByte))
		start = time.Now()
		requ, err = http.NewRequest(rule.Method, URL, strings.NewReader(ruleBody))
		if err != nil {
			fmt.Printf("get failed, err:%v\n", err)
			return res, err
		}
	} else {
		start = time.Now()
		requ, err = http.NewRequest(rule.Method, URL, strings.NewReader(rule.Body))
		if err != nil {
			fmt.Printf("get failed, err:%v\n", err)
			return res, err
		}
	}

	//设置请求头
	for head := range rule.Headers {
		requ.Header.Set(head, rule.Headers[head])
	}

	resp, err = ClientNoRedirect.Do(requ)
	if err != nil {
		// 检查是否是超时错误
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			//fmt.Sprintf("请求超时！\n")
			err = errors.New("请求超时！")
			return res, err
		}
	}

	end := time.Since(start)

	if resp != nil {
		CheckData := resp.Header.Get("Date")
		ContentLength := resp.Header.Get("Content-Length")
		CheckServer := resp.Header.Get("Server")
		if CheckData != "" || ContentLength != "" || CheckServer != "" {
			defer resp.Body.Close()

			m := make(map[string]string)
			for k, v := range resp.Header {
				m[k] = v[0]
			}
			res.Headers = m
			res.Status = resp.Status

			if strings.Contains(resp.Status, rule.Status) {
				var body []byte
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("fetch: %v\n", err)
					return res, err
				}

				if strings.Contains(rule.Check, "{{Content-Length") {

					var cl int
					if ContentLength != "" {
						cl, err = strconv.Atoi(ContentLength)
						if err != nil {
							fmt.Printf("strconv.Atoi: %v\n", err)
							return res, err
						}
					} else {
						err = errors.New("漏洞不存在！")
						return res, err
					}
					if strings.Contains(rule.Check, "{{Content-LengthMax(") {
						regex := regexp.MustCompile(`{{Content-LengthMax\("(.+?)"\)}}`)
						regexAry := regex.FindAllString(rule.Check, -1)
						CheckContentLength := strings.ReplaceAll(regexAry[0], "{{Content-LengthMax(\"", "")
						CheckContentLength = strings.ReplaceAll(CheckContentLength, "\")}}", "")

						clCheck, errCL := strconv.Atoi(CheckContentLength)
						if errCL != nil {
							fmt.Printf("strconv.Atoi: %v\n", err)
							return res, err
						}
						if cl > clCheck {
							err = errors.New("漏洞不存在！")
							return res, err
						}

					} else if strings.Contains(rule.Check, "{{Content-LengthMin(") {
						regex := regexp.MustCompile(`{{Content-LengthMin\("(.+?)"\)}}`)
						regexAry := regex.FindAllString(rule.Check, -1)
						CheckContentLength := strings.ReplaceAll(regexAry[0], "{{Content-LengthMin(\"", "")
						CheckContentLength = strings.ReplaceAll(CheckContentLength, "\")}}", "")

						clCheck, errCL := strconv.Atoi(CheckContentLength)
						if errCL != nil {
							fmt.Printf("strconv.Atoi: %v\n", err)
							return res, err
						}
						if cl <= clCheck {
							err = errors.New("漏洞不存在！")
							return res, err
						}
					} //else if strings.Contains(rule.Check, "{{Content-Length(") {
					//	regex := regexp.MustCompile(`{{Content-Length\("(.+?)"\)}}`)
					//	regexAry := regex.FindAllString(rule.Check, -1)
					//	CheckContentLength := strings.ReplaceAll(regexAry[0], "{{Content-Length(\"", "")
					//	CheckContentLength = strings.ReplaceAll(CheckContentLength, "\")}}", "")
					//
					//	clCheck, errCL := strconv.Atoi(CheckContentLength)
					//	if errCL != nil {
					//		fmt.Printf("strconv.Atoi: %v\n", err)
					//		return res, err
					//	}
					//	if cl > clCheck {
					//		err = errors.New("漏洞不存在！")
					//		return res, err
					//	}
					//
					//}

				}
				//判断漏洞存在
				if strings.Contains(rule.Check, "&&") {
					checkAry := strings.Split(rule.Check, "&&")
					for i := 0; i < len(checkAry); i++ {
						if strings.Contains(checkAry[i], "{{Time(") {
							ResTime := strings.ReplaceAll(checkAry[i], "{{Time(\"", "")
							ResTime = strings.ReplaceAll(ResTime, "\")}}", "")

							ResTimeCheck, errCL := strconv.Atoi(ResTime)
							if errCL != nil {
								fmt.Printf("strconv.Atoi: %v\n", err)
								return res, err
							}
							if end < time.Duration(ResTimeCheck)*time.Second {
								err = errors.New("漏洞不存在！")
								res.IsCheck = false
								return res, err
							}

						} else if !strings.Contains(string(body), checkAry[i]) && !strings.Contains(checkAry[i], "{{Content-Length") {
							err = errors.New("漏洞不存在！")
							res.IsCheck = false
							return res, err
						}
					}
					res.Body = string(body)
					res.IsCheck = true
					return res, err
				} else {
					checkAry := strings.Split(rule.Check, "||")
					for i := 0; i < len(checkAry); i++ {
						if strings.Contains(string(body), checkAry[i]) || (checkAry[i] == "NULL" && string(body) == "") {
							res.Body = string(body)
							res.IsCheck = true
							return res, err
						}
					}
					err = errors.New("漏洞不存在！")
				}

			} else {
				err = errors.New("漏洞不存在！")
			}
		} else {
			err = errors.New("请求被拒绝！")
		}
	} else {
		err = errors.New("请求出错！")
	}

	return res, err

}

func httpRequestToExploitation(rule publicCode.Rule, url string) (publicCode.RespondData, error) {

	//定义返回包数据结构体
	var res publicCode.RespondData

	//定义http响应
	var resp *http.Response
	var requ *http.Request
	var err error

	URL := url + rule.Path

	rule.Body = strings.ReplaceAll(rule.Body, "\n", "\r\n")
	rule.Body = strings.ReplaceAll(rule.Body, "\r\r\n", "\r\n")

	if strings.Contains(rule.Body, "{{hex_decode(\"") {
		regex := regexp.MustCompile(`{{hex_decode\("(.+?)"\)}}`)
		regexAry := regex.FindAllString(rule.Body, -1)
		ruleBody := strings.ReplaceAll(regexAry[0], "{{hex_decode(\"", "")
		ruleBody = strings.ReplaceAll(ruleBody, "\")}}", "")
		ruleBodyByte, err := hex.DecodeString(ruleBody)
		if err != nil {
			fmt.Println("解码失败：", err)
			return res, err
		}
		requ, err = http.NewRequest(rule.Method, URL, bytes.NewReader(ruleBodyByte))
		if err != nil {
			fmt.Printf("get failed, err:%v\n", err)
			return res, err
		}
	} else {
		requ, err = http.NewRequest(rule.Method, URL, strings.NewReader(rule.Body))
		if err != nil {
			fmt.Printf("get failed, err:%v\n", err)
			return res, err
		}
	}

	//设置请求头
	for head := range rule.Headers {
		requ.Header.Set(head, rule.Headers[head])
	}

	resp, err = ClientNoRedirect.Do(requ)
	if err != nil {
		// 检查是否是超时错误
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			//fmt.Sprintf("请求超时！\n")
			err = errors.New("请求超时！")
			return res, err
		}
	}

	if resp != nil {

		defer resp.Body.Close()

		m := make(map[string]string)
		for k, v := range resp.Header {
			m[k] = v[0]
		}
		res.Headers = m
		res.Status = resp.Status

		if strings.Contains(resp.Status, rule.Status) {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("fetch: %v\n", err)
				return res, err
			}
			res.Body = string(body)
			res.IsCheck = true
			return res, err

		} else {
			err = errors.New("漏洞不存在！")
		}

	} else {
		err = errors.New("请求出错！")
	}

	return res, err

}

func FingerprintEasyHttp(url string) publicCode.FingerprintScanRespond {

	var wg sync.WaitGroup
	//定义http响应
	var resp1 *http.Response
	var resp2 *http.Response

	var fsRespond publicCode.FingerprintScanRespond

	wg.Add(1)
	go func() {
		defer wg.Done()
		requ1, err := http.NewRequest("GET", url, strings.NewReader(""))
		if err != nil {
			fmt.Printf("get failed, err:%v\n", err)
			return
		}

		requ1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
		requ1.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		resp1, err = Client.Do(requ1)
		if err != nil {
			// 检查是否是超时错误
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				//fmt.Sprintf("请求超时！\n")
				err = errors.New("请求超时！")
				return
			}
		}

		if resp1 != nil {

			defer resp1.Body.Close()
			m := make(map[string]string)
			for k, v := range resp1.Header {
				m[k] = v[0]
			}
			fsRespond.Headers = m
			fsRespond.StatusCode = resp1.Status
			if err != nil {
				fmt.Printf("FingerprintEasyHttp fetch: %v\n", err)
				return
			}
			var body []byte
			body, err = ioutil.ReadAll(resp1.Body)
			if err != nil {
				fmt.Printf("fetch: %v\n", err)
				return
			}
			//if strings.Contains(fsRespond.Headers["Content-Type"], "GBK") || strings.Contains(fsRespond.Headers["Content-Type"], "gbk") {
			//	//enc := mahonia.NewEncoder("gbk")
			//	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(body), simplifiedchinese.GBK.NewDecoder()))
			//	fsRespond.Body = string(data)
			//} else if strings.Contains(fsRespond.Headers["Content-Type"], "GB18030") || strings.Contains(fsRespond.Headers["Content-Type"], "gb18030") {
			//	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(body), simplifiedchinese.GB18030.NewDecoder()))
			//	fsRespond.Body = string(data)
			//} else {
			//	fsRespond.Body = string(body)
			//}

			// 获取Content-Type头部字段
			contentType := resp1.Header.Get("Content-Type")

			// 判断编码并解码
			encoding, _, _ := charset.DetermineEncoding(body, contentType)
			if encoding != nil {
				data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(body), encoding.NewDecoder()))
				if err != nil {
					fmt.Println("解码失败：", err)
				}
				fsRespond.Body = string(data)
			} else {
				fsRespond.Body = string(body)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		URL := url + "/favicon.ico"
		requ2, err := http.NewRequest("GET", URL, strings.NewReader(""))
		if err != nil {
			fmt.Printf("get failed, err:%v\n", err)
			return
		}

		requ2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
		requ2.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		resp2, err = Client.Do(requ2)
		if err != nil {
			// 检查是否是超时错误
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				//fmt.Sprintf("请求超时！\n")
				err = errors.New("请求超时！")
				return
			}
		}

		if resp2 != nil {

			defer resp2.Body.Close()
			var body []byte
			body, err = ioutil.ReadAll(resp2.Body)
			if err != nil {
				fmt.Printf("fetch: %v\n", err)
				return
			}

			hash := md5.Sum(body)
			fsRespond.FaviconHash = hex.EncodeToString(hash[:])
		}
	}()

	wg.Wait()
	return fsRespond
}

func FingerprintEasyHttp2(url string, data publicCode.FingerprintScanData) publicCode.FingerprintScanRespond {

	//定义http响应
	var resp *http.Response

	var fsRespond publicCode.FingerprintScanRespond

	URL := url + data.Path

	requ, err := http.NewRequest(strings.ToUpper(data.RequestMethod), URL, strings.NewReader(data.RequestData))
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
	}

	//设置请求头
	for head := range data.RequestHeaders {
		requ.Header.Set(head, data.RequestHeaders[head])
	}

	requ.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	requ.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	resp, err = Client.Do(requ)
	if err != nil {
		// 检查是否是超时错误
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			//fmt.Sprintf("请求超时！\n")
			err = errors.New("请求超时！")
		}
	}

	if resp != nil {

		defer resp.Body.Close()
		m := make(map[string]string)
		for k, v := range resp.Header {
			m[k] = v[0]
		}
		fsRespond.Headers = m
		fsRespond.StatusCode = resp.Status
		if err != nil {
			fmt.Printf("FingerprintEasyHttp fetch: %v\n", err)
		}
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("fetch: %v\n", err)
		}
		// 获取Content-Type头部字段
		contentType := resp.Header.Get("Content-Type")

		// 判断编码并解码
		encoding, _, _ := charset.DetermineEncoding(body, contentType)
		if encoding != nil {
			bodyData, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(body), encoding.NewDecoder()))
			if err != nil {
				fmt.Println("解码失败：", err)
			}
			fsRespond.Body = string(bodyData)
		} else {
			fsRespond.Body = string(body)
		}
	}

	return fsRespond
}
