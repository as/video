package atom

type Sound struct {
	HDR      HDR
	SoundHDR SoundHDR
	Handler  *Handler
	DataInfo Atom
	Sample   Atom
}
type SoundHDR struct {
	HDR     HDR
	Version byte
	Flags   [3]byte
	Balance uint16
	_       uint16
}
