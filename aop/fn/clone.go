package fn

import "encoding/json"

func Copy(src interface{}, desc interface{}) error {
	data, _ := json.Marshal(src)
	return json.Unmarshal(data, desc)
}
