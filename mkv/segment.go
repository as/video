package main

type Segment struct {
	ChapterTranslateCodec      string
	ChapterTranslateID         []byte
	TimecodeScale              int
	MuxingApp                  string
	WritingApp                 string
	SegmentUID                 []byte
	SegmentFilename            string
	PrevUID                    []byte
	PrevFilename               string
	NextUID                    []byte
	NextFilename               string
	Duration                   []byte
	DateUTC                    int
	Title                      string
	ChapterTranslateEditionUID string
	SegmentFamily              []byte
	ChapterTranslate           []byte
}
