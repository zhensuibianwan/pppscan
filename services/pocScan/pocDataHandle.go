package pocScan

import (
	"changeme/services/publicCode"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func dataHandle(poc publicCode.Poc) []publicCode.Rule {

	var wg sync.WaitGroup
	Rule := make([]publicCode.Rule, len(poc.Request))

	for M, pocRequest := range poc.Request {
		wg.Add(1)
		go func(m int, pocReq publicCode.RequestData) {
			defer wg.Done()
			request := strings.Split(pocReq.PocString, "\n")
			if len(request) >= 3 {
				strHandle1 := strings.Split(request[0], " ")
				if len(strHandle1) >= 3 {
					Rule[m].Method = strHandle1[0]
					Rule[m].Path = strHandle1[1]
					var n int
					var headers []string
					if Rule[m].Method == "GET" {
						n = len(request)
						Rule[m].Body = ""
						request = strings.Split(pocReq.PocString, "\n")
						Rule[m].Body = ""
						headers = request[1:]
					} else {
						for i, t := range request {
							if t == "" && i != 0 {
								n = i
								break
							}
						}
						if n == 0 {
							request = strings.Split(pocReq.PocString, "\n")
							Rule[m].Body = ""
							headers = request[1:]
						} else {
							request = strings.SplitN(pocReq.PocString, "\n", n+2)
							Rule[m].Body = request[n+1]
							headers = request[1:n]
						}
					}

					headersData := make(map[string]string)
					for _, Header := range headers {
						if strings.Contains(Header, ":") {
							header := strings.SplitN(Header, ":", 2)
							if header[0] != "Host" && header[0] != "Connection" && header[0] != "Accept-Encoding" {
								headersData[header[0]] = header[1]
							}
						}
					}
					Rule[m].Headers = headersData
					Rule[m].Check = pocReq.Check
					Rule[m].Status = pocReq.Status
					Rule[m].Print = pocReq.Print
				}
			}

		}(M, pocRequest)
	}
	wg.Wait()
	return Rule
}

func NeedDataHandle(need map[string]string, rule publicCode.Rule, url string, res []publicCode.RespondData) publicCode.Rule {

	var wg sync.WaitGroup

	need["~request.url.0~"] = url
	regex := regexp.MustCompile(`~[^~]*~`)
	wg.Add(1)
	go func() {
		bodyNeeds := regex.FindAllString(rule.Body, -1)
		for _, bodyNeed := range bodyNeeds {
			if strings.Contains(bodyNeed, "url") || strings.Contains(bodyNeed, "input") {
				rule.Body = strings.ReplaceAll(rule.Body, bodyNeed, need[bodyNeed])
			} else if strings.Contains(bodyNeed, "body") {
				bodyNeedAry := strings.Split(bodyNeed, ".")

				resNum := strings.Trim(strings.ToLower(bodyNeedAry[0]), "~request")
				resN, err := strconv.Atoi(resNum)
				ResNum := resN - 1
				if err != nil {
					fmt.Printf("fetch: %v\n", err)
					os.Exit(1)
				}
				if res[ResNum].Status != "" {
					if need[bodyNeed] != "ALL" {
						need[bodyNeed] = strings.ReplaceAll(need[bodyNeed], "\r ", "\n")
						bodyPatternAry := strings.Split(need[bodyNeed], "~")
						if bodyPatternAry[0] == "" || bodyPatternAry[1] == "" {
							if bodyPatternAry[0] == "" {
								newStr := strings.TrimRight(res[ResNum].Body, bodyPatternAry[1])
								rule.Body = strings.ReplaceAll(rule.Body, bodyNeed, newStr)
							}
							if bodyPatternAry[1] == "" {
								newStr := strings.TrimLeft(res[ResNum].Body, bodyPatternAry[0])
								rule.Body = strings.ReplaceAll(rule.Body, bodyNeed, newStr)
							}
						} else {
							bodyPattern := regexp.MustCompile(fmt.Sprintf(`%s(.*?)%s`, regexp.QuoteMeta(bodyPatternAry[0]), regexp.QuoteMeta(bodyPatternAry[1])))
							if bodyPattern.MatchString(res[ResNum].Body) {
								newStr := bodyPattern.FindStringSubmatch(res[ResNum].Body)[1]
								rule.Body = strings.ReplaceAll(rule.Body, bodyNeed, newStr)
							}
						}

					} else {
						rule.Body = strings.ReplaceAll(rule.Body, bodyNeed, res[ResNum].Body)
					}
				}
			} else if strings.Contains(bodyNeed, "header") {
				headerNeedAry := strings.Split(bodyNeed, ".")

				resNum := strings.Trim(strings.ToLower(headerNeedAry[0]), "request")
				resN, err := strconv.Atoi(resNum)
				ResNum := resN - 1
				if err != nil {
					fmt.Printf("fetch: %v\n", err)
					os.Exit(1)
				}
				if res[ResNum].Status != "" {
					rule.Body = strings.ReplaceAll(rule.Body, bodyNeed, res[ResNum].Headers[need[bodyNeed]])
				}
			}

		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for headerKey, headerValue := range rule.Headers {
			if strings.Contains(headerValue, "~request") || strings.Contains(headerValue, "~input") {
				headerNeeds := regex.FindAllString(headerValue, -1)
				if strings.Contains(headerNeeds[0], "url") || strings.Contains(headerNeeds[0], "input") {
					rule.Headers[headerKey] = strings.ReplaceAll(headerValue, headerNeeds[0], need[headerNeeds[0]])
				} else if strings.Contains(headerNeeds[0], "body") {
					bodyNeedAry := strings.Split(headerNeeds[0], ".")

					resNum := strings.Trim(strings.ToLower(bodyNeedAry[0]), "~request")
					resN, err := strconv.Atoi(resNum)
					ResNum := resN - 1
					if err != nil {
						fmt.Printf("fetch: %v\n", err)
						os.Exit(1)
					}
					if res[ResNum].Status != "" {
						if need[headerNeeds[0]] != "ALL" {
							need[headerNeeds[0]] = strings.ReplaceAll(need[headerNeeds[0]], "\r ", "\n")
							bodyPatternAry := strings.Split(need[headerNeeds[0]], "~")
							if bodyPatternAry[0] == "" || bodyPatternAry[1] == "" {
								if bodyPatternAry[0] == "" {
									newStr := strings.TrimRight(res[ResNum].Body, bodyPatternAry[1])
									rule.Headers[headerKey] = strings.ReplaceAll(headerValue, headerNeeds[0], newStr)
								}
								if bodyPatternAry[1] == "" {
									newStr := strings.TrimLeft(res[ResNum].Body, bodyPatternAry[0])
									rule.Headers[headerKey] = strings.ReplaceAll(headerValue, headerNeeds[0], newStr)
								}
							} else {
								bodyPattern := regexp.MustCompile(fmt.Sprintf(`%s(.*?)%s`, regexp.QuoteMeta(bodyPatternAry[0]), regexp.QuoteMeta(bodyPatternAry[1])))
								if bodyPattern.MatchString(res[ResNum].Body) {
									newStr := bodyPattern.FindStringSubmatch(res[ResNum].Body)[1]
									rule.Headers[headerKey] = strings.ReplaceAll(headerValue, headerNeeds[0], newStr)
								}
							}

						} else {
							rule.Headers[headerKey] = strings.ReplaceAll(headerValue, headerNeeds[0], res[ResNum].Body)
						}
					}
				} else if strings.Contains(headerNeeds[0], "header") {
					headerNeedAry := strings.Split(headerNeeds[0], ".")

					resNum := strings.Trim(strings.ToLower(headerNeedAry[0]), "~request")
					resN, err := strconv.Atoi(resNum)
					ResNum := resN - 1
					if err != nil {
						fmt.Printf("fetch: %v\n", err)
						os.Exit(1)
					}
					if res[ResNum].Status != "" {
						rule.Headers[headerKey] = strings.ReplaceAll(headerValue, headerNeeds[0], res[ResNum].Headers[need[headerNeeds[0]]])
					}
				}

			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		pathNeeds := regex.FindAllString(rule.Path, -1)
		for _, pathNeed := range pathNeeds {
			if strings.Contains(pathNeed, "url") || strings.Contains(pathNeed, "input") {
				rule.Path = strings.ReplaceAll(rule.Path, pathNeed, need[pathNeed])
			} else if strings.Contains(pathNeed, "body") {
				bodyNeedAry := strings.Split(pathNeed, ".")
				resNum := strings.Trim(strings.ToLower(bodyNeedAry[0]), "~request")
				resN, err := strconv.Atoi(resNum)
				ResNum := resN - 1
				if err != nil {
					fmt.Printf("fetch: %v\n", err)
					os.Exit(1)
				}
				if res[ResNum].Status != "" {
					if need[pathNeed] != "ALL" {
						need[pathNeed] = strings.ReplaceAll(need[pathNeed], "\r ", "\n")
						bodyPatternAry := strings.Split(need[pathNeed], "~")
						if bodyPatternAry[0] == "" || bodyPatternAry[1] == "" {
							if bodyPatternAry[0] == "" {
								newStr := strings.TrimRight(res[ResNum].Body, bodyPatternAry[1])
								rule.Path = strings.ReplaceAll(rule.Path, pathNeed, newStr)
							}
							if bodyPatternAry[1] == "" {
								newStr := strings.TrimLeft(res[ResNum].Body, bodyPatternAry[0])
								rule.Path = strings.ReplaceAll(rule.Path, pathNeed, newStr)
							}
						} else {
							bodyPattern := regexp.MustCompile(fmt.Sprintf(`%s(.*?)%s`, regexp.QuoteMeta(bodyPatternAry[0]), regexp.QuoteMeta(bodyPatternAry[1])))
							if bodyPattern.MatchString(res[ResNum].Body) {
								newStr := bodyPattern.FindStringSubmatch(res[ResNum].Body)[1]
								rule.Path = strings.ReplaceAll(rule.Path, pathNeed, newStr)
							}
						}
					} else {
						rule.Path = strings.ReplaceAll(rule.Path, pathNeed, res[ResNum].Body)
					}
				}
			} else if strings.Contains(pathNeed, "header") {
				headerNeedAry := strings.Split(pathNeed, ".")

				resNum := strings.Trim(strings.ToLower(headerNeedAry[0]), "~request")
				resN, err := strconv.Atoi(resNum)
				ResNum := resN - 1
				if err != nil {
					fmt.Printf("fetch: %v\n", err)
					os.Exit(1)
				}
				if res[ResNum].Status != "" {
					rule.Path = strings.ReplaceAll(rule.Path, pathNeed, res[ResNum].Headers[need[pathNeed]])
				}
			}

		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		checkNeeds := regex.FindAllString(rule.Check, -1)
		for _, checkNeed := range checkNeeds {
			rule.Check = strings.ReplaceAll(rule.Check, checkNeed, need[checkNeed])
		}
		wg.Done()
	}()

	wg.Wait()

	//fmt.Println(rule)
	return rule
}

func printNeedHandle(need map[string]string, rule publicCode.Rule, url string, res []publicCode.RespondData) publicCode.Rule {

	need["~request.url.0~"] = url
	regex := regexp.MustCompile(`~[^~]*~`)
	printNeeds := regex.FindAllString(rule.Print, -1)
	for _, printNeed := range printNeeds {
		if strings.Contains(printNeed, "url") || strings.Contains(printNeed, "input") {
			rule.Print = strings.ReplaceAll(rule.Print, printNeed, need[printNeed])
		} else if strings.Contains(printNeed, "body") {
			bodyNeedAry := strings.Split(printNeed, ".")
			resNum := strings.Trim(strings.ToLower(bodyNeedAry[0]), "~request")
			resN, err := strconv.Atoi(resNum)
			ResNum := resN - 1
			if err != nil {
				fmt.Printf("fetch: %v\n", err)
				os.Exit(1)
			}
			if res[ResNum].Status != "" {
				if need[printNeed] != "ALL" {
					need[printNeed] = strings.ReplaceAll(need[printNeed], "\r ", "\n")
					bodyPatternAry := strings.Split(need[printNeed], "~")
					if bodyPatternAry[0] == "" || bodyPatternAry[1] == "" {
						if bodyPatternAry[0] == "" {
							newStr := strings.TrimRight(res[ResNum].Body, bodyPatternAry[1])
							rule.Print = strings.ReplaceAll(rule.Print, printNeed, newStr)
							//bodyPattern := regexp.MustCompile(fmt.Sprintf(`(.*?)%s`, regexp.QuoteMeta(bodyPatternAry[1])))
							//if bodyPattern.MatchString(res[ResNum].Body) {
							//	newStr := bodyPattern.FindStringSubmatch(res[ResNum].Body)[1]
							//	rule.Print = strings.ReplaceAll(rule.Print, printNeed, newStr)
							//}
						}
						if bodyPatternAry[1] == "" {
							newStr := strings.TrimLeft(res[ResNum].Body, bodyPatternAry[0])
							rule.Print = strings.ReplaceAll(rule.Print, printNeed, newStr)
							//bodyPattern := regexp.MustCompile(fmt.Sprintf(`%s(.*?)`, regexp.QuoteMeta(bodyPatternAry[0])))
							//if bodyPattern.MatchString(res[ResNum].Body) {
							//	newStr := bodyPattern.FindStringSubmatch(res[ResNum].Body)[1]
							//	rule.Print = strings.ReplaceAll(rule.Print, printNeed, newStr)
							//}
						}
					} else {
						bodyPattern := regexp.MustCompile(fmt.Sprintf(`%s(.*?)%s`, regexp.QuoteMeta(bodyPatternAry[0]), regexp.QuoteMeta(bodyPatternAry[1])))
						if bodyPattern.MatchString(res[ResNum].Body) {
							newStr := bodyPattern.FindStringSubmatch(res[ResNum].Body)[1]
							rule.Print = strings.ReplaceAll(rule.Print, printNeed, newStr)
						}
					}
				} else {
					rule.Print = strings.ReplaceAll(rule.Print, printNeed, res[ResNum].Body)
				}
			}
		} else if strings.Contains(printNeed, "header") {
			headerNeedAry := strings.Split(printNeed, ".")

			resNum := strings.Trim(strings.ToLower(headerNeedAry[0]), "~request")
			resN, err := strconv.Atoi(resNum)
			ResNum := resN - 1
			if err != nil {
				fmt.Printf("fetch: %v\n", err)
				os.Exit(1)
			}
			if res[ResNum].Status != "" {
				rule.Print = strings.ReplaceAll(rule.Print, printNeed, res[ResNum].Headers[need[printNeed]])
			}
		}
	}

	return rule
}
