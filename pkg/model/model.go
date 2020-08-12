package model

type FSStat struct {
	Path       string
	FreeSpace  float64
	TotalSpace float64
}

type FSStats []FSStat
