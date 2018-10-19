package main

import (
	"net/http"
)

func doReq(request *http.Request) (response *http.Response, err error) {
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
