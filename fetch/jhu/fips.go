package jhu

import (
	"io"
	"os"
	"strings"
    "encoding/csv"
)

func csvFile() string {
    path, _ := os.Getwd()
    return path + "/fetch/jhu/fips-counties.csv"
}

func contents() io.Reader {
    contents, _ := os.Open(csvFile())
    return contents
}

type FipsData struct {
    Name string
    Fips string
    StateCode string
}

func ParseFips() map[string]FipsData {
	reader := csv.NewReader(contents())
    result := make(map[string]FipsData)

    for {
        reader, err := reader.Read()
		if err == io.EOF { break }
        if err != nil { panic("Error reading fips csv\n" + err.Error())}

        stateCode := reader[0]
        countyName := reader[3]
        fips := reader[1] + reader[2]


        key := fips

        result[strings.ToLower(key)] = FipsData{
            Name: countyName,
            Fips: fips,
            StateCode: stateCode,
        }
    }
    return result
}

var FipsMap = ParseFips()
