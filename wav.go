/*

Package wav is a package allowing developers to decode and encode audio PCM
data using the Waveform Audio File Format https://en.wikipedia.org/wiki/WAV

*/
package wav

import "errors"

var (
	// ErrPCMChunkNotFound indicates a bad audio file without data
	ErrPCMChunkNotFound = errors.New("PCM Chunk not found in audio file")
)

func nullTermStr(b []byte) string {
	return string(b[:clen(b)])
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
