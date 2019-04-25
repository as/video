package atom

type TrackName struct {
	HDR  HDR
	_    uint32
	Lang uint16
	Name []byte
}

type TrackHDR struct {
	HDR      HDR
	Ver      byte
	Flag     Flag
	Create   Time
	Mod      Time
	ID       int32
	_        [4]byte
	Duration Duration
	_        [8]byte
	Layer    uint16
	AltGroup uint16
	Volume   uint16
	_        [2]byte
	Matrix   Matrix
	Dx       int32
	Dy       int32
}

type Track struct {
	HDR       HDR
	Profile   *Profile
	THDR      HDR
	Clip      *Clip
	Matte     *Matte
	Edit      *Edit
	Reference *Reference
	Txas      *Txas
	Load      *Load
	Imap      *Imap
	Media     Media
	User      *User
}

type Aperture struct {
	HDR
	Clean  Dim
	Prod   Dim
	Pixels Dim
}
type Dim struct {
	HDR
	Ver  byte
	Flag Flag
	Dx   uint32
	Dy   uint32
}
