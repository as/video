package mkv

import (
	"bytes"
	"testing"
)

func TestPygmy(t *testing.T) {
	s := NewScanner(bytes.NewReader([]byte(mkv)))
	for i := 0; ; i++ {
		e, a, err := s.Next()
		if have, want := name(e), pygtab[i].name; have != want {
			t.Logf("element: have: %q want: %q", have, want)
			t.Fail()
		}
		if have, want := a, pygtab[i].advance; have != want {
			t.Logf("advance: have: %d want: %d", have, want)
			t.Fail()
		}
		if err != nil {
			return
		}
		if name(e) == "Void" {
			s.ParseVoid()
		} else {
			s.Advance(a)
		}
	}
}

var pygtab = []struct {
	name    string
	advance int64
	err     error
}{
	{name: "EBML", advance: 0},
	{name: "EBMLVersion", advance: 1},
	{name: "EBMLReadVersion", advance: 1},
	{name: "EBMLMaxIDLength", advance: 1},
	{name: "EBMLMaxSizeLength", advance: 1},
	{name: "DocType", advance: 8},
	{name: "DocTypeVersion", advance: 1},
	{name: "DocTypeReadVersion", advance: 1},
	{name: "Segment", advance: 0},
	{name: "SeekHead", advance: 65},
	{name: "Void", advance: 0},
	{name: "Info", advance: 0},
	{name: "CRC-32", advance: 4},
	{name: "TimecodeScale", advance: 3},
	{name: "MuxingApp", advance: 13},
	{name: "WritingApp", advance: 13},
	{name: "SegmentUID", advance: 16},
	{name: "Duration", advance: 8},
	{name: "Tracks", advance: 0},
	{name: "CRC-32", advance: 4},
	{name: "TrackEntry", advance: 0},
	{name: "TrackNumber", advance: 1},
	{name: "TrackUID", advance: 1},
	{name: "FlagLacing", advance: 1},
	{name: "Language", advance: 3},
	{name: "CodecID", advance: 15},
	{name: "TrackType", advance: 1},
	{name: "DefaultDuration", advance: 4},
	{name: "Video", advance: 0},
	{name: "PixelWidth", advance: 1},
	{name: "PixelHeight", advance: 1},
	{name: "FlagInterlaced", advance: 1},
	{name: "DisplayUnit", advance: 1},
	{name: "CodecPrivate", advance: 39},
	{name: "Tags", advance: 0},
	{name: "CRC-32", advance: 4},
	{name: "Tag", advance: 0},
	{name: "Targets", advance: 0},
	{name: "SimpleTag", advance: 0},
	{name: "TagName", advance: 7},
	{name: "TagString", advance: 13},
	{name: "Tag", advance: 0},
	{name: "Targets", advance: 0},
	{name: "TagTrackUID", advance: 1},
	{name: "SimpleTag", advance: 0},
	{name: "TagName", advance: 7},
	{name: "TagString", advance: 21},
	{name: "Tag", advance: 0},
	{name: "Targets", advance: 0},
	{name: "TagTrackUID", advance: 1},
	{name: "SimpleTag", advance: 0},
	{name: "TagName", advance: 8},
	{name: "TagString", advance: 20},
	{name: "Cluster", advance: 0},
	{name: "CRC-32", advance: 4},
	{name: "Timecode", advance: 1},
	{name: "SimpleBlock", advance: 713},
	{name: "Cues", advance: 0},
	{name: "CRC-32", advance: 4},
	{name: "CuePoint", advance: 15},
	{name: "0", advance: 0},
}

var mkv = "\x1aEߣ\x01\x00\x00\x00\x00\x00\x00#B\x86\x81\x01B\xf7\x81\x01B\xf2\x81\x04B\xf3\x81\bB\x82\x88matroskaB\x87\x81\x04B\x85\x81\x02\x18S\x80g\x01\x00\x00\x00\x00\x00\x05\x9a\x11M\x9bt@A\xbf\x84\xd8G\xd51M\xbb\x8bS\xab\x84\x15I\xa9fS\xac\x81\xe5M\xbb\x8cS\xab\x84\x16T\xaekS\xac\x82\x01<M\xbb\x8cS\xab\x84\x12T\xc3gS\xac\x82\x01\xc4M\xbb\x8cS\xab\x84\x1cS\xbbkS\xac\x82\x05w\xec\x01\x00\x00\x00\x00\x00\x00\x95\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x15I\xa9f\x01\x00\x00\x00\x00\x00\x00K\xbf\x84\f\xa3\x9cD*ױ\x83\x0fB@M\x80\x8dLavf57.65.100WA\x8dLavf57.65.100s\xa4\x90\xb9\xe1<\x16\xab\xf9c\xaatk\xdb\x14\x8b\xe3\x9d*D\x89\x88@@\x80\x00\x00\x00\x00\x00\x16T\xaek\x01\x00\x00\x00\x00\x00\x00|\xbf\x84\b=\xd1L\xae\x01\x00\x00\x00\x00\x00\x00mׁ\x01sŁ\x01\x9c\x81\x00\"\xb5\x9c\x83und\x86\x8fV_MPEG4/ISO/AVC\x83\x81\x01#ツ\x01\xfc\xa04\xe0\x01\x00\x00\x00\x00\x00\x00\r\xb0\x81\x02\xba\x81\x02\x9a\x81\x02T\xb2\x81\x04c\xa2\xa7\x01z\x00\x0a\xff\xe1\x00\x16gz\x00\x0a\xbc\xd9_\x88\x8f\x84\x00\x14XT\x04Ĵ\x00<H\x96X\x01\x00\x06h\xeb\xe3\xcb\"\xc0\x12T\xc3g\x01\x00\x00\x00\x00\x00\x00ƿ\x84i\xf7\xcdWss\x01\x00\x00\x00\x00\x00\x00.c\xc0\x01\x00\x00\x00\x00\x00\x00\x00g\xc8\x01\x00\x00\x00\x00\x00\x00\x1aE\xa3\x87ENCODERD\x87\x8dLavf57.65.100ss\x01\x00\x00\x00\x00\x00\x00:c\xc0\x01\x00\x00\x00\x00\x00\x00\x04cŁ\x01g\xc8\x01\x00\x00\x00\x00\x00\x00\"E\xa3\x87ENCODERD\x87\x95Lavc57.75.100 libx264ss\x01\x00\x00\x00\x00\x00\x00:c\xc0\x01\x00\x00\x00\x00\x00\x00\x04cŁ\x01g\xc8\x01\x00\x00\x00\x00\x00\x00\"E\xa3\x88DURATIOND\x87\x9400:00:00.033000000\x00\x00\x1fC\xb6u\x01\x00\x00\x00\x00\x00\x02տ\x84NϘ\xdc\xe7\x81\x00\xa3BɁ\x00\x00\x80\x00\x00\x02\xae\x06\x05\xff\xff\xaa\xdcE\xe9\xbd\xe6\xd9H\xb7\x96,\xd8 \xd9#\xee\xefx264 - core 148 r2744 b97ae06 - H.264/MPEG-4 AVC codec - Copyleft 2003-2016 - http://www.videolan.org/x264.html - options: cabac=1 ref=3 deblock=1:0:0 analyse=0x3:0x113 me=hex subme=7 psy=1 psy_rd=1.00:0.00 mixed_ref=1 me_range=16 chroma_me=1 trellis=1 8x8dct=1 cqm=0 deadzone=21,11 fast_pskip=1 chroma_qp_offset=-2 threads=1 lookahead_threads=1 sliced_threads=0 nr=0 decimate=1 interlaced=0 bluray_compat=0 constrained_intra=0 bframes=3 b_pyramid=2 b_adapt=1 b_bias=0 direct=1 weightb=1 open_gop=0 weightp=2 keyint=250 keyint_min=25 scenecut=40 intra_refresh=0 rc_lookahead=40 rc=crf mbtree=1 crf=23.0 qcomp=0.60 qpmin=0 qpmax=69 qpstep=4 ip_ratio=1.40 aq=1:1.00\x00\x80\x00\x00\x00\x0fe\x88\x84\x00o\xfe\xbdg\xe6Y=\b\xdeS\xc1\x1cS\xbbk\x01\x00\x00\x00\x00\x00\x00\x17\xbf\x84[o\x1b经\xb3\x81\x00\xb7\x8a\xf7\x81\x01\xf1\x82\x02\x96\xf0\x81\t"
