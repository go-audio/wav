package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-audio/wav"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You must pass the pass the path of the file to decode")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := wav.NewDecoder(f)
	dec.ReadMetadata()
	if err := dec.Err(); err != nil {
		log.Fatal(err)
	}
	if dec.Metadata != nil {
		fmt.Printf("SampleInfo: %+v\n", dec.Metadata.SamplerInfo)
		for i, l := range dec.Metadata.SamplerInfo.Loops {
			fmt.Printf("\tloop [%d]:\t%+v\n", i, l)
		}
		for i, c := range dec.Metadata.CuePoints {
			fmt.Printf("\tcue point [%d]:\t%+v\n", i, c)
		}
	}
}
