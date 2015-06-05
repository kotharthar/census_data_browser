package main

import (
    "strings"
    "database/sql"
    "log"
    "fmt"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func queryApi(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    table := params.Get("t")

    db, err := sql.Open("sqlite3", "../data/main.db")
    checkErr(err)

    rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s;",table))
    checkErr(err)
    defer rows.Close()

    w.Header().Set("Content-Type", "text/html")

    columnNames, err := rows.Columns()
    checkErr(err)

    cells := make([]string, len(columnNames))
    columns := make([]interface{}, len(columnNames))
    columnPointers := make([]interface{}, len(columnNames))
    for i := 0; i < len(columnNames); i++ {
        columnPointers[i] = &columns[i]
    }

    for rows.Next() {
        err := rows.Scan(columnPointers...)
        checkErr(err)
        for j, col := range columns {
            switch v := col.(type) {
            case []byte:
                cells[j] = fmt.Sprintf("%s",v)
            case string:
                cells[j] = fmt.Sprintf("%v",v)
            case int32, int64:
                cells[j] = fmt.Sprintf("%v",v)
            case float32, float64:
                cells[j] = fmt.Sprintf("%v",v)
            default:
                cells[j] = "-"
            }
        }
        result := fmt.Sprintf("%s\n",strings.Join(cells,","))
        fmt.Fprintf(w,result)
    }
}

func main() {
    fs := http.FileServer(http.Dir("static"))

    http.Handle("/", fs)
    http.HandleFunc("/q", queryApi)
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}
