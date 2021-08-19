package main

import (
	"encoding/csv"
)

// Reads problems into the Problem type
type ProblemReader interface {
	// Reads a single problem from the source
	Read() (problem *Problem, err error)

	// Reads all problems from the source
	ReadAll() (problems []*Problem, err error)
}

type csvProblemReader struct {
	csv *csv.Reader
}

// make sure that the csvProblemReader implements ProblemReader
var _ ProblemReader = (*csvProblemReader)(nil)

// Creates a new problem reader based on a csv reader. The csv reader must have two fields
// with the first one being the question and the second one being the answer to that question.
func NewReaderFromCsv(r *csv.Reader) ProblemReader {
	return &csvProblemReader{
		csv: r,
	}
}

func (r *csvProblemReader) Read() (problem *Problem, err error) {
	r.csv.FieldsPerRecord = 2

	row, err := r.csv.Read()
	if err != nil {
		return nil, err
	}

	return &Problem{Question: row[0], Answer: row[1]}, nil
}

func (r *csvProblemReader) ReadAll() (problems []*Problem, err error) {
	r.csv.FieldsPerRecord = 2

	rows, err := r.csv.ReadAll()
	if err != nil {
		return nil, err
	}

	problems = make([]*Problem, len(rows))

	for i, row := range rows {
		problems[i] = &Problem{Question: row[0], Answer: row[1]}
	}

	return problems, nil
}
