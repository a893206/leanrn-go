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
	req.Header.Add("cookie", "sid=caf9cc2a-748f-4b33-9057-1a1d2fcb6217; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642907537; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642907543; ec=Xbz4oAEm-1642907546431-c44c526fc5977-1861293040; _exid=6A7e+BYwQE/DC7lL6ZRWPvNGV2bkL8Dd3ij5bqpetrre8sTMTtyQyXe676SvRWwS79k9CjvRZChXSRl5Omhe1Q==; FSSBBIl1UgzbN7NO=5F.anrCfNKU88rGutyp6FnolrWkpuj3DJKMPHTSG2aGWQt1KuIvvMeUUBKJq.c6QLc0jjhmEgnFIF.fdkbZURva; _efmdata=AC1b8gZVOJPWPMytDn4AEIKZ5Tj9U8Uaj7rA2hLfjlavSZi+wUg0BBqu6oSJBDBOOybevtKNIHwk6+O/cdAmgF/+01sZqwo1IUNV5HUgpio=; FSSBBIl1UgzbN7NP=538uojDo7S_Eqqqm5fiBNxaPs3VmEJYEG_F1eVIFupEgg5Y5gwRWvZafHiqdJj8HTShNLpCZ2gt6YxM5WNP8l1RgxHc46I1fNVBG.1vBobcskr8ISlA54PgUxV1p_j0SYQkPIKTduP_3bPc.IL816aZIvJ7WCKRrqzeEfsTzvmAfwdbKuqvPsQOZ4hWUbBWeQGg13ec.xWuaezYprbRjheit_7hpBLtd58wKwkDjNcFX1byoqsjZkKorxpOvhVSLRa")
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
