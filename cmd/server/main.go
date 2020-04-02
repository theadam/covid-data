package main

import (
	"covid-tracker/server"
	"covid-tracker/utils"
	"os"
    "fmt"
)

func main() {
    fmt.Println(os.Getenv("PORT"))
    db := utils.OpenDB();
    defer db.Close()

    r := server.Router(db)
    r.Run()
}
