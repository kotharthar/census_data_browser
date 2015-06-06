package main

import (
    "encoding/json"
    "encoding/base64"
    "strings"
    "database/sql"
    "log"
    "fmt"
    "net/http"
    "html/template"
    _ "github.com/mattn/go-sqlite3"
    "github.com/serenize/snaker"
)

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

type TableTpl struct{
    Name string
    Fields []string
}

func ToCamel(args ...interface{}) string {
    return args[0].(string)
}

func pageTemplate(w http.ResponseWriter, r *http.Request) {
    var command string
    params := r.URL.Query()
    table := params.Get("a")

    db, err := sql.Open("sqlite3", "../data/main.db")
    checkErr(err)

    command = fmt.Sprintf("SELECT * FROM %s limit 1;",table)
    rows, err := db.Query(command)
    checkErr(err)
    defer rows.Close()

    columnNames, err := rows.Columns()
    checkErr(err)
    rows.Close()

    tplData := TableTpl{Name: table, Fields: columnNames, }

    fm := template.FuncMap{
        "toCamel": snaker.SnakeToCamel,
    }
    t, _ := template.New("aok").Funcs(fm).ParseFiles("t.html")
    w.Header().Set("Content-Type", "text/html")
    t.ExecuteTemplate(w, "t.html", &tplData)
}

func queryApi(w http.ResponseWriter, r *http.Request) {
    var command string
    params := r.URL.Query()
    table := params.Get("t")
    format := params.Get("f")
    filter := params.Get("w")

    db, err := sql.Open("sqlite3", "../data/main.db")
    checkErr(err)

    if(len(filter) > 0){
        fdata, err := base64.StdEncoding.DecodeString(filter)
        checkErr(err);
        command = fmt.Sprintf("SELECT * FROM %s WHERE %s;",table,fdata)

    }else{
        command = fmt.Sprintf("SELECT * FROM %s",table)
    }

    rows, err := db.Query(command)
    checkErr(err)
    defer rows.Close()
    // rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s;",table, filter))
    // checkErr(err)
    // defer rows.Close()

    columnNames, err := rows.Columns()
    checkErr(err)

    cells := make([]string, len(columnNames))
    columns := make([]interface{}, len(columnNames))
    columnPointers := make([]interface{}, len(columnNames))
    for i := 0; i < len(columnNames); i++ {
        columnPointers[i] = &columns[i]
    }

    if(format == "csv"){
        w.Header().Set("Content-Type", "text/csv")
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
    if(format == "json"){
        w.Header().Set("Content-Type", "text/json")
        jrows := []map[string]interface{}{}
        for rows.Next() {
            err := rows.Scan(columnPointers...)
            checkErr(err)
            jcell := make(map[string]interface{})
            for j, col := range columns {
                switch v := col.(type) {
                case []byte:
                    jcell[columnNames[j]] = fmt.Sprintf("%s",v)
                case string:
                    jcell[columnNames[j]] = fmt.Sprintf("%v",v)
                default:
                    jcell[columnNames[j]] = v
                }
            }
            jrows = append(jrows,jcell)
        }
        enc := json.NewEncoder(w)
        enc.Encode(jrows)
    }
}

func main() {
    fs := http.FileServer(http.Dir("static"))

    http.Handle("/", fs)
    http.HandleFunc("/t", pageTemplate)
    http.HandleFunc("/q", queryApi)
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Println("Listening...")
    http.ListenAndServe(":3000", nil)
}
