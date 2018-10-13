package wav

import (
	"os"
	"testing"
)

func TestEncoderRoundTrip(t *testing.T) {
	os.Mkdir("testOutput", 0777)
	testCases := []struct {
		in       string
		out      string
		metadata *Metadata
		desc     string
	}{
		{"fixtures/kick.wav", "testOutput/kick.wav", nil, "22050 Hz @ 16 bits, 1 channel(s), 44100 avg bytes/sec, duration: 204.172335ms"},
		{"fixtures/kick-16b441k.wav", "testOutput/kick-16b441k.wav", nil, "2 ch,  44100 Hz, 'lpcm' 16-bit little-endian signed integer"},
		{"fixtures/bass.wav", "testOutput/bass.wav", nil, "44100 Hz @ 24 bits, 2 channel(s), 264600 avg bytes/sec, duration: 543.378684ms"},
		{"fixtures/8bit.wav", "testOutput/8bit.wav", &Metadata{
			Artist: "Matt", Copyright: "copyleft", Comments: "A comment", CreationDate: "2017-12-12", Engineer: "Matt A", Technician: "Matt Aimonetti",
			Genre: "test", Keywords: "go code", Medium: "Virtual", Title: "Titre", Product: "go-audio", Subject: "wav codec",
			Software: "go-audio codec", Source: "Audacity generator", Location: "Los Angeles", TrackNbr: "42",
		}, "1 ch,  44100 Hz, 8-bit unsigned integer"},
		{"fixtures/32bit.wav", "testOutput/32bit.wav", nil, "1 ch, 44100 Hz, 32-bit little-endian signed integer"},
	}

	for i, tc := range testCases {
		t.Logf("%d - in: %s", i, tc.in)
		in, err := os.Open(tc.in)
		if err != nil {
			t.Fatalf("couldn't open %s %v", tc.in, err)
		}
		d := NewDecoder(in)
		buf, err := d.FullPCMBuffer()
		if err != nil {
			t.Fatalf("couldn't read buffer %s %v", tc.in, err)
		}

		in.Close()

		out, err := os.Create(tc.out)
		if err != nil {
			t.Fatalf("couldn't create %s %v", tc.out, err)
		}

		e := NewEncoder(out,
			buf.Format.SampleRate,
			int(d.BitDepth),
			buf.Format.NumChannels,
			int(d.WavAudioFormat))
		if err = e.Write(buf); err != nil {
			t.Fatal(err)
		}
		if tc.metadata != nil {
			e.Metadata = tc.metadata
		}
		if err = e.Close(); err != nil {
			t.Fatal(err)
		}
		out.Close()

		nf, err := os.Open(tc.out)
		if err != nil {
			t.Fatal(err)
		}

		nd := NewDecoder(nf)
		nBuf, err := nd.FullPCMBuffer()
		if err != nil {
			t.Fatalf("couldn't extract the PCM from %s - %v", nf.Name(), err)
		}
		if tc.metadata != nil {
			nd.ReadMetadata()
			if nd.Metadata == nil {
				t.Errorf("expected some metadata, got a nil value")
			}
			if tc.metadata.Artist != nd.Metadata.Artist {
				t.Errorf("expected Artist to be %s, but was %s", tc.metadata.Artist, nd.Metadata.Artist)
			}
			if tc.metadata.Comments != nd.Metadata.Comments {
				t.Errorf("expected Comments to be %s, but was %s", tc.metadata.Comments, nd.Metadata.Comments)
			}
			if tc.metadata.Copyright != nd.Metadata.Copyright {
				t.Errorf("expected Copyright to be %s, but was %s", tc.metadata.Copyright, nd.Metadata.Copyright)
			}
			if tc.metadata.CreationDate != nd.Metadata.CreationDate {
				t.Errorf("expected CreationDate to be %s, but was %s", tc.metadata.CreationDate, nd.Metadata.CreationDate)
			}
			if tc.metadata.Engineer != nd.Metadata.Engineer {
				t.Errorf("expected Engineer to be %s, but was %s", tc.metadata.Engineer, nd.Metadata.Engineer)
			}
			if tc.metadata.Technician != nd.Metadata.Technician {
				t.Errorf("expected Technician to be %s, but was %s", tc.metadata.Technician, nd.Metadata.Technician)
			}
			if tc.metadata.Genre != nd.Metadata.Genre {
				t.Errorf("expected Genre to be %s, but was %s", tc.metadata.Genre, nd.Metadata.Genre)
			}
			if tc.metadata.Keywords != nd.Metadata.Keywords {
				t.Errorf("expected Keywords to be %s, but was %s", tc.metadata.Keywords, nd.Metadata.Keywords)
			}
			if tc.metadata.Medium != nd.Metadata.Medium {
				t.Errorf("expected Medium to be %s, but was %s", tc.metadata.Medium, nd.Metadata.Medium)
			}
			if tc.metadata.Title != nd.Metadata.Title {
				t.Errorf("expected Title to be %s, but was %s", tc.metadata.Title, nd.Metadata.Title)
			}
			if tc.metadata.Product != nd.Metadata.Product {
				t.Errorf("expected Product to be %s, but was %s", tc.metadata.Product, nd.Metadata.Product)
			}
			if tc.metadata.Subject != nd.Metadata.Subject {
				t.Errorf("expected Subject to be %s, but was %s", tc.metadata.Subject, nd.Metadata.Subject)
			}
			if tc.metadata.Software != nd.Metadata.Software {
				t.Errorf("expected Software to be %s, but was %s", tc.metadata.Software, nd.Metadata.Software)
			}
			if tc.metadata.Source != nd.Metadata.Source {
				t.Errorf("expected Source to be %s, but was %s", tc.metadata.Source, nd.Metadata.Source)
			}
			if tc.metadata.Location != nd.Metadata.Location {
				t.Errorf("expected Location to be %s, but was %s", tc.metadata.Location, nd.Metadata.Location)
			}
			if tc.metadata.TrackNbr != nd.Metadata.TrackNbr {
				t.Errorf("expected TrackNbr to be %s, but was %s", tc.metadata.TrackNbr, nd.Metadata.TrackNbr)
			}
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
