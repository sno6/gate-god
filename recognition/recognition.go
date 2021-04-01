package recognition

import "io"

type PlateResult struct {
	Score float64 `json:"score"`
	Plate string  `json:"plate"`
}

type VehicleResult struct {
	Score float64 `json:"score"`
}

type PlaterRecognizer interface {
	RecognizePlate(r io.Reader) (*PlateResult, error)
}

type VehicleRecognizer interface {
	RecognizeVehicle(r io.Reader) (*VehicleResult, error)
}
