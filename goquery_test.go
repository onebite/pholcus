package main

import (
	"testing"
	"io/ioutil"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"bytes"
	"strings"
	"net/http"
	"os"
	"io"
	"github.com/henrylee2cn/pholcus/app/pipeline/collector/data"
	"time"
	"regexp"
	"fmt"
)

func TestGoqueryFirst(t *testing.T){
	data,err:= ioutil.ReadFile("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\logs\\12.html")
	if err != nil {
		t.Logf("%v", err)
	}

	query,err:=goquery.NewDocumentFromReader(bytes.NewReader(data))
	wrapper := query.Find("td.folder")
	wrapper.Each(func(i int, selection *goquery.Selection) {
		html,_ := selection.Html()
		t.Logf("text: %v",html)
		href,_ := selection.Find("a").Eq(0).Attr("href")
		t.Logf("href: %v",href)
	})

	dd,err:= ioutil.ReadFile("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\logs\\33.html")
	if err != nil {
		t.Logf("%v", err)
	}

	qq,err:=goquery.NewDocumentFromReader(bytes.NewReader(dd))
	detail,_ := qq.Find("div.t_msgfont").Html()
	brs := strings.Split(detail,"<br/>")
	t.Logf("title: %v",brs[0])
	t.Logf("format: %v",brs[1])
	qq.Find("div.t_msgfont").Find("img").Each(func(i int, selection *goquery.Selection) {
		src,_ := selection.Attr("src")
		t.Logf("src: %v",src)
	})

	//http://68.168.16.147/bbs/attachment.php?aid=3166464
	hq,_ := qq.Find("dl.t_attachlist>dt>a").Eq(1).Attr("href")
	t.Logf("hq: %v",hq)

}


func TestGoquery(t *testing.T){
	data,err:= ioutil.ReadFile("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\logs\\hh.html")
	if err != nil {
		t.Logf("%v", err)
	}

	query,err:=goquery.NewDocumentFromReader(bytes.NewReader(data))
	wrapper := query.Find("div.wrapper-product")
	html,_ := wrapper.Html()
	t.Logf("text: %v",html)
	wrapper.Find("table>tbody>tr").Each(func(i int, selection *goquery.Selection) {
		ht,_ := selection.Children().Eq(0).Html()
		t.Logf("tr html :%v",ht)
		tdText := selection.Children().Eq(0).Text()
		if strings.Contains(tdText,"発売"){
			t.Logf("発売: %v",selection.Children().Eq(1).Text())
		}
		if strings.Contains(tdText,"出演者"){
			t.Logf("出演者: %v",selection.Children().Eq(1).Text())
		}
		if strings.Contains(tdText,"品番"){
			t.Logf("品番: %v",selection.Children().Eq(1).Text())
		}
	})

	query.Find("table.mg-b20").Find("tr").Each(func(i int, selection *goquery.Selection) {
		//ht,_ := selection.Children().Eq(0).Html()
		tdText := selection.Children().Eq(0).Text()
		if strings.Contains(tdText,"発売"){
			t.Logf("発売2: %v",selection.Children().Eq(1).Text())
		}
		if strings.Contains(tdText,"出演者"){
			t.Logf("出演者2: %v",selection.Children().Eq(1).Text())
		}
		if strings.Contains(tdText,"品番"){
			t.Logf("品番2: %v",selection.Children().Eq(1).Text())
		}
	})

	query.Find("img.tdmm").Each(func(i int, selection *goquery.Selection) {
		src,ok:= selection.Attr("src")
		if ok{
			t.Logf("img tdmm:%v",src)
		}
	})

	query.Find("img.mg-b6").Each(func(i int, selection *goquery.Selection) {
		src,ok:= selection.Attr("src")
		if ok{
			t.Logf("img mg-b6:%v",src)
		}
	})
}

func TestGoqueryOne(t *testing.T){
	data,err:= ioutil.ReadFile("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\logs\\one.html")
	if err != nil {
		t.Logf("%v", err)
	}

	query,err:=goquery.NewDocumentFromReader(bytes.NewReader(data))
	wrapper := query.Find("div.columns")
	wrapper.Each(func(i int, selection *goquery.Selection) {
		fanhao := selection.Find("div>h5>a").Text()
		publish := selection.Find("div>p>a").Text()
		actress := selection.Find("a.panel-block").Text()
		src,_ := selection.Find("div>img.image").Eq(0).Attr("src")
		torrent,_ := selection.Find("a[rel='nofollow']").Eq(0).Attr("href")
		t.Logf("品番: %v,时间:%v, 演员:%v,图片:%v,文件:%v",
			strings.TrimSpace(fanhao),
				strings.Replace(strings.TrimSpace(publish),"\n","",-1),
			strings.Replace(strings.TrimSpace(actress),"\n","",-1),src,torrent)
	})
}


func TestDownload(t *testing.T) {
	rep,err :=http.Get("https://pics.dmm.co.jp/mono/movie/adult/tknnpj312/tknnpj312pt.jpg")
	if err == nil {
		body,err:=ioutil.ReadAll(rep.Body)
		if err == nil {
			out,_ := os.Create("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\file_out\\dmm12.jpg")
			io.Copy(out,bytes.NewReader(body))
			out.Close()
		}
	}
}

func TestDownload2(t *testing.T) {
	rep,err :=http.Get("https://onejav.com/torrent/gvg768/download/558410.torrent")
	if err == nil {
		body,err:=ioutil.ReadAll(rep.Body)
		dis := rep.Header["Content-Disposition"]
		t.Logf("文件：%v",dis)
		if err == nil {
			out,_ := os.Create("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\file_out\\558410.torrent")
			io.Copy(out,bytes.NewReader(body))
			out.Close()
		}
	}
}

func TestChanDownload(t *testing.T){
	rep,err :=http.Get("https://pics.dmm.co.jp/mono/movie/adult/tknnpj312/tknnpj312pt.jpg")
	var fc chan data.FileCell

	fc = make(chan data.FileCell,10)
	if err == nil {
		body,err:= ioutil.ReadAll(rep.Body)
		if err == nil {
			go func() {
				for rc := range fc {
					fileName := rc["Name"].(string)
					out,_ := os.Create(fileName)
					io.Copy(out,bytes.NewReader(rc["Bytes"].([]byte)))
					out.Close()
				}
			}()

			f := data.FileCell{
				"RuleName":"testtest",
				"Name":"C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\file_out\\dmm112.jpg",
				"Bytes":body,
			}
			fc <- f

			ff := data.FileCell{
				"RuleName":"testtest",
				"Name":"C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\file_out\\dmm1432.jpg",
				"Bytes":body,
			}
			fc <- ff
		}

	}

	time.Sleep(10*time.Second)
}


func TestRegexp(t *testing.T) {
	text := `2018/12/20 12:18:39 [C] [
	<img src="http://img1.uploadhouse.com/fileuploads/23733/23733491cd87f974683a6694c823c87c5224cd69.jpg" border="0" onclick="zoom(this)" onload="attachimg(this, &#39;load&#39;)" alt="" width="565" style="cursor: pointer;"/>
	sfsjflsf;sjf;
<a href="http://www.uploadhouse.com/viewfile.php?id=23733490&amp;showlnk=0" target="_blank"><img src="http://img0.uploadhouse.com/fileuploads/23733/23733490e2e20feeb69024152d8ca01a9fdf376e.jpg" border="0" onclick="zoom(this)" onload="attachimg(this, &#39;load&#39;)" alt="" width="565" style="cursor: pointer;"/></a><br/>
<a href="http://www.uploadhouse.com/viewfile.php?id=23733491&amp;showlnk=0" target="_blank"><img src="http://img1.uploadhouse.com/fileuploads/23733/23733491cd87f974683a6694c823c87c5224cd69.jpg" border="0" onclick="zoom(this)" onload="attachimg(this, &#39;load&#39;)" alt="" width="565" style="cursor: pointer;"/></a></font></font>`
   reg :=regexp.MustCompile("<a.*/a>|</font>|<img.*src.*/>")
   text = reg.ReplaceAllString(text,"")
   fmt.Println("result:%s",text)
}
