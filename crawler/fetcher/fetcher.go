package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var rateLimiter = time.Tick(10 * time.Millisecond)

func Fetch(url string) ([]byte, error) {
	<-rateLimiter
	client := &http.Client{}
	newUrl := strings.Replace(url, "http://", "https://", 1)
	req, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	req.Header.Add("cookie", "sid=e2dfdde9-a796-4864-9c83-9780584a43ee; ec=WNfofO7o-1642522510650-aff210bca6d541352711054; FSSBBIl1UgzbN7NO=52OsBuNaDzQteWqO0cRH.NAVPm92J4iMUd9WPtcsGti6nBXGGgWItVGA6P5JZN.I9q2a97W6BEBClrrLVANvLIq; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642522526; _exid=HbtnWDiftrTREtAWgEk1oHYeDkpESwVDzhDVLG1kmfzVrwM8UtwcE/H8gSmXxzZOtiEFo//jEe0KZ0mPerZ8bA==; _efmdata=ISFqPTI7W2q1c09hdpwQGkBS2ZvD4dEQqhGAWW0ajeNhVYjyVBD9IUcExOlaWK2P84mIeOmKGLq49OeE0uA9ejBM40CDdHoHr9s99dntQnc=; FSSBBIl1UgzbN7NP=538MTOKobe3Qqqqm5MPZSyG7coZL8KTXOjsqtDWVDGQsKcDMvb3yFt8xq7k8ctWpZkaq9L.SufA9Y6kVhPa0diuc5f5yF8rgQn5084hmpFHw0PFHsBYjxt.yw3W0smKPzMC.R1T8aDOVIe2CXczRhll7u4POqVDi.0RCiMu44kePriRtSe8a4N3cRDKjYNBodQQzGPc9in7yqs4MJkHXWU2PTdFEsjHxbIfPDFe9_9ZwDL9Y1bcda1Z_nCFAu9.2da; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642526088")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
