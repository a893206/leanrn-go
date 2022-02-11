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
	req.Header.Add("cookie", "sid=8b7c2261-7ed5-4838-abba-562fc11507b1; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1644603887; ec=9rUSJZoo-1644603894919-9fe63703ce621-1502667137; FSSBBIl1UgzbN7NO=5IXzbntbJ87h5UnPdL.g2oB_qDEPN0kSOmE0RbhnNOgRIsiVGc2NNwnzUfxJYJJTm7WipjWe.85FoxGt9kHP_gG; _exid=Z/xFNJVbwwq/NHVrspvWQlIt4x2JEsqf6Wb0b6tpYFcsx3WL15cQiGOBVeSWFRo7i8BHgED/loa9xuWTYbzZ3g==; _efmdata=6+d7P05WU/IcoKeRKgZGMBxNQSm81jkwATq4tKYdNpDWIEfmg3Q6aiLbMQJwtarMCrb23pnNqFk8zY4j9rc+Z3qFfNnRKzdqVc8QnpDnT+w=; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1644603951; FSSBBIl1UgzbN7NP=53GFjHCEDjOLqqqm5_G9xaa9X048bRqQPV5Z.XVUcYbcED0H_d4JquDRQrUTT8id_rKABD6_ZKn38UokcYaMTbWgurYqBFWiLObJNiT0vUTj.xwwJGdwGyY5w3rspklWVP9MEU7MF1MjTDYumIc064ZeghM_tvzl_YJFj79yyliLl3YfoaWh8fghdgr2QRBk1.6B2u.skrc6viDBnr6HOARB5rJ2K7JlS_rp58pgKEB7irDWDdnu4AN5SMgY7Td9KW")
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
