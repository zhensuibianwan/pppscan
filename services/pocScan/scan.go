package pocScan

import (
	"changeme/services/mysqldb"
	"changeme/services/publicCode"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func PocEzScan(pocs []publicCode.Poc, url string, PocChanStart chan string, PocChanEnd chan int) {

	var msg string
	var err error

	if ClientNoRedirect == nil {
		CreatPocScanClient()
	}
	msg = fmt.Sprintf("对%s的扫描结果如下：\n", url)
	PocChanStart <- msg
	for _, poc := range pocs {
		Rules := make([]publicCode.Rule, len(poc.Request))
		Rules = dataHandle(poc)
		if Rules[0].Method == "" || Rules[0].Path == "" {
			msg = fmt.Sprintf("[-] %s：该poc存在问题！\n", poc.Name)
			PocChanStart <- msg
		} else {
			NeedData := poc.NeedData
			need := make(map[string]string, len(NeedData))
			for _, needData := range NeedData {
				label := "~" + needData.Label + "~"
				need[label] = needData.Value
			}
			Respond := make([]publicCode.RespondData, poc.OptionValue)

			for i := 0; i < poc.OptionValue; i++ {
				Rules[i] = NeedDataHandle(need, Rules[i], url, Respond)
				Respond[i], err = httpRequest(Rules[i], url)
				Rules[i] = printNeedHandle(need, Rules[i], url, Respond)
				if err != nil {
					errStr := fmt.Sprintf("%s", err)
					if errStr == "漏洞不存在！" {
						if i != 0 {
							if Respond[i-1].IsCheck == true && Rules[i-1].Print != "" {
								msg = fmt.Sprintf("[*] %s：%s\n", poc.Name, Rules[i-1].Print)
								PocChanStart <- msg
							} else if !strings.Contains(Respond[i].Status, Rules[i].Status) || i == poc.OptionValue-1 {
								msg = fmt.Sprintf("[-] %s：%s\n", poc.Name, errStr)
								PocChanStart <- msg
								break
							}
						} else {
							if !strings.Contains(Respond[i].Status, Rules[i].Status) || i == poc.OptionValue-1 {
								msg = fmt.Sprintf("[-] %s：%s\n", poc.Name, errStr)
								PocChanStart <- msg
								break
							}
						}
					} else {
						msg = fmt.Sprintf("[-] %s：%s\n", poc.Name, errStr)
						PocChanStart <- msg
						break
					}

				}
				if i == poc.OptionValue-1 && Respond[i].IsCheck == true {
					msg = fmt.Sprintf("[+] %s：%s\n", poc.Name, Rules[i].Print)
					PocChanStart <- msg
				}

			}
		}

	}
	PocChanEnd <- 1

}

type pocList struct {
	pocs   []publicCode.Poc
	listId int
}

type urlScan struct {
	url     string
	isScan  []int
	scanNum int
}

var ScanResultChan chan publicCode.PocScanResult
var ScanEndChan chan int
var ScanPauseChan chan int
var ScanProgressChan chan int
var ScanPauseChan2 chan int

func BatchScan1(urls []string, pocs []publicCode.Poc) {

	var thread int
	var wg sync.WaitGroup
	var lock sync.Mutex

	if ClientNoRedirect == nil {
		CreatPocScanClient()
	}

	if len(pocs)%publicCode.ThreadNum != 0 {
		thread = len(pocs)/publicCode.ThreadNum + 1
	} else {
		thread = len(pocs) / publicCode.ThreadNum
	}

	urlList := make([]urlScan, len(urls))
	for u, url := range urls {
		wg.Add(1)
		go func(u int, url string) {
			defer wg.Done()
			url = strings.ReplaceAll(url, "\r", "")
			url = strings.TrimRight(url, "/")
			if !strings.Contains(url, "http") {
				url = "http://" + url
			}
			urlList[u].url = url
			urlList[u].isScan = make([]int, thread)
			urlList[u].scanNum = 0
		}(u, url)
	}

	pocsList := make([]pocList, thread)
	for i := 0; i < thread; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i < thread-1 {
				pocsList[i].pocs = pocs[(i)*publicCode.ThreadNum : (i+1)*publicCode.ThreadNum]
				pocsList[i].listId = i + 1

			} else {
				pocsList[i].pocs = pocs[i*publicCode.ThreadNum:]
				pocsList[i].listId = i + 1
			}
		}(i)
	}
	wg.Wait()

	urlChan := make(chan urlScan, thread)
	pocListChan := make(chan pocList, thread)
	end := make(chan int)

	go func() {
		for i := 0; i < thread; i++ {
			urlChan <- urlList[i]
		}
	}()

	//go func() {
	//	for _, UL := range urlList {
	//		urlChan <- UL
	//	}
	//}()

	go func() {
		for _, list := range pocsList {
			pocListChan <- list
		}
	}()

	//t := time.Now()
	//msg := fmt.Sprintf("[*] 开始扫描：%s\n", t.Format("2006-01-02 15:04:05"))
	//ScanResultChan <- msg
	num := 0
	urlNum := thread
	existNum := 0

	for true {
		select {
		case <-end:
			num++
			if num == len(urls) {
				wg.Wait()
				//endTime := time.Since(t)
				//msg := fmt.Sprintf("[*] 共扫描%d次！共发现%d个漏洞！耗时：%s\n", num*len(pocs), existNum, endTime)
				//ScanResultChan <- msg
				ScanEndChan <- 1
				return
			}
			if urlNum < len(urls) {
				urlChan <- urlList[urlNum]
				urlNum++
			}

		case list := <-pocListChan:
			go func() {
				for true {
					ScanUrl := <-urlChan
					if ScanUrl.isScan[list.listId-1] != list.listId {
						ScanUrl.isScan[list.listId-1] = list.listId
						Ruless := make([][]publicCode.Rule, len(list.pocs))
						for p, poc := range list.pocs {
							Ruless[p] = make([]publicCode.Rule, len(poc.Request))
							wg.Add(1)
							go func(p int, poc publicCode.Poc) {
								defer wg.Done()
								Ruless[p] = dataHandle(poc)
								if Ruless[p][0].Method == "" || Ruless[p][0].Path == "" {
									Ruless = append(Ruless[:p], Ruless[p+1:]...)
								} else {
									NeedData := poc.NeedData
									need := make(map[string]string, len(NeedData))
									for _, needData := range NeedData {
										label := "~" + needData.Label + "~"
										need[label] = needData.Value
									}
									Respond := make([]publicCode.RespondData, poc.OptionValue)
									var err error
									for i := 0; i < poc.OptionValue; i++ {
										Ruless[p][i] = NeedDataHandle(need, Ruless[p][i], ScanUrl.url, Respond)
										Respond[i], err = httpRequest(Ruless[p][i], ScanUrl.url)
										Ruless[p][i] = printNeedHandle(need, Ruless[p][i], ScanUrl.url, Respond)
										if err != nil {
											errStr := fmt.Sprintf("%s", err)
											if errStr == "漏洞不存在！" {
												if i != 0 {
													if Respond[i-1].IsCheck == true {
														//msg := fmt.Sprintf("[*] %s：%s漏洞可能存在！\n", ScanUrl.url, poc.Name)
														var msg publicCode.PocScanResult
														msg.URL = ScanUrl.url
														msg.PocName = poc.Name
														msg.Print = "漏洞可能存在！"
														ScanResultChan <- msg
														lock.Lock()
														existNum++
														lock.Unlock()
														ScanProgressChan <- 1
													} else if !strings.Contains(Respond[i].Status, Ruless[p][i].Status) {
														ScanProgressChan <- 1
														break
													}

												} else {
													ScanProgressChan <- 1
													break
												}
											} else {
												ScanProgressChan <- 1
												break
											}

										}
										if i == poc.OptionValue-1 && Respond[i].IsCheck == true {
											//msg := fmt.Sprintf("[+] %s：%s\n", ScanUrl.url, poc.Name)
											var msg publicCode.PocScanResult
											msg.URL = ScanUrl.url
											msg.PocName = poc.Name
											msg.Print = "漏洞存在！"
											ScanResultChan <- msg
											lock.Lock()
											existNum++
											lock.Unlock()
											ScanProgressChan <- 1
										}
									}
								}
							}(p, poc)
						}

						ScanUrl.scanNum++

					} // if-1

					if ScanUrl.scanNum != thread {
						urlChan <- ScanUrl
					} else {
						end <- 1
					}

				} //for true-2
			}()

		} //select
	} //for true-1

}

func BatchScan2(urls []string, pocs []publicCode.Poc) {

	var wg sync.WaitGroup
	var lock sync.Mutex
	num := 0
	existNum := 0

	if ClientNoRedirect == nil {
		CreatPocScanClient()
	}
	worker := make(chan int, publicCode.ThreadNum)
	//t := time.Now()
	//msg := fmt.Sprintf("[*] 开始扫描：%s\n", t.Format("2006-01-02 15:04:05"))
	//ScanResultChan <- msg
	for _, url := range urls {
		url = strings.ReplaceAll(url, "\r", "")
		url = strings.TrimSpace(url)
		url = strings.TrimRight(url, "/")
		if !strings.Contains(url, "http") {
			url = "http://" + url
		}
		for _, poc := range pocs {
			wg.Add(1)
			go func(url string, poc publicCode.Poc) {
				defer wg.Done()
				worker <- 1
				Rules := make([]publicCode.Rule, len(poc.Request))
				Rules = dataHandle(poc)
				if Rules[0].Method == "" || Rules[0].Path == "" {
					//msg = fmt.Sprintf("[-] %s：该poc存在问题！\n", poc.Name)
					//PocChanStart <- msg

				} else {
					NeedData := poc.NeedData
					need := make(map[string]string, len(NeedData))
					for _, needData := range NeedData {
						label := "~" + needData.Label + "~"
						need[label] = needData.Value
					}
					Respond := make([]publicCode.RespondData, poc.OptionValue)
					var err error
					for i := 0; i < poc.OptionValue; i++ {
						Rules[i] = NeedDataHandle(need, Rules[i], url, Respond)
						Respond[i], err = httpRequest(Rules[i], url)
						Rules[i] = printNeedHandle(need, Rules[i], url, Respond)
						if err != nil {
							errStr := fmt.Sprintf("%s", err)
							if errStr == "漏洞不存在！" {
								if i != 0 {
									if Respond[i-1].IsCheck == true && Rules[i-1].Print != "" {
										var msg publicCode.PocScanResult
										msg.URL = url
										msg.PocName = poc.Name
										//msg.Print = "漏洞可能存在！"
										msg.Print = Rules[i-1].Print
										ScanResultChan <- msg
										ScanProgressChan <- 1
										lock.Lock()
										num++
										existNum++
										lock.Unlock()
									} else if !strings.Contains(Respond[i].Status, Rules[i].Status) || i == poc.OptionValue-1 {
										//msg = fmt.Sprintf("[-] %s：%s\n", poc.Name, errStr)
										//PocChanStart <- msg
										ScanProgressChan <- 1
										lock.Lock()
										num++
										lock.Unlock()
										break
									}

								} else {
									if !strings.Contains(Respond[i].Status, Rules[i].Status) || i == poc.OptionValue-1 {
										//msg = fmt.Sprintf("[-] %s：%s\n", poc.Name, errStr)
										//PocChanStart <- msg
										//fmt.Printf("[-] %s漏洞不存在！", url)
										ScanProgressChan <- 1
										lock.Lock()
										num++
										lock.Unlock()
										break
									}

								}
							} else {
								//msg = fmt.Sprintf("[-] %s：%s\n", poc.Name, errStr)
								//PocChanStart <- msg
								//fmt.Printf("[-] %s漏洞不存在！", url)
								ScanProgressChan <- 1
								lock.Lock()
								num++
								lock.Unlock()
								break
							}

						}
						if i == poc.OptionValue-1 && Respond[i].IsCheck == true {
							//msg = fmt.Sprintf("[+] %s：%s\n", poc.Name, Rules[i].Print)
							//PocChanStart <- msg
							var msg publicCode.PocScanResult
							msg.URL = url
							msg.PocName = poc.Name
							msg.Print = Rules[i].Print
							//msg := fmt.Sprintf("[+] %s：%s%s\n", url, poc.Name, Rules[i].Print)
							ScanResultChan <- msg
							ScanProgressChan <- 1
							lock.Lock()
							num++
							existNum++
							lock.Unlock()
						}

					}
				}
				<-worker
			}(url, poc)

		}
	}

	wg.Wait()
	//endTime := time.Since(t)
	//msg = fmt.Sprintf("[*] 共扫描%d次！共发现%d个漏洞！耗时：%s\n", num, existNum, endTime)
	//ScanResultChan <- msg
	ScanEndChan <- 1
}

func ExploitationScan(poc publicCode.Poc, url string, PocChanStart chan string, PocChanEnd chan int) {

	var msg string

	if ClientNoRedirect == nil {
		CreatPocScanClient()
	}

	Rules := make([]publicCode.Rule, len(poc.Request))
	Rules = dataHandle(poc)
	if Rules[0].Method == "" || Rules[0].Path == "" {
		msg = fmt.Sprintf("[-] %s：该poc存在问题！\n", poc.Name)
		PocChanStart <- msg
	} else {
		NeedData := poc.NeedData
		need := make(map[string]string, len(NeedData))
		for _, needData := range NeedData {
			label := "~" + needData.Label + "~"
			need[label] = needData.Value
		}
		Respond := make([]publicCode.RespondData, poc.OptionValue)

		for i := 0; i < poc.OptionValue; i++ {
			Rules[i] = NeedDataHandle(need, Rules[i], url, Respond)
			Respond[i], _ = httpRequestToExploitation(Rules[i], url)
			Rules[i] = printNeedHandle(need, Rules[i], url, Respond)
			if i == poc.OptionValue-1 && Respond[i].IsCheck == true {
				msg = strings.ReplaceAll(Rules[i].Print, "漏洞存在！", "")
				msg = strings.ReplaceAll(msg, "漏洞存在!", "")
				PocChanStart <- msg
			}

		}
	}

	PocChanEnd <- 1

}

var FSResultChan chan publicCode.FingerprintScanResult
var fsWorker chan int
var FSScanEndChan chan int
var FSScanPauseChan chan int
var FSScanProgressChan chan int
var FSScanPauseChan2 chan int
var fsPocScanWorker chan int

func FingerprintScan(urls []string, linkPoc bool, isEasyScan bool) {
	var wg sync.WaitGroup

	if Client == nil {
		CreatPocScanClient()
	}
	fsWorker = make(chan int, publicCode.ThreadNum)
	fingerprints := mysqldb.SearchFingerprintByName("")

	//var StartOREnd publicCode.FingerprintScanResult
	//t := time.Now()
	//StartOREnd.URL = fmt.Sprintf("[*] 开始扫描：%s\n", t.Format("2006-01-02 15:04:05"))
	//FSResultChan <- StartOREnd
	for i, url := range urls {
		wg.Add(1)
		var result publicCode.FingerprintScanResult
		result.URL = url
		go func(i int, url string, result publicCode.FingerprintScanResult) {
			defer wg.Done()
			fsWorker <- 1
			fSRespond := FingerprintEasyHttp(url)
			if fSRespond.Body == "" {
				FSScanProgressChan <- 1
				<-fsWorker
				return
			}
			if strings.Contains(fSRespond.Body, "<title>") && strings.Contains(fSRespond.Body, "</title>") {
				regex := regexp.MustCompile(`<title>.*</title>`)
				regexAry := regex.FindAllString(fSRespond.Body, -1)
				if len(regexAry) != 0 {
					result.Title = strings.ReplaceAll(regexAry[0], "<title>", "")
					result.Title = strings.ReplaceAll(result.Title, "</title>", "")
				}
			} else if strings.Contains(fSRespond.Body, "<TITLE>") && strings.Contains(fSRespond.Body, "</TITLE>") {
				regex := regexp.MustCompile(`<TITLE>.*</TITLE>`)
				regexAry := regex.FindAllString(fSRespond.Body, -1)
				if len(regexAry) != 0 {
					result.Title = strings.ReplaceAll(regexAry[0], "<TITLE>", "")
					result.Title = strings.ReplaceAll(result.Title, "</TITLE>", "")
				}
			}

			result = fingerprintCompare(fingerprints, url, fSRespond, result, linkPoc)
			if isEasyScan {
				FSResultChan <- result
			} else if !isEasyScan && len(result.Fingerprint) > 0 {
				FSResultChan <- result
			}
			FSScanProgressChan <- 1
			<-fsWorker
		}(i, url, result)
	}
	wg.Wait()
	//endTime := time.Since(t)
	//StartOREnd.URL = fmt.Sprintf("[*] 扫描结束！耗时：%s\n", endTime)
	//FSResultChan <- StartOREnd
	FSScanEndChan <- 1

}

func fingerprintCompare(fingerprints []publicCode.Fingerprint, url string, fSRespond publicCode.FingerprintScanRespond, result publicCode.FingerprintScanResult, linkPoc bool) publicCode.FingerprintScanResult {

	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, fingerprint := range fingerprints {
		wg.Add(1)
		go func(fingerprint publicCode.Fingerprint) {
			defer wg.Done()
		fingerprintScan:
			for _, fingerprintScan := range fingerprint.FingerprintScan {
				if fingerprintScan.Path != "/" || fingerprintScan.RequestData != "" || len(fingerprintScan.RequestHeaders) != 0 {
					fs := FingerprintEasyHttp2(url, fingerprintScan)
					if len(fingerprintScan.Headers) != 0 {
						for key, vue := range fingerprintScan.Headers {
							if !strings.Contains(fs.Headers[key], vue) {
								continue fingerprintScan
							}
						}
					}
					if len(fingerprintScan.Keyword) != 0 {
						for _, vue := range fingerprintScan.Keyword {
							if !strings.Contains(fs.Body, vue) {
								continue fingerprintScan
							}
						}
					}
					if fingerprintScan.StatusCode != 0 {
						if !strings.Contains(fs.StatusCode, strconv.Itoa(fingerprintScan.StatusCode)) {
							continue fingerprintScan
						}
					}
					lock.Lock()
					result.Fingerprint = append(result.Fingerprint, fingerprintScan.Name)
					lock.Unlock()
					if linkPoc && len(fingerprint.PocsInfo) != 0 {
						for _, pocInfo := range fingerprint.PocsInfo {
							poc := mysqldb.SearchPocByUUID(pocInfo.UUID)
							rs := linkPocScan(poc[0], url, result)
							lock.Lock()
							result = rs
							lock.Unlock()
						}
					}
					return

				} else {
					if len(fingerprintScan.Headers) != 0 {
						for key, vue := range fingerprintScan.Headers {
							if !strings.Contains(fSRespond.Headers[key], vue) {
								continue fingerprintScan
							}
						}
					}
					if len(fingerprintScan.Keyword) != 0 {
						for _, vue := range fingerprintScan.Keyword {
							if !strings.Contains(fSRespond.Body, vue) {
								continue fingerprintScan
							}
						}
					}
					if fingerprintScan.StatusCode != 0 {
						if !strings.Contains(fSRespond.StatusCode, strconv.Itoa(fingerprintScan.StatusCode)) {
							continue fingerprintScan
						}
					}
					if len(fingerprintScan.FaviconHash) != 0 {
						faviconHashCheck := 0
						for _, faviconHash := range fingerprintScan.FaviconHash {
							if fSRespond.FaviconHash == faviconHash {
								faviconHashCheck = 1
							}
						}
						if faviconHashCheck != 1 {
							continue fingerprintScan
						}
					}
					lock.Lock()
					result.Fingerprint = append(result.Fingerprint, fingerprintScan.Name)
					lock.Unlock()
					if linkPoc && len(fingerprint.PocsInfo) != 0 {
						for _, pocInfo := range fingerprint.PocsInfo {
							poc := mysqldb.SearchPocByUUID(pocInfo.UUID)
							if len(poc) > 0 {
								rs := linkPocScan(poc[0], url, result)
								lock.Lock()
								result = rs
								lock.Unlock()
							}
						}
					}
					return
				}
			}
		}(fingerprint)

	}
	wg.Wait()
	return result
}

func linkPocScan(poc publicCode.Poc, url string, result publicCode.FingerprintScanResult) publicCode.FingerprintScanResult {

	var err error
	Rules := make([]publicCode.Rule, len(poc.Request))
	Rules = dataHandle(poc)
	if Rules[0].Method == "" || Rules[0].Path == "" {

	} else {
		NeedData := poc.NeedData
		need := make(map[string]string, len(NeedData))
		for _, needData := range NeedData {
			label := "~" + needData.Label + "~"
			need[label] = needData.Value
		}
		Respond := make([]publicCode.RespondData, poc.OptionValue)

		for i := 0; i < poc.OptionValue; i++ {
			Rules[i] = NeedDataHandle(need, Rules[i], url, Respond)
			Respond[i], err = httpRequest(Rules[i], url)
			Rules[i] = printNeedHandle(need, Rules[i], url, Respond)
			if err != nil {
				errStr := fmt.Sprintf("%s", err)
				if errStr == "漏洞不存在！" {
					if i != 0 {
						if Respond[i-1].IsCheck == true && Rules[i-1].Print != "" {
							result.Vulnerability = append(result.Vulnerability, poc.Name)
						} else if !strings.Contains(Respond[i].Status, Rules[i].Status) || i == poc.OptionValue-1 {
							break
						}
					} else {
						if !strings.Contains(Respond[i].Status, Rules[i].Status) || i == poc.OptionValue-1 {
							break
						}
					}
				} else {
					break
				}

			}
			if i == poc.OptionValue-1 && Respond[i].IsCheck == true {
				result.Vulnerability = append(result.Vulnerability, poc.Name)
			}

		}
	}

	return result

}
