package discovery

import "encoding/json"

// 端点：
// ipconfig中的端点是gateway server
type Endpoint struct {
	Ip       string                 `json:ip`
	Port     string                 `json:port`
	MetaData map[string]interface{} `json:metadata`
}

func Marshal(endpoint *Endpoint) string {
	data, err := json.Marshal(endpoint)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Unmarshal(data []byte) (*Endpoint, error) {
	endpoint := &Endpoint{}
	err := json.Unmarshal(data, endpoint)
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}
