package model

type ObjectStat struct {
	Path        string
	ObjectCount float64
	ObjectSize  float64
}

type ObjectStats []ObjectStat
