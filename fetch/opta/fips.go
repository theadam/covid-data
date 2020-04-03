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
    pattern, _ := regexp.Compile(" (County|Borough|Municipality|Parish|City and Borough)$")
    c = pattern.ReplaceAllString(c, "")

    pattern, _ = regexp.Compile(" city$")
    return strings.TrimSpace(pattern.ReplaceAllString(c, " City"))
}

type FipsData struct {
    Name string
    Fips string
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


        key := stateCode + "-" + cleanCounty(countyName)

        result[strings.ToLower(key)] = FipsData{
            Name: countyName,
            Fips: fips,
        }
    }
    return result
}
