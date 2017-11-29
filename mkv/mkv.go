// usage (before this is seperated into it's own package)
// go run mkv.go file.mkv
//
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/bits"
	"os"
	"strings"
)

var file = os.Args[1]

func main() {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	tags := make(map[string]string)

	key := ""
	r := bytes.NewReader(data)

	readstring := func(len int64) string {
		m := make([]byte, len)
		r.Read(m)
		return strings.Trim(string(m), "\x00")
	}

	// Scan through the MKV, looking for elements we want
	// including their predecessors.
Scan:
	for i := 0; i < 1e7; i++ {
		e, a, err := NextItem(r)
		nm := name(e)
		// uncomment for lots of information
		// fmt.Printf("%s:#%d,+#%d	%q\n", file, at(r), a, nm)
		if err != nil {
			break
		}

		switch nm {
		case "SimpleTag", "Tags", "Tag":
			// nothing, just descend without advancing
		case "TagName":
			key = readstring(a)
		case "TagString":
			tags[key] = readstring(a)
		case "CodecID":
			tags["CodecID"] = readstring(a)
		default:
			_, err = r.Seek(a, io.SeekCurrent)
			if err != nil {
				if err != nil && err != io.EOF {
					log.Fatalln(err)
				}
				break Scan
			}
		}
	}
	for k, v := range tags {
		fmt.Printf("%q=%q\n", k, v)
	}
}

func extract(r *bytes.Reader) (nz int, advance int64, v uint32) {
	err := binary.Read(r, binary.BigEndian, &v)
	if err != nil {
		return -1, -1, 0
	}
	nz = int(bits.LeadingZeros32(v))
	r.Seek(-4, io.SeekCurrent)
	advance = int64(nz + 1)
	return int(nz), advance, v
}

func at(r *bytes.Reader) int64 {
	n, _ := r.Seek(0, io.SeekCurrent)
	return int64(n)
}

func NextItem(r *bytes.Reader) (elem int, advance int64, err error) {
	_, advance, v := extract(r)
	if advance == -1 {
		return 0, 0, io.EOF
	}
	_, err = r.Seek(advance, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}

	elem = int(v)
	elem = elem >> ((4 - uint(advance)) * 8)

	_, advance, v = extract(r)
	p0 := at(r)
	p1, err := r.Seek(advance, io.SeekCurrent)
	L := uint(advance)
	v = v << L >> L
	v >>= (4 - L) * 8

	if p0 == p1 {
		err = io.EOF
	}
	return elem, int64(v), err
}

func name(e int) string {
	s, ok := ElementName[e]
	if !ok {
		return fmt.Sprintf("%x", e)
	}
	return s
}

var ElementName = map[int]string{
	0xBF:       "CRC-32",
	0x1A45DFA3: "EBML",
	0x4286:     "EBMLVersion",
	0x42F7:     "EBMLReadVersion",
	0x42F2:     "EBMLMaxIDLength",
	0x42F3:     "EBMLMaxSizeLength",
	0x4282:     "DocType",
	0x4287:     "DocTypeVersion",
	0x4285:     "DocTypeReadVersion",
	0xEC:       "Void",
	0x4DBB:     "Seek",
	0x53AB:     "SeekID",
	0x53AC:     "SeekPosition",
	0x18538067: "Segment",
	0x6532:     "SignedElement",
	0x7384:     "SegmentFilename",
	0x4444:     "SegmentFamily",
	0x6924:     "ChapterTranslate",
	0x5854:     "SilentTracks",
	0x7446:     "AttachmentLink",
	0x114D9B74: "SeekHead",
	0x6624:     "TrackTranslate",
	0x6240:     "ContentEncoding",
	0x5031:     "ContentEncodingOrder",
	0x5032:     "ContentEncodingScope",
	0x5033:     "ContentEncodingType",
	0x5034:     "ContentCompression",
	0x4254:     "ContentCompAlgo",
	0x4255:     "ContentCompSettings",
	0x5035:     "ContentEncryption",
	0x96:       "CueRefTime",
	0x97:       "CueRefCluster",
	0x4660:     "FileMimeType",
	0x4675:     "FileReferral",
	0x4661:     "FileUsedStartTime",
	0x4662:     "FileUsedEndTime",
	0x45DD:     "EditionFlagOrdered",
	0x92:       "ChapterTimeEnd",
	0x98:       "ChapterFlagHidden",
	0x4598:     "ChapterFlagEnabled",
	0x6E67:     "ChapterSegmentUID",
	0x6EBC:     "ChapterSegmentEditionUID",
	0x63C3:     "ChapterPhysicalEquiv",
	0x8F:       "ChapterTrack",
	0x89:       "ChapterTrackNumber",
	0x437E:     "ChapCountry",
	0x6944:     "ChapProcess",
	0x450D:     "ChapProcessPrivate",
	0x6911:     "ChapProcessCommand",
	0x6922:     "ChapProcessTime",
	0x6933:     "ChapProcessData",
	0x1254C367: "Tags",
	0x7373:     "Tag",
	0x63C0:     "Targets",
	0x68CA:     "TargetTypeValue",
	0x63CA:     "TargetType",
	0x63C5:     "TagTrackUID",
	0x63C9:     "TagEditionUID",
	0x63C4:     "TagChapterUID",
	0x63C6:     "TagAttachmentUID",
	0x67C8:     "SimpleTag",
	0x45A3:     "TagName",
	0x447A:     "TagLanguage",
	0x4484:     "TagDefault",
	0x4487:     "TagString",
	0x1B538667: "SignatureSlot",
	0x7E8A:     "SignatureAlgo",
	0x7E9A:     "SignatureHash",
	0x7EA5:     "SignaturePublicKey",
	0x7EB5:     "Signature",
	0x7E5B:     "SignatureElements",
	0x7E7B:     "SignatureElementList",
	0x1549A966: "Info",
	0x73A4:     "SegmentUID",
	0x3CB923:   "PrevUID",
	0x3C83AB:   "PrevFilename",
	0x3EB923:   "NextUID",
	0x3E83BB:   "NextFilename",
	0x69FC:     "ChapterTranslateEditionUID",
	0x69BF:     "ChapterTranslateCodec",
	0x69A5:     "ChapterTranslateID",
	0x2AD7B1:   "TimecodeScale",
	0x4489:     "Duration",
	0x4461:     "DateUTC",
	0x7BA9:     "Title",
	0x4D80:     "MuxingApp",
	0x5741:     "WritingApp",
	0x1F43B675: "Cluster",
	0xE7:       "Timecode",
	0x58D7:     "SilentTrackNumber",
	0xA7:       "Position",
	0xAB:       "PrevSize",
	0xA3:       "SimpleBlock",
	0xA0:       "BlockGroup",
	0xA1:       "Block",
	0xA2:       "BlockVirtual",
	0x75A1:     "BlockAdditions",
	0xA6:       "BlockMore",
	0xEE:       "BlockAddID",
	0xA5:       "BlockAdditional",
	0x9B:       "BlockDuration",
	0xFA:       "ReferencePriority",
	0xFB:       "ReferenceBlock",
	0xFD:       "ReferenceVirtual",
	0xA4:       "CodecState",
	0x75A2:     "DiscardPadding",
	0x8E:       "Slices",
	0xE8:       "TimeSlice",
	0xCC:       "LaceNumber",
	0xCD:       "FrameNumber",
	0xCB:       "BlockAdditionID",
	0xCE:       "Delay",
	0xCF:       "SliceDuration",
	0xC8:       "ReferenceFrame",
	0xC9:       "ReferenceOffset",
	0xCA:       "ReferenceTimeCode",
	0xAF:       "EncryptedBlock",
	0x1654AE6B: "Tracks",
	0xAE:       "TrackEntry",
	0xD7:       "TrackNumber",
	0x73C5:     "TrackUID",
	0x83:       "TrackType",
	0xB9:       "FlagEnabled",
	0x88:       "FlagDefault",
	0x55AA:     "FlagForced",
	0x9C:       "FlagLacing",
	0x6DE7:     "MinCache",
	0x6DF8:     "MaxCache",
	0x23E383:   "DefaultDuration",
	0x234E7A:   "DefaultDecodedFieldDuration",
	0x23314F:   "TrackTimecodeScale",
	0x537F:     "TrackOffset",
	0x55EE:     "MaxBlockAdditionID",
	0x536E:     "Name",
	0x22B59C:   "Language",
	0x86:       "CodecID",
	0x63A2:     "CodecPrivate",
	0x258688:   "CodecName",
	0x3A9697:   "CodecSettings",
	0x3B4040:   "CodecInfoURL",
	0x26B240:   "CodecDownloadURL",
	0xAA:       "CodecDecodeAll",
	0x6FAB:     "TrackOverlay",
	0x56AA:     "CodecDelay",
	0x56BB:     "SeekPreRoll",
	0x66FC:     "TrackTranslateEditionUID",
	0x66BF:     "TrackTranslateCodec",
	0x66A5:     "TrackTranslateTrackID",
	0xE0:       "Video",
	0x9A:       "FlagInterlaced",
	0x9D:       "FieldOrder",
	0x53B8:     "StereoMode",
	0x53C0:     "AlphaMode",
	0x53B9:     "OldStereoMode",
	0xB0:       "PixelWidth",
	0xBA:       "PixelHeight",
	0x54AA:     "PixelCropBottom",
	0x54BB:     "PixelCropTop",
	0x54CC:     "PixelCropLeft",
	0x54DD:     "PixelCropRight",
	0x54B0:     "DisplayWidth",
	0x54BA:     "DisplayHeight",
	0x54B2:     "DisplayUnit",
	0x54B3:     "AspectRatioType",
	0x2EB524:   "ColourSpace",
	0x2FB523:   "GammaValue",
	0x2383E3:   "FrameRate",
	0x55B0:     "Colour",
	0x55B1:     "MatrixCoefficients",
	0x55B2:     "BitsPerChannel",
	0x55B3:     "ChromaSubsamplingHorz",
	0x55B4:     "ChromaSubsamplingVert",
	0x55B5:     "CbSubsamplingHorz",
	0x55B6:     "CbSubsamplingVert",
	0x55B7:     "ChromaSitingHorz",
	0x55B8:     "ChromaSitingVert",
	0x55B9:     "Range",
	0x55BA:     "TransferCharacteristics",
	0x55BB:     "Primaries",
	0x55BC:     "MaxCLL",
	0x55BD:     "MaxFALL",
	0x55D0:     "MasteringMetadata",
	0x55D1:     "PrimaryRChromaticityX",
	0x55D2:     "PrimaryRChromaticityY",
	0x55D3:     "PrimaryGChromaticityX",
	0x55D4:     "PrimaryGChromaticityY",
	0x55D5:     "PrimaryBChromaticityX",
	0x55D6:     "PrimaryBChromaticityY",
	0x55D7:     "WhitePointChromaticityX",
	0x55D8:     "WhitePointChromaticityY",
	0x55D9:     "LuminanceMax",
	0x55DA:     "LuminanceMin",
	0xE1:       "Audio",
	0xB5:       "SamplingFrequency",
	0x78B5:     "OutputSamplingFrequency",
	0x9F:       "Channels",
	0x7D7B:     "ChannelPositions",
	0x6264:     "BitDepth",
	0xE2:       "TrackOperation",
	0xE3:       "TrackCombinePlanes",
	0xE4:       "TrackPlane",
	0xE5:       "TrackPlaneUID",
	0xE6:       "TrackPlaneType",
	0xE9:       "TrackJoinBlocks",
	0xED:       "TrackJoinUID",
	0xC0:       "TrickTrackUID",
	0xC1:       "TrickTrackSegmentUID",
	0xC6:       "TrickTrackFlag",
	0xC7:       "TrickMasterTrackUID",
	0xC4:       "TrickMasterTrackSegmentUID",
	0x6D80:     "ContentEncodings",
	0x47E1:     "ContentEncAlgo",
	0x47E2:     "ContentEncKeyID",
	0x47E3:     "ContentSignature",
	0x47E4:     "ContentSigKeyID",
	0x47E5:     "ContentSigAlgo",
	0x47E6:     "ContentSigHashAlgo",
	0x1C53BB6B: "Cues",
	0xBB:       "CuePoint",
	0xB3:       "CueTime",
	0xB7:       "CueTrackPositions",
	0xF7:       "CueTrack",
	0xF1:       "CueClusterPosition",
	0xF0:       "CueRelativePosition",
	0xB2:       "CueDuration",
	0x5378:     "CueBlockNumber",
	0xEA:       "CueCodecState",
	0xDB:       "CueReference",
	0x535F:     "CueRefNumber",
	0xEB:       "CueRefCodecState",
	0x1941A469: "Attachments",
	0x61A7:     "AttachedFile",
	0x467E:     "FileDescription",
	0x466E:     "FileName",
	0x465C:     "FileData",
	0x46AE:     "FileUID",
	0x1043A770: "Chapters",
	0x45B9:     "EditionEntry",
	0x45BC:     "EditionUID",
	0x45BD:     "EditionFlagHidden",
	0x45DB:     "EditionFlagDefault",
	0xB6:       "ChapterAtom",
	0x73C4:     "ChapterUID",
	0x5654:     "ChapterStringUID",
	0x91:       "ChapterTimeStart",
	0x80:       "ChapterDisplay",
	0x85:       "ChapString",
	0x437C:     "ChapLanguage",
	0x6955:     "ChapProcessCodecID",
	0x4485:     "TagBinary",
}
