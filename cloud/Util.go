package cloud

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// timestamp_microseconds_instanceName_regionCoordinates_numberOfChanges_token
// 1564859471_6-474162_oprit_577-283-727-375_1153_27.mp4
// - Timestamp
// - Size + - + microseconds
// - device
// - Region
// - Number of changes
// - Token
// startRecording = time.Now().Unix() // we mark the current time when the record started.ss
//
//	s := strconv.FormatInt(startRecording, 10) + "_" +
//		"6" + "-" +
//		"967003" + "_" +
//		config.Name + "_" +
//		"200-200-400-400" + "_0_" +
//		"769"
//
// filename formatL %Y-%m-%d_%H-%M-%S.%f
func ToKerberosFormat(streamKey string, filename string) (string, error) {
	ext := filepath.Ext(filename)
	fileNoExt := strings.TrimSuffix(filename, ext)

	layout := "2006-01-02_15-04-05.000000"
	t, err := time.Parse(layout, fileNoExt)
	if err != nil {
		// try with pattern: 2006-01-02_15-04-05-000000
		if idx := strings.LastIndex(fileNoExt, "-"); idx != -1 {
			chars := []rune(fileNoExt)
			chars[idx] = '.'
			t, err = time.Parse(layout, string(chars))
		}
	}

	if err == nil {
		startRecording := t.Unix()
		s := strconv.FormatInt(startRecording, 10) + "_" +
			"6" + "-" +
			"967003" + "_" +
			streamKey + "_" +
			"200-200-400-400" + "_0_" +
			"769" +
			ext

		return s, err
	}

	return filename, err
}
