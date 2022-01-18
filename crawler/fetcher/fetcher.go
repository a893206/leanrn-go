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
)

func Fetch(url string) ([]byte, error) {
	client := &http.Client{}
	newUrl := strings.Replace(url, "http://", "https://", 1)
	req, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	req.Header.Add("cookie", "sid=e2dfdde9-a796-4864-9c83-9780584a43ee; ec=WNfofO7o-1642522510650-aff210bca6d541352711054; FSSBBIl1UgzbN7NO=52OsBuNaDzQteWqO0cRH.NAVPm92J4iMUd9WPtcsGti6nBXGGgWItVGA6P5JZN.I9q2a97W6BEBClrrLVANvLIq; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642522526; _exid=hAXoXW4vrBeQ+hWfYR2e4hrZIY5hfz4QNMdCr9ILSvO5D4ifjWuvqlU646vBjoOGPkaPYa1mW+IPs1zC7+yOjw==; _efmdata=ISFqPTI7W2q1c09hdpwQGkBS2ZvD4dEQqhGAWW0ajeNhVYjyVBD9IUcExOlaWK2P9Mhp3bY3wCXGoeEqwWZU6KFdYIqQKpy8iTIrHqbks9I=; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642522588; FSSBBIl1UgzbN7NP=538M.Jbob_ulqqqm5MGMicaX8a09muC0Hiuzcc7moAMeBqESChXTis7sKzJEME5Rfc86MNGm5HjyT_PE5lguGiT7S1RmbXTpgw92efGWY8dRzlSBiSyfKPremDxTE1JyTcgBbmm0.w4vdewlBKXe4Z03mmamk0ebeAXCrcyGKMxtrxGtElw5lYxduj9BE5spnGJOmLIpPkOsOdQcP4ipzRk3D91XQwWHL7_CPPCYo6EDAinewwrfBeIwgyQsC317gA")
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
