package utils

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"paypath/pkg/logger"
)

func ReadCSV(file string) ([]string, [][]string, bool) {
	f, err := os.Open(file)
	if err != nil {
		logger.Log.Warn().Str("file", file).Msg("seed file not found, skipping")
		return nil, nil, false
	}
	defer f.Close()
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		return nil, nil, false
	}
	headers := make([]string, len(records[0]))
	for i, h := range records[0] {
		headers[i] = strings.ToLower(strings.TrimSpace(h))
	}
	return headers, records[1:], true
}

func CSVVal(v string) string { return strings.TrimSpace(v) }

func CSVFloat(v string) *float64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &n
}

func CSVInt(v string) *int {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &n
}

func HeaderIndex(headers []string) map[string]int {
	m := make(map[string]int, len(headers))
	for i, h := range headers {
		m[h] = i
	}
	return m
}
