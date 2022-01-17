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
	req.Header.Add("cookie", "FSSBBIl1UgzbN7NO=5g0kms2DdrTKoPqMolN3VFY1Q7.f6U9LG9x21DqptLGyraGXUFMKo7WbPG1jE1OTaj8bFb2_F99o5S8XnEDePVG; sid=c0f3f23f-3b29-4ac6-8c37-824558f4f83a; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642442447; ec=XGwj9Rj9-1642442447539-2ea6d76763abf1393945480; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1642443379; _exid=kB3fKD8RXGQhR+sAZrgdi6ZAe6jltKkW3uM3H398poI4GSNhGVgZavtfLvwu6jayRqcBvEL0nbwFrlQeLZtoQw==; _efmdata=VZDIG1Hs6dUkoxnKWKoz8Qhw/YqJtBCtS7EmxwoSS/tBNhi81SfK3KHu9HUJ+t6MUPvy/0KDPSjOMq8g6O6m5QAZ+MEFRrj1dvHC/zd9H/0=; FSSBBIl1UgzbN7NP=538Kn2Do27Lgqqqm5FVDVrGEsyMnJx9_GOO1LaIvRGTKqUzRyt7De6Dw8hdpu24GL2oSuPK0pScW_t.3kFCsppmfCDOmLu.sO_uOMM_hBODxYQdVmE.FYBFKnXl2THfL_Q58Aoc.IreTqzU1GXwU4Ne5UkomFb2rXqpfk9NcBNlzdcqMIoIS1kZJxMUgVVoUyaPWhwT.td47ZAH3ikRewnXpAYqE6TsxnT7hXImJ4ISw_Z9uQe3hWgUH6xIBZ0w7PV")
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
