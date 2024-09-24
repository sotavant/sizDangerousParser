package main

import (
	"encoding/csv"
	"os"
	"regexp"
	"strings"
)

type Row struct {
	Danger,
	DangerID,
	DangerEvent,
	DangerEventID string
}

const (
	dataFileName   = "files/siz_by_dangerous.csv"
	resultFileName = "files/result.csv"
	separator      = ';'
)

const indexOfFirstRowWithData = 2

const (
	dangerIdCellId      = 0
	dangerCellId        = 1
	dangerEventIDCellId = 2
	dangerEventCellId   = 3
)

const (
	dangerEventIdPattern = `^\d+\.\d+\.\d+`
	stringPattern        = `^[а-яА-Я]+`
)

func main() {
	parseFile()
}

func parseFile() {
	f, err := os.ReadFile(dataFileName)
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(strings.NewReader(string(f)))
	r.Comma = rune(separator)

	records, err := r.ReadAll()

	var rows []Row
	var prevRow Row

	for i, record := range records {
		if i < indexOfFirstRowWithData {
			continue
		}

		if row, ok := fillRow(record, &prevRow); ok {
			prevRow = row
			rows = append(rows, row)
		}
	}

	saveResult(rows)
}

func fillRow(record []string, prevRow *Row) (Row, bool) {
	var row Row

	for i, cell := range record {
		record[i] = strings.TrimSpace(cell)
	}

	if isFullRecord(record) {
		row.Danger = record[dangerCellId]
		row.DangerID = record[dangerIdCellId]
		row.DangerEvent = record[dangerEventCellId]
		row.DangerEventID = record[dangerEventIDCellId]

		return row, true
	} else if isDangerEvent(record) && prevRow != nil {
		row = *prevRow
		row.DangerEventID = record[dangerEventIDCellId]
		row.DangerEvent = record[dangerEventCellId]

		return row, true
	}

	return row, false
}

func isFullRecord(record []string) bool {
	rules := []map[string]string{
		{
			"pattern": `^\d+\.\d+`,
			"data":    record[dangerIdCellId],
		},
		{
			"pattern": stringPattern,
			"data":    record[dangerCellId],
		},
		{
			"pattern": dangerEventIdPattern,
			"data":    record[dangerEventIDCellId],
		},
		{
			"pattern": stringPattern,
			"data":    record[dangerEventCellId],
		},
	}

	for _, rule := range rules {
		ok, err := regexp.MatchString(rule["pattern"], rule["data"])
		if err != nil {
			panic(err)
		}

		if !ok {
			return false
		}
	}

	return true
}

func isDangerEvent(record []string) bool {
	rules := []map[string]string{
		{
			"pattern": `^$`,
			"data":    record[dangerIdCellId],
		},
		{
			"pattern": `^$`,
			"data":    record[dangerCellId],
		},
		{
			"pattern": dangerEventIdPattern,
			"data":    record[dangerEventIDCellId],
		},
		{
			"pattern": stringPattern,
			"data":    record[dangerEventCellId],
		},
	}

	for _, rule := range rules {
		ok, err := regexp.MatchString(rule["pattern"], rule["data"])
		if err != nil {
			panic(err)
		}

		if !ok {
			return false
		}
	}

	return true
}

func saveResult(rows []Row) {
	file, err := os.OpenFile(resultFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var data string
	var rowData [4]string

	for _, row := range rows {
		rowData[0] = row.DangerID
		rowData[1] = row.Danger
		rowData[2] = row.DangerEventID
		rowData[3] = row.DangerEvent

		data = strings.Join(rowData[:], string(separator))
		_, err = file.WriteString(data + "\n")
		if err != nil {
			panic(err)
		}
	}
}
