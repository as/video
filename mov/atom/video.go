package atom

type Video struct {
	HDR      HDR
	VideoHDR VideoHDR
	Handler  *Handler
	DataInfo Atom
	Sample   Atom
}
type VideoHDR struct {
	HDR          HDR
	Version      byte
	Flags        [3]byte
	GraphicsMode uint16
	OpColor      [6]byte
}
