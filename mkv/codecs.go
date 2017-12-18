package mkv

type CodecInfo struct {
	ID          string
	Name        string
	Kind        string
	Description string
	URL         string
	Download    string
	Settings    string
}

func init() {
	for k, v := range codectab {
		v.ID = k
		if len(k) > 0 {
			switch k[0] {
			case 'A':
				v.Kind = "Audio"
			case 'V':
				v.Kind = "Video"
			case 'S':
				v.Kind = "Subtitle"
			case 'B':
				v.Kind = "Button"
			}
		}
		codectab[k] = v
	}
}

var codectab = map[string]CodecInfo{
	"A_AAC/MPEG2/LC":     {},
	"A_AAC/MPEG2/LC/SBR": {},
	"A_AAC/MPEG2/MAIN":   {},
	"A_AAC/MPEG2/SSR":    {},
	"A_AAC/MPEG4/LC":     {},
	"A_AAC/MPEG4/LC/SBR": {},
	"A_AAC/MPEG4/LTP":    {},
	"A_AAC/MPEG4/MAIN":   {},
	"A_AAC/MPEG4/SSR":    {},
	"A_AC3":              {},
	"A_AC3/BSID10":       {},
	"A_AC3/BSID9":        {},
	"A_ALAC":             {},
	"A_DTS":              {},
	"A_DTS/EXPRESS":      {},
	"A_DTS/LOSSLESS":     {},
	"A_FLAC":             {},
	"A_MPC":              {},
	"A_MPEG/L1":          {},
	"A_MPEG/L2":          {},
	"A_MPEG/L3":          {},
	"A_MS/ACM":           {},
	"A_PCM/FLOAT/IEEE":   {},
	"A_PCM/INT/BIG":      {},
	"A_PCM/INT/LIT":      {},
	"A_QUICKTIME":        {},
	"A_QUICKTIME/QDM2":   {},
	"A_REAL/14_4":        {},
	"A_REAL/28_8":        {},
	"A_REAL/ATRC":        {},
	"A_REAL/COOK":        {},
	"A_REAL/RALF":        {},
	"A_REAL/SIPR":        {},
	"A_TTA1":             {},
	"A_VORBIS":           {},
	"A_WAVPACK4":         {},
	"B_VOBBTN":           {},
	"S_IMAGE/BMP":        {},
	"S_KATE":             {},
	"S_TEXT/ASS":         {},
	"S_TEXT/SSA":         {},
	"S_TEXT/USF":         {},
	"S_TEXT/UTF8":        {},
	"S_TEXT/WEBVTT":      {},
	"S_VOBSUB":           {},
	"V_3IVX":             {},
	"V_COREYUV":          {},
	"V_DV":               {},
	"V_HUFFYUV":          {},
	"V_INDEO5":           {},
	"V_MJPEG":            {},
	"V_MJPEG2000":        {},
	"V_MJPEG2000LL":      {},
	"V_MPEG1":            {},
	"V_MPEG2":            {},
	"V_MPEG4/ISO/???":    {},
	"V_MPEG4/ISO/AP":     {},
	"V_MPEG4/ISO/ASP":    {},
	"V_MPEG4/ISO/SP":     {},
	"V_MPEG4/MS/V3":      {},
	"V_MS/VFW/FOURCC":    {},
	"V_MSWMV":            {},
	"V_ON2VP4":           {},
	"V_ON2VP5":           {},
	"V_PRORES":           {},
	"V_QUICKTIME":        {},
	"V_REAL/RV10":        {},
	"V_REAL/RV20":        {},
	"V_REAL/RV30":        {},
	"V_REAL/RV40":        {},
	"V_RUDUDUV_TARKIN":   {},
	"V_THEORA":           {},
	"V_UNCOMPRESSED":     {},
}
