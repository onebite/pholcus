package surfer

import (
	"testing"
	"fmt"
	"strings"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"os"
	"bufio"
	"github.com/steakknife/bloomfilter"
	"hash/fnv"
	"regexp"
)

func TestSeleniumOpen(t *testing.T){
	req := &DefaultRequest{
		Url: "https://v.qq.com/",DownloaderID:3,
		DialTimeout:3000}

	sf := NewSel()
	sf.Download(req)

	req = &DefaultRequest{
		Url: "https://www.youtube.com",DownloaderID:3,
		DialTimeout:1000}

	sf.Download(req)
}


func TestDownload(t *testing.T) {
	req := &DefaultRequest{
		Url: "https://pics.dmm.co.jp/digital/video/mide00616/mide00616pl.jpg",
		//Url:"https://v.qq.com/",
		DownloaderID:3,
		DialTimeout:3000}
	sf := NewSel()
	_,err := sf.Download(req)
	if err != nil {
		fmt.Println("download error %v",err)
		return
	}

	//_,s := path.Split(req.GetUrl())
	//f, err := os.OpenFile(filepath.Join("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\file_out",s),
	//	os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer f.Close()
	//size, err := io.Copy(f,resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("download size:%d",size)

}

func TestDownload2(t *testing.T) {
	req := &DefaultRequest{
		Url: "http://wx1.sinaimg.cn/large/ed5e6a1dly1fy8b3rduxhj20go0oo0yq.jpg",
		//Url:"https://v.qq.com/",
		DownloaderID:3,
		DialTimeout:30000}

	Download(req)
}

func TestReplace(t *testing.T)  {
	st := "strings.Replace<br/><br><br/>2222##<br>"
	st = strings.Replace(st,"<br/>","\n",-1)
	fmt.Println(st)
}

func TestRegexp(t *testing.T) {
	r := regexp.MustCompile("[\u0391-\uFFE5]+")
	splits := r.Split("中国人afjls日本人pppp",-1)
	fmt.Println(splits)
	fmt.Println(splits[0])
}


func TestBloomFilter(t *testing.T) {
	opt,_:=bloomfilter.NewOptimal(300000000,0.0000001)
	h := fnv.New64()
	h.Write([]byte("jfsjflsjfsljfsljf"))

	for i := 0;i < 1000;i++{
		fmt.Printf("contains:%v \n",opt.Contains(h))
	}
	fmt.Println(regexp.MustCompile("[-_A-Za-z0-9]+$").Match([]byte("ace-123-a")))
	fmt.Println(regexp.MustCompile("[-_A-Za-z0-9]+$").Match([]byte("    ")))
	fmt.Println(regexp.MustCompile("[-_A-Za-z0-9]+$").Match([]byte("12345")))
	fmt.Println(regexp.MustCompile("[-_A-Za-z0-9]+$").Match([]byte("__----")))
}

func TestGoQuery(t *testing.T) {
	f, err := os.Open("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\logs\\bbs_forum-58-1.html")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		fmt.Println(err)
	}

	a := doc.Find("td.folder>a")
	a.Each(func(i int, selection *goquery.Selection) {
		href,_ := selection.Attr("href")
		fmt.Printf("href %d %v \n",i,href)
	})

	ff, err := os.Open("C:\\workspace\\gosource\\src\\github.com\\henrylee2cn\\pholcus\\pholcus_pkg\\logs\\33.html")
	defer ff.Close()
	if err != nil {
		fmt.Println(err)
	}
	query, _ := goquery.NewDocumentFromReader(bufio.NewReader(ff))
	detail,_ := query.Find("div.t_msgfont").Html()

	post,_:= query.Find("td.postcontent>div>h2").Html()
	zimu := "日语"
	if strings.Contains(post,"中文") || strings.Contains(post,"中字"){
		zimu = "中文"
	}

	var fanhao string
	if post != ""{
		splits := strings.Split(post," ")
		for _,value := range splits {
			if strings.Contains(value,"-") {
				fanhao = value
				break
			}
		}
	}

	brs := strings.Split(detail,"<br/>")
	var out map[int]interface{}

	out = make(map[int]interface{})
	out[0] = fanhao
	out[1] = zimu

	for i,br := range brs {
		out[i+2] = strings.Replace(strings.TrimSpace(br),"\n","",-1)
	}

	query.Find("div.t_msgfont").Find("img").Each(func(i int, selection *goquery.Selection) {
		src,_ := selection.Attr("src")
		fmt.Printf("image:%v \n",src)
	})

	torrent,_:= query.Find("dl.t_attachlist>dt>a").Eq(1).Attr("href")
	if !strings.Contains(torrent,"http") {
		torrent = "http://68.168.16.147/bbs/" + torrent
	}

	fmt.Printf("torrent:%v \n",torrent)

	fmt.Println(out)
}