package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func MakeUserRequest(address string) (int, error) {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		log.Println("Error:", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:97.0) Gecko/20100101 Firefox/97.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "ru,en-US;q=0.7,en;q=0.3")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Proxy-Authorization", "Basic LnN2QDU4NTM0Mzt1YS46SER2TkZYMm5ubVVPZTYxMEFxeFlGeDNFKzdlb1hXSjBxOE5Pd0RDa1ZUWT0=")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", "csrf-token-name=csrftoken; csrf-token-value=16d74d201b416182f109d9108f39c5ddcb8e4ce5dd8c3875b3efb1d1f818fb05c008cf36d403bc34; sputnik_session=1645868989999|1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func GetMyIp2() (it IntroduceType, err error) {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}

	json.Unmarshal(body, &it)

	return
}

func GetMyIp() (it IntroduceType, err error) {
	client := fasthttp.Client{}

	body := make([]byte, 0)

	var status int

	status, body, err = client.Get(body, "http://ip-api.com/json")
	if err != nil {
		err = errors.Wrap(err, "cannot perform request for getting IP address")
		return
	}

	if status != 200 {
		err = fmt.Errorf("got status %v", status)
		return
	}

	err = json.Unmarshal(body, &it)

	if err != nil {
		err = errors.Wrap(err, "cannot parse response from server")
	}

	return
}
