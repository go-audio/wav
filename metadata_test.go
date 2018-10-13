package wav

import (
	"os"
	"reflect"
	"testing"
)

func TestDecoder_ReadMetadata(t *testing.T) {
	testCases := []struct {
		in       string
		metadata *Metadata
	}{
		{in: "fixtures/listinfo.wav",
			metadata: &Metadata{
				Artist: "artist", Title: "track title", Product: "album title",
				TrackNbr: "42", CreationDate: "2017", Genre: "genre", Comments: "my comment",
			},
		},
		{in: "fixtures/kick.wav"},
		{in: "fixtures/flloop.wav", metadata: &Metadata{
			Software: "FL Studio (beta)",
			CuePoints: []*CuePoint{
				0:  {ID: [4]uint8{0x1, 0x0, 0x0, 0x0}, Position: 0x0, DataChunkID: [4]uint8{'d', 'a', 't', 'a'}},
				1:  {ID: [4]uint8{0x2, 0x0, 0x0, 0x0}, Position: 0x1a5e, DataChunkID: [4]uint8{'d', 'a', 't', 'a'}, SampleOffset: 0x1a5e},
				2:  {ID: [4]uint8{0x3, 0x0, 0x0, 0x0}, Position: 0x34bc, DataChunkID: [4]uint8{'d', 'a', 't', 'a'}, SampleOffset: 0x34bc},
				3:  {ID: [4]uint8{0x4, 0x0, 0x0, 0x0}, Position: 0x4f1a, DataChunkID: [4]uint8{'d', 'a', 't', 'a'}, SampleOffset: 0x4f1a},
				4:  {ID: [4]uint8{0x5, 0x0, 0x0, 0x0}, Position: 0x6978, DataChunkID: [4]uint8{'d', 'a', 't', 'a'}, SampleOffset: 0x6978},
				5:  {ID: [4]uint8{0x6, 0x0, 0x0, 0x0}, Position: 0x83d6, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x83d6},
				6:  {ID: [4]uint8{0x7, 0x0, 0x0, 0x0}, Position: 0x9e34, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x9e34},
				7:  {ID: [4]uint8{0x8, 0x0, 0x0, 0x0}, Position: 0xb892, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0xb892},
				8:  {ID: [4]uint8{0x9, 0x0, 0x0, 0x0}, Position: 0xd2f0, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0xd2f0},
				9:  {ID: [4]uint8{0xa, 0x0, 0x0, 0x0}, Position: 0xed4e, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0xed4e},
				10: {ID: [4]uint8{0xb, 0x0, 0x0, 0x0}, Position: 0x107ac, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x107ac},
				11: {ID: [4]uint8{0xc, 0x0, 0x0, 0x0}, Position: 0x1220a, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x1220a},
				12: {ID: [4]uint8{0xd, 0x0, 0x0, 0x0}, Position: 0x13c68, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x13c68},
				13: {ID: [4]uint8{0xe, 0x0, 0x0, 0x0}, Position: 0x156c6, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x156c6},
				14: {ID: [4]uint8{0xf, 0x0, 0x0, 0x0}, Position: 0x17124, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x17124},
				15: {ID: [4]uint8{0x10, 0x0, 0x0, 0x0}, Position: 0x18b82, DataChunkID: [4]uint8{0x64, 0x61, 0x74, 0x61}, SampleOffset: 0x18b82},
			},
			SamplerInfo: &SamplerInfo{SamplePeriod: 22676, MIDIUnityNote: 60, NumSampleLoops: 1,
				Loops: []*SampleLoop{
					{CuePointID: [4]byte{0, 0, 2, 0}, Type: 1024, Start: 0, End: 107999, Fraction: 0, PlayCount: 0},
				}},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			f, err := os.Open(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			d := NewDecoder(f)
			d.ReadMetadata()
			if err = d.Err(); err != nil {
				t.Fatal(err)
			}
			if tc.metadata != nil {
				if tc.metadata.SamplerInfo != nil {
					if !reflect.DeepEqual(tc.metadata.SamplerInfo, d.Metadata.SamplerInfo) {
						t.Fatalf("Expected sampler info\n%#v to equal\n%#v\n", d.Metadata.SamplerInfo, tc.metadata.SamplerInfo)
					}
				}
				if tc.metadata.CuePoints != nil {
					if !reflect.DeepEqual(tc.metadata.CuePoints, d.Metadata.CuePoints) {
						for i, c := range d.Metadata.CuePoints {
							if !reflect.DeepEqual(c, tc.metadata.CuePoints[i]) {
								t.Errorf("[%d] expected %#v got %#v", i, tc.metadata.CuePoints[i], c)
							}
						}
						t.Errorf("Expected cue points\n%#v to equal\n%#v\n", d.Metadata.CuePoints, tc.metadata.CuePoints)
					}
				}

				if !reflect.DeepEqual(tc.metadata, d.Metadata) {
					t.Fatalf("Expected\n%#v\n to equal\n%#v\n", d.Metadata, tc.metadata)
				}
			}
			f.Close()
		})
	}
}
