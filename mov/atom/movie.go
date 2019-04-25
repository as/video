package atom

type Movie struct {
	//	HDR       HDR
	Ver       byte
	Flag      Flag
	Create    Time
	Mod       Time
	TimeScale int32
	Duration
	Rate      uint32
	Volume    uint16
	_         [10]byte
	Matrix    Matrix
	Preview   Range
	Poster    Time
	Selection Range
	Current   Time
	NextID    int32
}
