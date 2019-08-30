package bank

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

type Transaction struct {
	Payee    string
	Price    int
	Date     string
	Category string
}

type Bank interface {
	Fetch(filename, sinceDate string) (txs []Transaction, err error)
}

type CSV struct{}

func New() *CSV {
	return &CSV{}
}

func (c *CSV) Fetch(filename, sinceDate string) (txs []Transaction, err error) {
	if filename == "" {
		filename = "/Users/luisvieira/Desktop/n26-csv-transactions-2.csv"
	}
	csvfile, err := os.Open(filename)
	if err != nil {
		return
	}

	r := csv.NewReader(csvfile)
	r.Comma = ';'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return txs, err
		}

		if since, err := time.Parse("2006-01-02", sinceDate); err == nil {
			if recordDate, err := time.Parse("2006-01-02", record[0]); err == nil {
				if since.Equal(recordDate) || since.Before(recordDate) {
					f, err := strconv.ParseFloat(record[3], 64)
					if err != nil {
						return txs, err
					}
					if f < 0 {
						tx := Transaction{
							Date:     record[0],
							Payee:    record[1],
							Category: record[2],
							Price:    int(f * 1000),
						}

						txs = append(txs, tx)
					}
				}
			}
		}

	}

	return txs, nil
}
