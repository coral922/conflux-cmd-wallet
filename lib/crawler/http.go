package crawler

import (
	"cfxWorld/lib/util"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

var httpClient = util.DefaultHttpClient()

func SetClient(client *http.Client) {
	httpClient = client
}

func Get(url string, params ...map[string][]string) (content []byte, err error) {
	req, err := http.NewRequestWithContext(context.TODO(), "GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params[0] {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "请求失败%s", req.URL)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("返回码错误: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "未知响应%s", req.URL)
	}
	return body, nil
}