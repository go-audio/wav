package main

import (
	"flag"
	"log"
	"math"
	"os"

	"github.com/go-audio/wav"
)

func main() {
	output := flag.String("output", "output.wav", "filename to write to")
	frequency := flag.Float64("frequency", 440, "frequency in hertz to generate")
	length := flag.Float64("lenght", 5, "length in seconds of output file")
	flag.Parse()

	log.Printf("generating a %f sec sine wav at %f hz", *length, *frequency)
	f, err := os.Create(*output)
	if err != nil {
		log.Fatalf("error creating %s: %s", *output, err)
	}
	defer f.Close()

	const sampleRate = 48000
	wavOut := wav.NewEncoder(f, sampleRate, 16, 1, 1)
	numSamples := int(sampleRate * *length)
	defer wavOut.Close()

	for i := 0; i < numSamples; i++ {
		fv := math.Sin(float64(i) / sampleRate * *frequency * 2 * math.Pi)
		v := uint16(fv * 32767)
		wavOut.WriteFrame(v)
	}
}
