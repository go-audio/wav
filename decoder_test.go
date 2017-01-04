package wav_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func TestDecoder_Duration(t *testing.T) {
	testCases := []struct {
		in       string
		duration time.Duration
	}{
		{"fixtures/kick.wav", time.Duration(204172335 * time.Nanosecond)},
	}

	for _, tc := range testCases {
		f, err := os.Open(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		dur, err := wav.NewDecoder(f).Duration()
		if err != nil {
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		if dur != tc.duration {
			t.Fatalf("expected duration to be: %s but was %s", tc.duration, dur)
		}
	}

}

func TestDecoder_IsValidFile(t *testing.T) {
	testCases := []struct {
		in      string
		isValid bool
	}{
		{"fixtures/kick.wav", true},
		{"fixtures/bass.wav", true},
		{"fixtures/dirty-kick-24b441k.wav", true},
		{"fixtures/sample.avi", false},
		{"fixtures/bloop.aif", false},
	}

	for _, tc := range testCases {
		f, err := os.Open(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		d := wav.NewDecoder(f)
		if d.IsValidFile() != tc.isValid {
			t.Fatalf("validation of the wav files doesn't match expected %t, got %t", tc.isValid, d.IsValidFile())
		}
	}

}

func TestDecoder_Attributes(t *testing.T) {
	testCases := []struct {
		in             string
		numChannels    int
		sampleRate     int
		avgBytesPerSec int
		bitDepth       int
	}{
		{in: "fixtures/kick.wav",
			numChannels:    1,
			sampleRate:     22050,
			avgBytesPerSec: 44100,
			bitDepth:       16,
		},
	}

	for _, tc := range testCases {
		f, err := os.Open(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		d := wav.NewDecoder(f)
		d.ReadInfo()
		f.Close()
		if int(d.NumChans) != tc.numChannels {
			t.Fatalf("expected info to have %d channels but it has %d", tc.numChannels, d.NumChans)
		}
		if int(d.SampleRate) != tc.sampleRate {
			t.Fatalf("expected info to have a sample rate of %d but it has %d", tc.sampleRate, d.SampleRate)
		}
		if int(d.AvgBytesPerSec) != tc.avgBytesPerSec {
			t.Fatalf("expected info to have %d avg bytes per sec but it has %d", tc.avgBytesPerSec, d.AvgBytesPerSec)
		}
		if int(d.BitDepth) != tc.bitDepth {
			t.Fatalf("expected info to have %d bits per sample but it has %d", tc.bitDepth, d.BitDepth)
		}
	}
}

func TestDecoder_PCMBuffer(t *testing.T) {
	testCases := []struct {
		input      string
		desc       string
		samples    []int
		numSamples int
	}{
		{"fixtures/bass.wav",
			"2 ch,  44100 Hz, 'lpcm' 24-bit little-endian signed integer",
			[]int{0, 0, 28160, 26368, 16128, 14848, -746240, -705536, 596480, 565504, 2161408, 2050304, -306944, -271872, -607488, -537856, -1624064, -1477376, -4554752, -4233472, -16599808, -15644160, -21150208, -20034560, -6344192, -6146816, 28578048, 26699520, 60357888, 56626176, 70529280},
			47914,
		},
		{"fixtures/kick-16b441k.wav",
			"2 ch,  44100 Hz, 'lpcm' (0x0000000C) 16-bit little-endian signed integer",
			[]int{0, 0, 0, 0, 0, 0, 3, 3, 28, 28, 130, 130, 436, 436, 1103, 1103, 2140, 2140, 3073, 3073, 2884, 2884, 760, 760, -2755, -2755, -5182, -5182, -3860, -3860, 1048, 1048, 5303, 5303, 3885, 3885, -3378, -3378, -9971, -9971, -8119, -8119, 2616, 2616, 13344, 13344, 13297, 13297, 553, 553, -15013, -15013, -20341, -20341, -10692, -10692, 6553, 6553, 18819, 18819, 18824, 18824, 8617, 8617, -4253, -4253, -13305, -13305, -16289, -16289, -13913, -13913, -7552, -7552, 1334, 1334, 10383, 10383, 16409, 16409, 16928, 16928, 11771, 11771, 3121, 3121, -5908, -5908, -12829, -12829, -16321, -16321, -15990, -15990, -12025, -12025, -5273, -5273, 2732, 2732, 10094, 10094, 15172, 15172, 17038, 17038, 15563, 15563, 11232, 11232, 4973, 4971, -2044, -2044, -8602, -8602, -13659, -13659, -16458, -16458, -16574, -16575, -14012, -14012, -9294, -9294, -3352, -3352, 2823, 2823, 8485, 8485, 13125, 13125, 16228, 16228, 17214, 17214, 15766, 15766, 12188, 12188, 7355, 7355, 2152, 2152, -2973, -2973, -7929, -7929, -12446, -12446, -15806, -15806, -17161, -17161, -16200, -16200, -13407, -13407, -9681, -9681, -5659, -5659, -1418, -1418, 3212, 3212, 8092, 8092, 12567, 12567, 15766, 15766, 17123, 17123, 16665, 16665, 14863, 14863, 12262, 12262, 9171, 9171, 5644, 5644, 1636, 1636, -2768, -2768, -7262, -7262, -11344, -11344, -14486, -14486, -16310, -16310, -16710, -16710, -15861, -15861, -14093, -14093, -11737, -11737, -8974, -8974, -5840, -5840, -2309, -2309, 1577, 1577, 5631, 5631, 9510, 9510, 12821, 12821, 15218, 15218, 16500, 16500, 16663, 16663, 15861, 15861, 14338, 14338, 12322, 12322, 9960, 9960},
			15564,
		},
	}

	for i, tc := range testCases {
		t.Logf("%d - %s\n", i, tc.input)
		path, _ := filepath.Abs(tc.input)
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		d := wav.NewDecoder(f)

		intBuf := make([]int, len(tc.samples))
		buf := &audio.IntBuffer{Data: intBuf}
		err = d.PCMBuffer(buf)
		if err != nil {
			t.Fatal(err)
		}
		if len(buf.Data) != len(tc.samples) {
			t.Fatalf("the length of the buffer (%d) didn't match what we expected (%d)", len(buf.Data), len(tc.samples))
		}
		for i := 0; i < len(buf.Data); i++ {
			if buf.Data[i] != tc.samples[i] {
				t.Fatalf("Expected %d at position %d, but got %d", tc.samples[i], i, buf.Data[i])
			}
		}
	}
}

func TestDecoder_FullPCMBuffer(t *testing.T) {
	testCases := []struct {
		input       string
		desc        string
		samples     []int
		numSamples  int
		totalFrames int
		numChannels int
		sampleRate  int
		bitDepth    int
	}{
		{"fixtures/bass.wav",
			"2 ch,  44100 Hz, 'lpcm' 24-bit little-endian signed integer",
			[]int{0, 0, 28160, 26368, 16128, 14848, -746240, -705536, 596480, 565504, 2161408, 2050304, -306944, -271872, -607488, -537856, -1624064, -1477376, -4554752, -4233472, -16599808, -15644160, -21150208, -20034560, -6344192, -6146816, 28578048, 26699520, 60357888, 56626176, 70529280},
			47914,
			23957,
			2,
			44100,
			24,
		},
		{"fixtures/kick-16b441k.wav",
			"2 ch,  44100 Hz, 'lpcm' (0x0000000C) 16-bit little-endian signed integer",
			[]int{0, 0, 0, 0, 0, 0, 3, 3, 28, 28, 130, 130, 436, 436, 1103, 1103, 2140, 2140, 3073, 3073, 2884, 2884, 760, 760, -2755, -2755, -5182, -5182, -3860, -3860, 1048, 1048, 5303, 5303, 3885, 3885, -3378, -3378, -9971, -9971, -8119, -8119, 2616, 2616, 13344, 13344, 13297, 13297, 553, 553, -15013, -15013, -20341, -20341, -10692, -10692, 6553, 6553, 18819, 18819, 18824, 18824, 8617, 8617, -4253, -4253, -13305, -13305, -16289, -16289, -13913, -13913, -7552, -7552, 1334, 1334, 10383, 10383, 16409, 16409, 16928, 16928, 11771, 11771, 3121, 3121, -5908, -5908, -12829, -12829, -16321, -16321, -15990, -15990, -12025, -12025, -5273, -5273, 2732, 2732, 10094, 10094, 15172, 15172, 17038, 17038, 15563, 15563, 11232, 11232, 4973, 4971, -2044, -2044, -8602, -8602, -13659, -13659, -16458, -16458, -16574, -16575, -14012, -14012, -9294, -9294, -3352, -3352, 2823, 2823, 8485, 8485, 13125, 13125, 16228, 16228, 17214, 17214, 15766, 15766, 12188, 12188, 7355, 7355, 2152, 2152, -2973, -2973, -7929, -7929, -12446, -12446, -15806, -15806, -17161, -17161, -16200, -16200, -13407, -13407, -9681, -9681, -5659, -5659, -1418, -1418, 3212, 3212, 8092, 8092, 12567, 12567, 15766, 15766, 17123, 17123, 16665, 16665, 14863, 14863, 12262, 12262, 9171, 9171, 5644, 5644, 1636, 1636, -2768, -2768, -7262, -7262, -11344, -11344, -14486, -14486, -16310, -16310, -16710, -16710, -15861, -15861, -14093, -14093, -11737, -11737, -8974, -8974, -5840, -5840, -2309, -2309, 1577, 1577, 5631, 5631, 9510, 9510, 12821, 12821, 15218, 15218, 16500, 16500, 16663, 16663, 15861, 15861, 14338, 14338, 12322, 12322, 9960, 9960},
			15564,
			7782,
			2,
			44100,
			16,
		},
		{
			"fixtures/kick.wav",
			"1 ch,  22050 Hz, 'lpcm' 16-bit little-endian signed integer",
			[]int{76, 75, 77, 73, 74, 69, 73, 68, 72, 66, 67, 71, 529, 1427, 2243, 2943, 3512, 3953, 4258, 4436, 4486, 4412, 4220, 3901, 3476, 2937, 2294, 1555, 709, -212, -1231, -2322},
			4484,
			4484,
			1,
			22050,
			16,
		},
	}

	for i, tc := range testCases {
		t.Logf("%d - %s\n", i, tc.input)
		path, _ := filepath.Abs(tc.input)
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		d := wav.NewDecoder(f)
		buf, err := d.FullPCMBuffer()
		if err != nil {
			t.Fatal(err)
		}
		if len(buf.Data) != tc.numSamples {
			t.Fatalf("the length of the buffer (%d) didn't match what we expected (%d)", len(buf.Data), tc.numSamples)
		}
		for i := 0; i < len(tc.samples); i++ {
			if buf.Data[i] != tc.samples[i] {
				t.Fatalf("Expected %d at position %d, but got %d", tc.samples[i], i, buf.Data[i])
			}
		}
		if buf.Format.SampleRate != tc.sampleRate {
			t.Fatalf("expected samplerate to be %d but got %d", tc.sampleRate, buf.Format.SampleRate)
		}
		if buf.Format.NumChannels != tc.numChannels {
			t.Fatalf("expected channel number to be %d but got %d", tc.numChannels, buf.Format.NumChannels)
		}
		framesNbr := buf.NumFrames()
		if framesNbr != tc.totalFrames {
			t.Fatalf("Expected %d frames, got %d\n", tc.totalFrames, framesNbr)
		}
	}
}
