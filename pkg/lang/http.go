package lang

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	emptyMapInter = map[string]interface{}{}
)

// PostForm http post (form) request
func PostForm(URL string,
	data map[string]interface{},
) (map[string]interface{}, error) {
	uvs := data2urlValues(data)
	resp, err := http.PostForm(URL, uvs)
	if err != nil {
		return emptyMapInter, err
	}

	var (
		respData map[string]interface{}
	)

	bs, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	respStr := string(bs)

	if !IsJSON(respStr) {
		respData = ParseURLQuery(respStr)
	} else {
		respData = ParseJSON(respStr)
	}

	return respData, nil
}

// IsJSON Judge a string is json string or not
func IsJSON(s string) bool {
	return json.Valid([]byte(s))
}

// ParseURLQuery parse url query string to map[string]interface{}
// k=v&a=1&b=2
func ParseURLQuery(query string) map[string]interface{} {
	uvs, err := url.ParseQuery(query)
	if err != nil {
		// println(err)
		return emptyMapInter
	}

	respData := make(map[string]interface{})
	for key := range uvs {
		val := uvs.Get(key)
		tmpData := map[string]string{}
		if err := json.Unmarshal([]byte(val), &tmpData); err != nil {
			respData[key] = val
		} else {
			// println("this called with key:", key)
			respData[key] = tmpData
		}
	}
	return respData
}

// ParseJSON parse json string to map[string]interface{}
//
// json string like:
// {"a":"a","b":"b"}
//
func ParseJSON(s string) map[string]interface{} {
	respData := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &respData); err != nil {
		println(err)
		return emptyMapInter
	}
	return respData
}

func data2urlValues(data map[string]interface{}) url.Values {
	uvs := url.Values{}
	for key, val := range data {
		if valStr, ok := val.(string); !ok {
			bs, _ := json.Marshal(val)
			uvs.Add(key, string(bs))
		} else {
			uvs.Add(key, valStr)
		}
	}
	return uvs
}

// CopyRequest means to copy a request for another handle-func or else
func CopyRequest(req *http.Request) *http.Request {
	body, _ := ioutil.ReadAll(req.Body)
	rdOnly := ioutil.NopCloser(bytes.NewBuffer(body))

	newReq, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	newReq.Header = req.Header
	req.Body = rdOnly
	return newReq
}
