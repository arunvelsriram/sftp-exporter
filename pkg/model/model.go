package model

type FSStat struct {
	Path       string
	FreeSpace  float64
	TotalSpace float64
}

type FSStats []FSStat

type ObjectStat struct {
	Path        string
	ObjectCount float64
	ObjectSize  float64
}

type ObjectStats []ObjectStat
