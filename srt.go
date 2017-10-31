package srt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Element is the wrapper for the element/s identified by the SRT parser
type Element struct {
	Index     uint
	Timestamp string
	Subtitles []string
}

type ElementHandler interface {
	onIndexLine(e Element)
	onNewLine(e Element)
}

const _ELEMENT_INDEX uint = 0
const _ELEMENT_TIMESTAMP uint = 1
const _ELEMENT_SUBTITLE uint = 2

// ReadFile accepts the srt file to be read and parses the file line by line
// Implement the methods onXXXLine to be notified as and when the elements are found by the parser
func Read(f string, channel chan Element) {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var index uint
	var element Element

	for scanner.Scan() { // internally, it advances token based on seperator
		token := scanner.Text()
		switch index {
		case _ELEMENT_INDEX:
			{
				row, err := strconv.Atoi(token)
				if err != nil {
					fmt.Println(err)
				}
				element.Index = uint(row)
				//channel <- element
				index++
				break
			}
		case _ELEMENT_TIMESTAMP:
			{
				element.Timestamp = token
				//channel <- element
				index++
				break
			}
		default:
			{
				// if it is empty it means that
				if token == "" {
					index = _ELEMENT_INDEX
					channel <- element
					element = Element{}
					break
				}

				element.Subtitles = append(element.Subtitles, token)
			}
		}
	}

	// if the index wasnt reset it means that the SRT file didn't end with a new line.
	// so reset the index and call the onNewLine for the last subtitle element.
	if index != _ELEMENT_INDEX {
		index = _ELEMENT_INDEX
		channel <- element
		element = Element{}
	}

	close(channel)
}

func Write(f string, elements chan Element) {
	file, err := os.Create(f)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	w := bufio.NewWriter(file)

	for element := range elements {
		// process the data...
		w.WriteString(fmt.Sprintf("%d\n", element.Index))
		w.WriteString(fmt.Sprintln(element.Timestamp))
		for _, subtitle := range element.Subtitles {
			w.WriteString(fmt.Sprintln(subtitle))
		}
		w.WriteString(fmt.Sprintln())
	}

	w.Flush()
}
