package surfer

import (
	"github.com/tebeka/selenium"
	"net/http"
	"net"
	"sync"
	"os"
	ch "github.com/tebeka/selenium/chrome"
	"fmt"
	"time"
	"strings"
	"io/ioutil"
	"github.com/henrylee2cn/pholcus/config"
	"github.com/henrylee2cn/pholcus/logs"
)

/** problems:
1、 运行久了，会报一下错误，说chromedriver no long running
	The process started from chrome location /opt/google/chrome/google-chrome is no longer running, so ChromeDriver is assuming that Chrome has crashed
	解决：如果你自定义路径下的exe, 需设置Chromedriver.exe路径参数，chrome默认启用的是local setttings里的chromedriver.exe
 */

type Sel struct {
	Service *selenium.Service
	caps selenium.Capabilities
	saveScript string
	port int
	host string
	initial bool
	maxReq  chan int      //最大请求数，过大会崩溃
	sessionGroup sync.WaitGroup
}
var ms sync.Mutex

func NewSel() Surfer{
	s := &Sel{
		initial: false,
		maxReq: make(chan int, 8),
	}

	return s
}

func (self *Sel)initIE(){
	if self.initial {
		return
	}
	ms.Lock()
	defer ms.Unlock()

	port,_ := pickUnusedPort()
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(false)
	service, err := selenium.NewIEDriverService(config.IEDRIVER,port,opts...)
	if err != nil {
		panic(err)
	}
	caps := selenium.Capabilities{"browserName":"internet explorer"}
	chromeCaps := ch.Capabilities{
		Path: config.IEDRIVER,
	}
	caps.AddChrome(chromeCaps)
	self.Service = service
	self.caps = caps
	self.port = port
	self.host = fmt.Sprintf("http://localhost:%d/wd/hub",port)
	self.saveScript = ""
}


func pickUnusedPort() (int, error){
	addr, err := net.ResolveTCPAddr("tcp","127.0.0.1:0")
	if err != nil {
		return 0,nil
	}
	l,err := net.ListenTCP("tcp",addr)
	if err != nil {
		return 0,err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0,err
	}

	return port, nil
}

func (self *Sel) getRemote2() (selenium.WebDriver,error) {
	mutex.Lock()
	defer mutex.Unlock()

	if self.initial {
		wb, err := selenium.NewRemote(self.caps,self.host)

		if err == nil {
			return wb,nil
		}
	}

	if self.initial && self.Service != nil {
		self.initial = false
		self.sessionGroup.Wait()
		self.Service.Stop()
	}

	port,_ := pickUnusedPort()
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(config.CHROMEDRIVER,port,opts...)

	if err != nil {
		logs.Log.Critical("[WebDriver Error] start on port:%d %v",port,err)
		self.initial = false
		return nil,err
	}

	logs.Log.Critical("[WebDriver] start on port:%d",port)

	caps := selenium.Capabilities{"browserName":"chrome"}
	chromeCaps := ch.Capabilities{
		Args: []string{
			"--headless",
			"--start-maximized",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-extensions",
			"--disable-accelerated-jpeg-decoding",
			"--test-type=ui",
		},
		Prefs: map[string]interface{} {
			"browser.download.dir": config.FILE_DIR,
			"download.default_directory": config.FILE_DIR,
			"download.prompt_for_download": false,
			"download.directory_upgrade": true,
			"safebrowsing.enabled": true,
		},
	}
	caps.AddChrome(chromeCaps)
	self.Service = service
	self.caps = caps
	self.port = port
	self.host = fmt.Sprintf("http://localhost:%d/wd/hub",port)
	self.saveScript = ""
	self.initial = true

	return selenium.NewRemote(self.caps,self.host)
}

func (self *Sel) getRemote() (selenium.WebDriver,error){
	//here no need to lock
	if self.initial {
		wb, err := selenium.NewRemote(self.caps,self.host)

		if err == nil {
			return wb,nil
		}
	}

	return self.getRemote2()
}

func (self *Sel) freeOne() {
	<- self.maxReq
}

func (self *Sel) Download(req Request) (resp *http.Response,err error) {
	self.maxReq <- 1
	defer self.freeOne()

	param, err := NewParam(req)
	if err != nil {
		return nil, err
	}

	resp = param.writeback(resp)
	resp.Request.URL = param.url

	self.sessionGroup.Add(1)
	defer self.sessionGroup.Done()

	for i := 0; i < param.tryTimes; i++ {
		if i != 0 {
			time.Sleep(param.retryPause)
		}

		html, err := self.GetPageSource(req)
		if err != nil {
			continue
		}

		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusOK)
		resp.Body = ioutil.NopCloser(strings.NewReader(html))

		break
	}

	if err != nil {
		resp.StatusCode = http.StatusBadGateway
		resp.Status = err.Error()
	}

	return
}


func (self *Sel) GetPageSource(req Request) (string,error) {
	//session create here not outside as time.sleep make session easy timeout
	wd, err := self.getRemote()

	if err != nil {
		//wait some time
		time.Sleep(10*time.Second)
		wd, err = self.getRemote()
		if err != nil {
			logs.Log.Informational("[Crash]WebDriver crash: %v",err)
			return "", err
		}
	}


	if wd != nil {
		defer wd.Close()
	}


	err = wd.Get(req.GetUrl())
	if err != nil {
		return "", err
	}

	html, err := wd.PageSource()

	return html, err
}