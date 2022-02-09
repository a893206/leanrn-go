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
	req.Header.Add("cookie", "sid=95e3e0d6-5508-412e-95b5-8756254ebf85; ec=EYFYbnjU-1644419471300-0789b186157031573668907; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1644419490; FSSBBIl1UgzbN7NO=5vCUGTT.kaf591FDb0I1.p4Bh2oet53.WRfOfJx0xH0_QiW39c2D7jcV.Qo6_8y6r1ZV6h6gtKTPY13W7rJQMPa; _exid=0G7/qWv0B4EqgQLCWVFRDtaMC3tEU4iDxpqW2XmSb+RqAUwhrGbOiWDmYocEt50weZXh6KTQ2Dzh4gUmDw15Tg==; _efmdata=4dBBtGbDCWd95Wq17uqANfMWNVLW48/orf+YWQIEy3bkuJrUuk5KJW2iVQtgCq0gHoHOJvU4ot7nJzz32WHB6tWUapkAc6EXoJm0o9QWopg=; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1644420238; FSSBBIl1UgzbN7NP=53GHQ9CEk6Egqqqm56SKPmavayJmJoPDtX42_6jo8Fc4u.DvtKPOSixGkuNhg5TROTXeNtMiQpPtP2KPgKHDgXIQzzcELNvxe1q7t8MQa1bFsRdiL1FISnTRtJV4uiZP4E8D0huv5sjfmSatV7AdK_8kmqOwcdjLWIQwbtdblS7BWHnK1TlUujDSK440r2N_94iWaHNNepVO114mwWBOnegP60thKkShme5KwL6dc6_MhwSZeXVJ.V2xO1HMDaG9tG")
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
