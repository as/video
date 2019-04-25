package atom

type Profile struct {
	HDR HDR
}
type Clip struct {
	HDR    HDR
	Region Atom // Clipping Region def. is confusing
}
type Matte struct {
	Atom
}
type Edit struct {
	Atom
}
type Txas struct {
	Atom
}
type Load struct {
	Atom
}
type Imap struct {
	Atom
}
