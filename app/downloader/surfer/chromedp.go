package surfer

import (
	"github.com/chromedp/chromedp"
	"context"
	"log"
	"net/http"
	"strings"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/dom"
	"io/ioutil"
	"time"
	"sync"
)

var mutex sync.Mutex

type Chrome struct {
	Cdp *chromedp.CDP
	Pool *chromedp.Pool
	ctx *context.Context
	cacel context.CancelFunc
	initial bool
}

func NewChromeDP() Surfer {
	chrome := &Chrome{
		initial:false,
	}


	return chrome
}

func (self *Chrome) initC() {
	if self.initial {
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	ctx,cancelF := context.WithCancel(context.Background())

	cdp,err := chromedp.New(ctx)
	if err != nil {
		log.Printf("[E] chromedp init error: %v",err)
	}

	self.Cdp = cdp
	self.ctx = &ctx
	self.cacel = cancelF
	self.initial = true
}


func (self *Chrome) Download(req Request) (resp *http.Response,err error) {
	//self.initC()
	//
	//defer self.cacel()
	self.initC()
	param, err := NewParam(req)
	if err != nil {
		return nil, err
	}

	resp = param.writeback(resp)
	resp.Request.URL = param.url


	for i := 0; i < param.tryTimes; i++ {
		if i != 0 {
			time.Sleep(param.retryPause)
		}

		err := self.Cdp.Run(*self.ctx,doTask(req,resp))

		if err != nil {
			continue
		}

		break
	}

	if err == nil {
		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusOK)
	} else {
		resp.StatusCode = http.StatusBadGateway
		resp.Status = err.Error()
	}


	return
}

//TODO chromedp current version do not support download files
func doTask(req Request,resp *http.Response) chromedp.Tasks {
	visibleById := req.GetTemp("VisibleById","")
	var visible chromedp.Action
	if visibleById == nil {
		visible = chromedp.Sleep(req.GetDialTimeout())
	}else {
		visible = chromedp.WaitVisible(visibleById,chromedp.ByID)
	}
	return chromedp.Tasks{
		//network.Enable(),
		//network.SetExtraHTTPHeaders(getHeaders(req)),
		chromedp.Navigate(req.GetUrl()),
		visible,
		chromedp.ActionFunc(func(ctx context.Context,h cdp.Executor) error{
			html,err := dom.GetOuterHTML().Do(ctx,h)
			if err != nil{
				resp.StatusCode = http.StatusNotFound
				resp.Status = err.Error()
			}else {
				resp.StatusCode = http.StatusOK
				resp.Status = http.StatusText(http.StatusOK)
				resp.Body = ioutil.NopCloser(strings.NewReader(html))
			}
			return nil
		}),
	}
}

func getHeaders(req Request) network.Headers {
	headers := make(map[string]interface{})
	for key,value := range req.GetHeader(){
		headers[key] = value
	}

	return headers
}