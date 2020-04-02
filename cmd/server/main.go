package main

import (
	"covid-tracker/server"
	"covid-tracker/utils"
	"os"
)

func main() {
    db := utils.OpenDB();
    defer db.Close()

    r := server.Router(db)
    r.Run("0.0.0.0:" + os.Getenv("PORT"))
}
