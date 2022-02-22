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

var rateLimiter = time.Tick(500 * time.Millisecond)

func Fetch(url string) ([]byte, error) {
	<-rateLimiter
	client := &http.Client{}
	newUrl := strings.Replace(url, "http://", "https://", 1)
	req, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	req.Header.Add("cookie", "sid=2fc02ffb-424d-4682-b1c3-31c4bd265bcf; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1645542452; ec=t7Na8bIT-1645542456924-f031238cf1e81534899366; FSSBBIl1UgzbN7NO=5aZEePwpRF8Gi5U3s1FiK4WUckXRx6VJAyjOlIG.TkbDFvy0YzPPhvkQlsgZ_iNcnA_WrAF5mlpkH5LbLwpN20a; _exid=3qX+oWTIVDvB/CIBWrIYSrMjfTCoLJL10sj/rc1KG8eeNgfHs0JFIGcYas9LydOTBwE/UQZHdDalJuCpyM+j/w==; _efmdata=tBm3Fl4DBeTs696v1kWZEbBzNbdddqEDbfqvO+HiXDMa5jKo6/JW4MiApZbfkANktXlIFn835DYk9FQJr6mNEjHaw3Z8UESmXC/cCC/rR3g=; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1645543095; FSSBBIl1UgzbN7NP=53fs5BCEVd.GqqqmdEGWMPAkF4Yu.b3wjzzkpKnt.5TdhsRsTOODtBiMcn3im41hbubbhB5KAR41soHS1yyzOq6ql8GaG3CrjABHMelSR1lLsla7ICERBmrkmd5tYYQuSN4fu.Y6RcqxbszZ7gPeuz1QgBtistqwc7jI6LgZd4vTSRvnwo_zpQcK5o1hhIH7pHG5ipd99x48HSALtTjwELVYWxjdYxzktKIvzdiKBuPCmteqIsarIkJLO7i1xzHXcE")
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
