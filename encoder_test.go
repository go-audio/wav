package wav_test

import (
	"os"
	"testing"

	"github.com/go-audio/wav"
)

func TestEncoderRoundTrip(t *testing.T) {
	os.Mkdir("testOutput", 0777)
	testCases := []struct {
		in   string
		out  string
		desc string
	}{
		{"fixtures/kick.wav", "testOutput/kick.wav", "22050 Hz @ 16 bits, 1 channel(s), 44100 avg bytes/sec, duration: 204.172335ms"},
		{"fixtures/kick-16b441k.wav", "testOutput/kick-16b441k.wav", "2 ch,  44100 Hz, 'lpcm' 16-bit little-endian signed integer"},
		{"fixtures/bass.wav", "testOutput/bass.wav", "44100 Hz @ 24 bits, 2 channel(s), 264600 avg bytes/sec, duration: 543.378684ms"},
	}

	for i, tc := range testCases {
		t.Logf("%d - in: %s, out: %s\n%s", i, tc.in, tc.out, tc.desc)
		in, err := os.Open(tc.in)
		if err != nil {
			t.Fatalf("couldn't open %s %v", tc.in, err)
		}
		d := wav.NewDecoder(in)
		buf, err := d.FullPCMBuffer()
		if err != nil {
			t.Fatalf("couldn't read buffer %s %v", tc.in, err)
		}

		// pcm := d.PCM()
		// numChannels, bitDepth, sampleRate, err := pcm.Info()
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// totalFrames := pcm.Size()
		// frames, err := d.FramesInt()
		// if err != nil {
		// 	t.Fatal(err)
		// }
		in.Close()
		// t.Logf("%s - total frames %d - total samples %d", tc.in, totalFrames, len(frames))

		out, err := os.Create(tc.out)
		if err != nil {
			t.Fatalf("couldn't create %s %v", tc.out, err)
		}

		e := wav.NewEncoder(out,
			buf.Format.SampleRate,
			buf.Format.BitDepth,
			buf.Format.NumChannels,
			int(d.WavAudioFormat))
		if err := e.Write(buf); err != nil {
			t.Fatal(err)
		}
		if err := e.Close(); err != nil {
			t.Fatal(err)
		}
		out.Close()

		nf, err := os.Open(tc.out)
		if err != nil {
			t.Fatal(err)
		}

		nd := wav.NewDecoder(nf)
		nBuf, err := nd.FullPCMBuffer()
		if err != nil {
			t.Fatalf("couldn't extract the PCM from %s - %v", nf.Name(), err)
		}

		nf.Close()
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := os.Remove(nf.Name()); err != nil {
				panic(err)
			}
		}()

		if nBuf.Format.SampleRate != buf.Format.SampleRate {
			t.Fatalf("sample rate didn't support roundtripping exp: %d, got: %d", buf.Format.SampleRate, nBuf.Format.SampleRate)
		}
		if nBuf.Format.BitDepth != buf.Format.BitDepth {
			t.Fatalf("sample size didn't support roundtripping exp: %d, got: %d", buf.Format.BitDepth, nBuf.Format.BitDepth)
		}
		if nBuf.Format.NumChannels != buf.Format.NumChannels {
			t.Fatalf("the number of channels didn't support roundtripping exp: %d, got: %d", buf.Format.NumChannels, nBuf.Format.NumChannels)
		}
		if len(nBuf.Data) != len(buf.Data) {
			t.Fatalf("the reported number of frames didn't support roundtripping, exp: %d, got: %d", len(buf.Data), len(nBuf.Data))
		}
		for i := 0; i < len(buf.Data); i++ {
			if buf.Data[i] != nBuf.Data[i] {
				max := len(buf.Data)
				if i+3 < max {
					max = i + 3
				}
				t.Fatalf("frame value at position %d: %d -> %d\ndidn't match new buffer position %d: %d -> %d", i, buf.Data[:i+1], buf.Data[i+1:max], i, nBuf.Data[:i+1], nBuf.Data[i+1:max])
			}
		}

	}
}
