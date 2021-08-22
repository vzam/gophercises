package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ChapterOption refers to an other chapter which will be reached
// when choosing that option and a text to display for the option.
type ChapterOption struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// Story is a paragraph in a Chapter.
type Story string

// Chapter has a title and contains the story paragraphs as well
// as the options for the next chapters. Note, that the options
// can be empty, in which case the adventure has come to an end.
type Chapter struct {
	Title           string          `json:"title"`
	StoryParagraphs []Story         `json:"story"`
	Options         []ChapterOption `json:"options"`
}

// Adventure is a collection of chapters and the definition of the
// initial chapter.
type Adventure struct {
	InitialChapter string             `json:"initial-chapter"`
	Chapters       map[string]Chapter `json:"chapters"`
}

// ParseJsonAdventure takes an io.Reader containing json, which
// will be parsed into an adventure and then returned.
// The error will be non-nil when the json is invalid or when the
// json does not fit the data structure.
func ParseJsonAdventure(jsonReader io.Reader) (*Adventure, error) {
	d := json.NewDecoder(jsonReader)

	var adventure Adventure
	err := d.Decode(&adventure)
	if err != nil {
		return nil, fmt.Errorf("unable to parse adventure: %v", err)
	}

	return &adventure, nil
}

func ParseJsonFileAdventure(filename string) (*Adventure, error) {
	adv, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open adventure file: %v", err)
	}
	defer adv.Close()

	return ParseJsonAdventure(adv)
}
