package main

import (
	"covid-tracker/server"
	"covid-tracker/utils"
)

func main() {
    db := utils.OpenDB();
    defer db.Close()

    r := server.Router(db)
    r.Run(":8080")
}
