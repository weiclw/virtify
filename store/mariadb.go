package store

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

func initMariaDB() {
    db, err := sql.Open("mysql", "username:password@tcp(hostname:port)/dbname")
    if err != nil {
        fmt.Println("Error: " + err.Error())
        return
    }

    defer db.Close()
}
