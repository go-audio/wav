package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mattetti/audio/riff"
)

var (
	// See http://bwfmetaedit.sourceforge.net/listinfo.html
	markerIART    = [4]byte{'I', 'A', 'R', 'T'}
	markerISFT    = [4]byte{'I', 'S', 'F', 'T'}
	markerICRD    = [4]byte{'I', 'C', 'R', 'D'}
	markerICOP    = [4]byte{'I', 'C', 'O', 'P'}
	markerIARL    = [4]byte{'I', 'A', 'R', 'L'}
	markerINAM    = [4]byte{'I', 'N', 'A', 'M'}
	markerIENG    = [4]byte{'I', 'E', 'N', 'G'}
	markerIGNR    = [4]byte{'I', 'G', 'N', 'R'}
	markerIPRD    = [4]byte{'I', 'P', 'R', 'D'}
	markerISRC    = [4]byte{'I', 'S', 'R', 'C'}
	markerISBJ    = [4]byte{'I', 'S', 'B', 'J'}
	markerICMT    = [4]byte{'I', 'C', 'M', 'T'}
	markerIAUT    = [4]byte{'I', 'A', 'U', 'T'}
	markerITRK    = [4]byte{'I', 'T', 'R', 'K'}
	markerITRKBug = [4]byte{'i', 't', 'r', 'k'}
	markerITCH    = [4]byte{'I', 'T', 'C', 'H'}
	markerIKEY    = [4]byte{'I', 'K', 'E', 'Y'}
	markerIMED    = [4]byte{'I', 'M', 'E', 'D'}
)

// DecodeListChunk decodes a LIST chunk
func DecodeListChunk(d *Decoder, ch *riff.Chunk) error {
	if ch == nil {
		return fmt.Errorf("can't decode a nil chunk")
	}
	if d == nil {
		return fmt.Errorf("nil decoder")
	}
	if ch.ID == CIDList {
		// read the entire chunk in memory
		buf := make([]byte, ch.Size)
		var err error
		if _, err = ch.Read(buf); err != nil {
			return fmt.Errorf("failed to read the LIST chunk - %v", err)
		}
		r := bytes.NewReader(buf)
		// INFO subchunk
		scratch := make([]byte, 4)
		if _, err = r.Read(scratch); err != nil {
			return fmt.Errorf("failed to read the INFO subchunk - %v", err)
		}
		if !bytes.Equal(scratch, []byte{'I', 'N', 'F', 'O'}) {
			return fmt.Errorf("expected an INFO subchunk but got %+v", scratch)
		}
		if d.Metadata == nil {
			d.Metadata = &Metadata{}
		}

		// the rest is a list of string entries
		var (
			id   [4]byte
			size uint32
		)
		readSubHeader := func() error {
			if err := binary.Read(r, binary.BigEndian, &id); err != nil {
				return err
			}
			if err := binary.Read(r, binary.LittleEndian, &size); err != err {
				return err
			}
			return nil
		}

		for err == nil {
			err = readSubHeader()
			if err != nil {
				break
			}
			scratch = make([]byte, size)
			r.Read(scratch)
			switch id {
			case markerIARL:
				d.Metadata.Location = nullTermStr(scratch)
			case markerIART:
				d.Metadata.Artist = nullTermStr(scratch)
			case markerISFT:
				d.Metadata.Software = nullTermStr(scratch)
			case markerICRD:
				d.Metadata.CreationDate = nullTermStr(scratch)
			case markerICOP:
				d.Metadata.Copyright = nullTermStr(scratch)
			case markerINAM:
				d.Metadata.Title = nullTermStr(scratch)
			case markerIENG:
				d.Metadata.Engineer = nullTermStr(scratch)
			case markerIGNR:
				d.Metadata.Genre = nullTermStr(scratch)
			case markerIPRD:
				d.Metadata.Product = nullTermStr(scratch)
			case markerISRC:
				d.Metadata.Source = nullTermStr(scratch)
			case markerISBJ:
				d.Metadata.Subject = nullTermStr(scratch)
			case markerICMT:
				d.Metadata.Comments = nullTermStr(scratch)
			case markerITRK, markerITRKBug:
				d.Metadata.TrackNbr = nullTermStr(scratch)
			case markerITCH:
				d.Metadata.Technician = nullTermStr(scratch)
			case markerIKEY:
				d.Metadata.Keywords = nullTermStr(scratch)
			case markerIMED:
				d.Metadata.Medium = nullTermStr(scratch)
			}
		}
	}
	ch.Drain()
	return nil
}
