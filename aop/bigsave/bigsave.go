package bigsave

import (
	"bytes"
	"encoding/gob"
)

func Encoder(data []interface{}) string {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(&data)
	if err != nil {
		return ""
	}
	b := GZipBytes(buf.Bytes())
	return string(b[:])
}

func Decoder(i string) []interface{} {
	byteEn := UGZipBytes([]byte(i))
	decoder := gob.NewDecoder(bytes.NewReader(byteEn))
	ret := []interface{}{}
	err := decoder.Decode(&ret)
	if err != nil {
		return nil
	}
	return ret
}
