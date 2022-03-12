package internal

import (
	"encoding/json"
	"fmt"
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
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Cache-Control", "max-age=0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
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
