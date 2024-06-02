package main

import (
	"changeme/services/config"
	"changeme/services/fingerprintManage"
	"changeme/services/mysqldb"
	"changeme/services/pocManage"
	"changeme/services/pocScan"
	"changeme/services/publicCode"
	"changeme/services/tools"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) PocList(key string) []publicCode.Poc {
	var pocs []publicCode.Poc
	//if publicCode.OnlyLocalLoading {
	//	pocs = publicCode.PocFileList(key)
	//} else {
	pocs = pocManage.PocObjects(key)
	//}
	return pocs
}

func (a *App) PocList2(key string) []publicCode.Poc {
	pocs := publicCode.PocFileList(key)
	return pocs
}

func (a *App) LocalList(key string) []publicCode.Poc {
	var pocs []publicCode.Poc
	if publicCode.LocalLoading {
		pocs = publicCode.PocFileList(key)
	}
	return pocs
}

func (a *App) AddPoc(poc publicCode.Poc) []string {
	status, msg := pocManage.SavePoc(poc)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) DelPoc(pocs []publicCode.Poc) []string {
	status, msg := pocManage.RemovePoc(pocs)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) DelPocFile(pocs []publicCode.Poc) []string {
	status, msg := pocManage.RemovePocFile(pocs)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) PocScan(pocs []publicCode.Poc, url string) {

	PocChanStart := make(chan string)
	PocChanEnd := make(chan int)
	url = strings.TrimRight(url, "/")
	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	go func() {
		pocScan.PocEzScan(pocs, url, PocChanStart, PocChanEnd)
	}()

	go func() {
		for true {
			go func() {
				end := <-PocChanEnd
				if end == 1 {
					return
				}
			}()
			msg := <-PocChanStart
			runtime.EventsEmit(a.ctx, "PocScan", msg)
		}
	}()

}

func (a *App) PocExploitation(poc publicCode.Poc, url string) {

	PocChanStart := make(chan string)
	PocChanEnd := make(chan int)
	url = strings.TrimRight(url, "/")
	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	go func() {
		pocScan.ExploitationScan(poc, url, PocChanStart, PocChanEnd)
	}()

	go func() {
		for true {
			go func() {
				end := <-PocChanEnd
				if end == 1 {
					return
				}
			}()
			msg := <-PocChanStart
			runtime.EventsEmit(a.ctx, "ExploitationScan", msg)
		}
	}()

}

func (a *App) SaveSetting(setting publicCode.AllSetting, dbChange int) []string {
	status, msg := config.SettingSave(setting, dbChange)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) ShowSetting() publicCode.AllSetting {
	return config.SettingShow()
}

func (a *App) DBInitialize() []string {

	status, msg := mysqldb.InitializeDB()
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result

}

func (a *App) BatchPocScan() publicCode.PocScanResult {

	var msg publicCode.PocScanResult
	runtime.EventsEmit(a.ctx, "BatchScan", msg)
	return msg
}

func (a *App) ScanBatch(URLS string, pocs []publicCode.Poc) {

	pocScan.ScanResultChan = make(chan publicCode.PocScanResult)
	pocScan.ScanEndChan = make(chan int)
	pocScan.ScanPauseChan = make(chan int)
	pocScan.ScanPauseChan2 = make(chan int, 2)
	pocScan.ScanProgressChan = make(chan int, publicCode.TimeOut-1)

	urls := strings.Split(URLS, "\n")
	go func() {
		//if len(pocs) < publicCode.ThreadNum {
		//	pocScan.BatchScan2(urls, pocs)
		//} else {
		//	pocScan.BatchScan1(urls, pocs)
		//}
		pocScan.BatchScan2(urls, pocs)
	}()

	num := 0
	go func() {
		for true {
			select {
			case end := <-pocScan.ScanEndChan:
				if end == 1 {
					time.Sleep(1 * time.Second)
					runtime.EventsEmit(a.ctx, "BatchScanProgress", 100)
					return
				}
				if end == 2 {
					//err := "[*] 扫描已取消！"
					//runtime.EventsEmit(a.ctx, "BatchScan", err)
					return
				}
			case msg := <-pocScan.ScanResultChan:

				runtime.EventsEmit(a.ctx, "BatchScan", msg)

			case p := <-pocScan.ScanProgressChan:
				num = num + p
				progress := num * 100 / (len(urls) * len(pocs))
				runtime.EventsEmit(a.ctx, "BatchScanProgress", progress)
			case <-pocScan.ScanPauseChan:
				<-pocScan.ScanPauseChan2
			}

		}
	}()

}

func (a *App) ScanPause() {

	pocScan.ScanPauseChan <- 1
}

func (a *App) ScanContinue() {

	pocScan.ScanPauseChan2 <- 1
}

func (a *App) ScanClose() {

	pocScan.ScanPauseChan2 <- 1
	pocScan.ScanEndChan <- 2
}

func (a *App) PocsDownload(pocs []publicCode.Poc) []string {

	status, msg := pocManage.DownLoadPocs(pocs)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) PocsUpload(pocs []publicCode.Poc) []string {

	status, msg := pocManage.UploadPocs(pocs)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) BackgroundSetting() string {
	return publicCode.BackgroundImage
	//status, msg := config.BackgroundImageCreat(image)
	//var result []string
	//result = append(result, strconv.Itoa(status), msg)
	//return result
}

func (a *App) BackgroundImageSave(image publicCode.Background) {

	config.BackgroundSave(image)
}

func (a *App) FingerprintList(key string) []publicCode.Fingerprint {

	var fingerprints []publicCode.Fingerprint
	//if publicCode.OnlyLocalLoading {
	//	fingerprints = publicCode.FingerprintFileList(key)
	//} else {
	fingerprints = fingerprintManage.FingerprintsList(key)
	//}
	return fingerprints

}

func (a *App) FingerprintFileList(key string) []publicCode.Fingerprint {

	fingerprints := fingerprintManage.UploadFingerprintsList(key)
	return fingerprints

}

func (a *App) FingerprintLocalList(key string) []publicCode.Fingerprint {

	var fingerprints []publicCode.Fingerprint
	if publicCode.LocalLoading {
		fingerprints = publicCode.FingerprintFileList(key)
	}
	return fingerprints

}

func (a *App) AddFingerprint(fingerprint publicCode.Fingerprint) []string {
	status, msg := fingerprintManage.SaveFingerprint(fingerprint)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) DelFingerprint(fingerprints []publicCode.Fingerprint) []string {
	status, msg := fingerprintManage.RemoveFingerprint(fingerprints)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) DelFingerprintFile(fingerprints []publicCode.Fingerprint) []string {
	status, msg := fingerprintManage.RemoveFingerprintFile(fingerprints)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) FingerprintsDownload(fingerprints []publicCode.Fingerprint) []string {

	status, msg := fingerprintManage.DownLoadFingerprintS(fingerprints)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) FingerprintsUpload(fingerprints []publicCode.Fingerprint) []string {

	status, msg := fingerprintManage.UploadFingerprints(fingerprints)
	var result []string
	result = append(result, strconv.Itoa(status), msg)
	return result
}

func (a *App) FingerprintScan() publicCode.FingerprintScanResult {

	var msg publicCode.FingerprintScanResult
	runtime.EventsEmit(a.ctx, "FingerprintScan", msg)
	return msg
}

func (a *App) FSScan(URLStr string, linkPoc bool, isEasyScan bool) {

	pocScan.FSResultChan = make(chan publicCode.FingerprintScanResult)
	pocScan.FSScanEndChan = make(chan int)
	pocScan.FSScanPauseChan = make(chan int)
	pocScan.FSScanPauseChan2 = make(chan int, 2)
	pocScan.FSScanProgressChan = make(chan int, publicCode.TimeOut-1)

	URLS := strings.Split(URLStr, "\n")
	var urls []string
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, url := range URLS {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			url = strings.ReplaceAll(url, "\r", "")
			url = strings.TrimSpace(url)
			url = strings.TrimRight(url, "/")
			if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
				url1 := "http://" + url
				url2 := "https://" + url
				lock.Lock()
				urls = append(urls, url1)
				urls = append(urls, url2)
				lock.Unlock()
			} else {
				lock.Lock()
				urls = append(urls, url)
				lock.Unlock()
			}
		}(url)

	}
	wg.Wait()
	go func() {
		pocScan.FingerprintScan(urls, linkPoc, isEasyScan)
	}()

	num := 0
	go func() {
		for true {
			select {
			case end := <-pocScan.FSScanEndChan:
				if end == 1 {
					time.Sleep(1 * time.Second)
					runtime.EventsEmit(a.ctx, "FingerprintScanProgress", 100)
					return
				}
				if end == 2 {
					//var err publicCode.FingerprintScanResult
					//err.URL = "[*] 扫描已取消！"
					//runtime.EventsEmit(a.ctx, "FingerprintScan", err)
					return
				}
			case msg := <-pocScan.FSResultChan:
				runtime.EventsEmit(a.ctx, "FingerprintScan", msg)

			case p := <-pocScan.FSScanProgressChan:
				num = num + p
				progress := num * 100 / len(urls)
				runtime.EventsEmit(a.ctx, "FingerprintScanProgress", progress)
			case <-pocScan.FSScanPauseChan:
				<-pocScan.FSScanPauseChan2
			}

		}
	}()

}

func (a *App) FSScanPause() {

	pocScan.FSScanPauseChan <- 1
}

func (a *App) FSScanContinue() {

	pocScan.FSScanPauseChan2 <- 1
}

func (a *App) FSScanClose() {

	pocScan.FSScanPauseChan2 <- 1
	pocScan.FSScanEndChan <- 2
}

func (a *App) FingerprintGetName(key string) []string {
	names := mysqldb.SearchFingerprintGetNameByName(key)
	return names
}

func (a *App) DefaultBrowserOpenUrl(url string) {
	tools.OpenUrlByDefaultBrowser(url)
}
