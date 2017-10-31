package srt

import (
	"log"
	"testing"
)

func TestReadFile(t *testing.T) {
	elements := make(chan Element)
	go Read("sample_eng.srt", elements)

	for element := range elements {
		// process the data...
		log.Println(element)
	}
}

func TestWriteFile(t *testing.T) {
	elements := make(chan Element)
	go Read("sample_eng.srt", elements)

	Write("tam.srt", elements)
}
