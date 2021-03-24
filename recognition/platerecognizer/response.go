package platerecognizer

import "time"

type Response struct {
	ProcessingTime float64 `json:"processing_time"`
	Results        []struct {
		Box struct {
			Xmin int `json:"xmin"`
			Ymin int `json:"ymin"`
			Xmax int `json:"xmax"`
			Ymax int `json:"ymax"`
		} `json:"box"`
		Plate  string `json:"plate"`
		Region struct {
			Code  string  `json:"code"`
			Score float64 `json:"score"`
		} `json:"region"`
		Score      float64 `json:"score"`
		Candidates []struct {
			Score float64 `json:"score"`
			Plate string  `json:"plate"`
		} `json:"candidates"`
		Dscore  float64 `json:"dscore"`
		Vehicle struct {
			Score float64 `json:"score"`
			Type  string  `json:"type"`
			Box   struct {
				Xmin int `json:"xmin"`
				Ymin int `json:"ymin"`
				Xmax int `json:"xmax"`
				Ymax int `json:"ymax"`
			} `json:"box"`
		} `json:"vehicle"`
	} `json:"results"`
	Filename  string      `json:"filename"`
	Version   int         `json:"version"`
	CameraID  interface{} `json:"camera_id"`
	Timestamp time.Time   `json:"timestamp"`
}
