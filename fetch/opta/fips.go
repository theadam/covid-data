package opta

import (
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strings"
)

func csvFile() string {
    path, _ := os.Getwd()
    return path + "/fetch/opta/fips-counties.csv"
}

func contents() io.Reader {
    contents, _ := os.Open(csvFile())
    return contents
}

func cleanCounty(c string) string {
    pattern, _ := regexp.Compile(" (County|Parish|Municipality|Borough|City and Borough)$")
    c = pattern.ReplaceAllString(c, "")

    pattern, _ = regexp.Compile(" city$")
    return pattern.ReplaceAllString(c, " City")
}

func ParseFips() map[string]string {
	reader := csv.NewReader(contents())
    result := make(map[string]string)

    for {
        reader, err := reader.Read()
		if err == io.EOF { break }
        if err != nil { panic("Error reading fips csv\n" + err.Error())}

        stateCode := reader[0]
        countyName := reader[3]
        fips := reader[1] + reader[2]


        key := stateCode + "-" + cleanCounty(countyName)

        result[strings.ToLower(key)] = fips
    }
    return result
}
