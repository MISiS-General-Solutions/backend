package client

import (
	"MGS/osmdata"
	"MGS/shared"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	apiEndpoint = "https://overpass-api.de/api/interpreter"
	cacheFolder = "data/cache/"
)

type Client struct {
	apiEndpoint string
	httpClient  *http.Client
	settings    string
}

func NewClient() *Client {
	return &Client{
		apiEndpoint: apiEndpoint,
		httpClient:  http.DefaultClient,
	}
}
func (r *Client) SetSettings(settings string) {
	r.settings = settings
}
func (r *Client) Query(query string) ([]byte, error) {
	resp, err := r.httpClient.PostForm(
		r.apiEndpoint,
		url.Values{"data": []string{query}},
	)
	if err != nil {
		return nil, err
	}
	defer shared.Discard(resp)
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%v: %s", resp.StatusCode, body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (r *Client) GetData(req string) (*osmdata.OsmData, error) {
	fullRequest := r.settings + req
	hashSum := md5.Sum([]byte(fullRequest))
	cache, err := LoadResult(fmt.Sprintf("%v%x.osm", cacheFolder, hashSum))
	if err != nil {
		fmt.Println("not cached")
		cache, err = r.Query(fullRequest)
		if err != nil {
			return nil, err
		}
		if err = SaveResult(cache, fmt.Sprintf("%v%x.osm", cacheFolder, hashSum)); err != nil {
			return nil, err
		}
	}
	return osmdata.ReadOSM(bytes.NewReader(cache))
}
func (r *Client) GetObstacles() (*osmdata.OsmData, error) {
	obstacleRequest := `(
		way["building"];
		relation["building"];
	  );
	  out skel qt;
	  >;
	  out skel qt;`
	return r.GetData(obstacleRequest)
}
func (r *Client) GetRoads() (*osmdata.OsmData, error) {
	roadRequest := `(way["highway"];node(w););out qt;`
	return r.GetData(roadRequest)
}
func SaveResult(res []byte, path string) error {
	return os.WriteFile(path, res, 0600)
}

func LoadResult(path string) ([]byte, error) {
	res, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return res, nil
}
