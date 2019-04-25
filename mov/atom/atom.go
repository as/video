package atom

type Atom struct {
	Data interface{}
}

func (a *Atom) Header() HDR {
	return HDR{}
}

type (
	Type     [4]byte
	Flag     [3]byte
	Time     int32
	Duration int32
	Range    struct {
		Time     Time
		Duration Duration
	}
	Matrix [36]byte
)

type AtomV struct {
	HDR
	VerFlag
	Data []byte
}
type HDR struct {
	Size int32
	Type Type
}
type VerFlag struct {
	V byte
	F [3]byte
}

type String struct {
	HDR  HDR
	Data []byte
}

type PTV struct {
	Size     int16
	_        int32
	Slide    byte
	Autoplay byte
}

type Color struct {
	HDR   HDR
	Seed  uint16
	Flag  uint16
	Size  uint16
	Array []byte
}

type User struct {
	HDR  HDR
	List []Atom
}
