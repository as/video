package atom

type Media struct {
	HDR HDR
	MediaHDR
	ExtLang Atom
	Handler Atom
	Info    Atom
	User    *User
}
type MediaHDR struct {
	HDR       HDR
	Ver       byte
	Flag      Flag
	CTime     Time
	MTime     Time
	TimeScale int32
	Duration  Duration
	Lang      uint16
	Quality   uint16
}
type Base struct {
	HDR     // minf
	BaseHDR struct {
		HDR           // gmhd
		Base BaseInfo // gmin
		Text Atom     // text
	}
}
type BaseInfo struct {
	VideoHDR
	Balance uint16
	_       uint16
}
type Handler struct {
	HDR  HDR
	Ver  byte
	Flag Flag
	Com
}
type Com struct {
	ComHdr
	Name []byte
}
type ComHdr struct {
	Type         uint32
	SubType      uint32
	Manufacturer uint32
	Flag         uint32
	Mask         uint32
}
