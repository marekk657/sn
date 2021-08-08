package snackable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type AllResponse struct {
	FileID           string `json:"fileId"`
	ProcessingStatus string `json:"processingStatus"`
}

type AllResponses []AllResponse

func (a AllResponses) FindByID(fileID string) (AllResponse, bool) {
	for _, v := range a {
		if v.FileID == fileID {
			return v, true
		}
	}
	return AllResponse{}, false
}

type DetailsResponse struct {
	FileID           string `json:"fileId"`
	FileName         string `json:"fileName"`
	MP3Path          string `json:"mp3Path"`
	OriginalFilePath string `json:"originalFilePath"`
	SeriesTitle      string `json:"seriesTitle"`
}

type SegmentResponse struct {
	FileSegmentID int    `json:"fileSegmentId"`
	FileID        string `json:"fileId"`
	SegmentText   string `json:"segmentText"`
	StartTime     int    `json:"startTime"`
	Endtime       int    `json:"endTime"`
}

type Interface interface {
	All(limit, offset int) (AllResponses, error)
	Details(fileID string) (DetailsResponse, error)
	Segments(fileID string) ([]SegmentResponse, error)
}

func NewClient(baseURL string) Client {
	return Client{
		baseUrl: baseURL,
	}
}

type Client struct {
	baseUrl string
}

func (c Client) All(limit, offset int) (AllResponses, error) {
	endpoint := fmt.Sprintf("%s/api/file/all", c.baseUrl)

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Error().Err(err).Str("raw_url", endpoint).Msg("failed to parse url")
		return nil, err
	}

	q := u.Query()
	if limit > 0 {
		q.Set("limit", fmt.Sprint(limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprint(offset))
	}
	u.RawQuery = q.Encode()

	var allResponse []AllResponse
	if err := c.doRequest(u.String(), &allResponse); err != nil {
		log.Error().Err(err).Msg("failed to query /all endpoint")
		return nil, err
	}

	return allResponse, nil
}

func (c Client) Details(fileID string) (DetailsResponse, error) {
	endpoint := fmt.Sprintf("%s/api/file/details/%s", c.baseUrl, fileID)

	var detailsResponse DetailsResponse
	if err := c.doRequest(endpoint, &detailsResponse); err != nil {
		log.Error().Err(err).Msg("failed to query details endpoint")
		return DetailsResponse{}, err
	}

	return detailsResponse, nil
}

func (c Client) Segments(fileID string) ([]SegmentResponse, error) {
	endpoint := fmt.Sprintf("%s/api/file/segments/%s", c.baseUrl, fileID)

	var segmentsResponse []SegmentResponse
	if err := c.doRequest(endpoint, &segmentsResponse); err != nil {
		log.Error().Err(err).Msg("failed to query segments endpoint")
		return nil, err
	}

	return segmentsResponse, nil
}

func (c Client) doRequest(url string, model interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("failed to query endpoint")
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(model); err != nil {
		log.Error().Err(err).Str("url", url).Msg("failed to read json input")
		return err
	}
	return nil
}
