package surfer

import "testing"

func TestOpen(t *testing.T){
	req := &DefaultRequest{
		Url: "http://www.baidu.com/",DownloaderID:2,
		DialTimeout:3000}

	sf := NewChromeDP()
	sf.Download(req)

	req = &DefaultRequest{
		Url: "http://www.google.com/",DownloaderID:2,}

	sf.Download(req)
}
