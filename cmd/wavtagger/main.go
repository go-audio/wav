// This command line tool helps the user tag wav files by injecting metadata in
// the file in a safe way.
// All files are copied and stored in the wavtagger folder by the original files.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-audio/wav"
)

var (
	flagFileToTag   = flag.String("file", "", "Path to the wave file to tag")
	flagDirToTag    = flag.String("dir", "", "Directory containing all the wav files to tag")
	flagTitleRegexp = flag.String("regexp", "", `submatch regexp to use to set the title dynamically by extracting it from the filename (ignoring the extension), example: 'my_files_\d\d_(.*)'`)
	//
	flagTitle     = flag.String("title", "", "File's title")
	flagArtist    = flag.String("artist", "", "File's artist")
	flagComments  = flag.String("comments", "", "File's comments")
	flagCopyright = flag.String("copyright", "", "File's copyright")
	flagGenre     = flag.String("genre", "", "File's genre")
	// TODO: add other supported metadata types
)

func main() {
	flag.Parse()
	if *flagFileToTag == "" && *flagDirToTag == "" {
		fmt.Println("You need to pass -file or -dir to indicate what file or folder content to tag.")
		os.Exit(1)
	}
	if *flagFileToTag != "" {
		if err := tagFile(*flagFileToTag); err != nil {
			fmt.Printf("Something went wrong when tagging %s - error: %v\n", *flagFileToTag, err)
			os.Exit(1)
		}
	}
	if *flagDirToTag != "" {
		var filePath string
		fileInfos, _ := ioutil.ReadDir(*flagDirToTag)
		for _, fi := range fileInfos {
			if strings.HasPrefix(
				strings.ToLower(filepath.Ext(fi.Name())),
				".wav") {
				filePath = filepath.Join(*flagDirToTag, fi.Name())
				if err := tagFile(filePath); err != nil {
					fmt.Printf("Something went wrong tagging %s - %v\n", filePath, err)
				}
			}
		}
	}
}

func tagFile(path string) error {
	in, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s - %v", path, err)
	}
	d := wav.NewDecoder(in)
	buf, err := d.FullPCMBuffer()
	if err != nil {
		return fmt.Errorf("couldn't read buffer %s %v", path, err)
	}
	in.Close()

	outputDir := filepath.Join(filepath.Dir(path), "wavtagger")
	outPath := filepath.Join(outputDir, filepath.Base(path))
	os.MkdirAll(outputDir, os.ModePerm)

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("couldn't create %s %v", outPath, err)
	}
	defer out.Close()

	e := wav.NewEncoder(out,
		buf.Format.SampleRate,
		int(d.BitDepth),
		buf.Format.NumChannels,
		int(d.WavAudioFormat))
	if err := e.Write(buf); err != nil {
		return fmt.Errorf("failed to write audio buffer - %v", err)
	}
	e.Metadata = &wav.Metadata{}
	if *flagArtist != "" {
		e.Metadata.Artist = *flagArtist
	}
	if *flagTitleRegexp != "" {
		filename := filepath.Base(path)
		filename = filename[:len(filename)-len(filepath.Ext(path))]
		re := regexp.MustCompile(*flagTitleRegexp)
		matches := re.FindStringSubmatch(filename)
		if len(matches) > 0 {
			e.Metadata.Title = matches[1]
		} else {
			fmt.Printf("No matches for title regexp %s in %s\n", *flagTitleRegexp, filename)
		}
	}
	if *flagTitle != "" {
		e.Metadata.Title = *flagTitle
	}

	if *flagComments != "" {
		e.Metadata.Comments = *flagComments
	}
	if *flagCopyright != "" {
		e.Metadata.Copyright = *flagCopyright
	}
	if *flagGenre != "" {
		e.Metadata.Genre = *flagGenre
	}
	if err := e.Close(); err != nil {
		return fmt.Errorf("failed to close %s - %v", outPath, err)
	}
	fmt.Println("Tagged file available at", outPath)

	return nil
}
