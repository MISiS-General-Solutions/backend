package shared

import (
	"encoding/json"
	"io"
	"net/http"
)

func Discard(resp *http.Response) {
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}
func MustMarshallToJSON(val interface{}) []byte {
	b, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return b
}
