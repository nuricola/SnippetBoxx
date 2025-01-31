package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nuricola/snippetbox/pkg/models/mysql"
)


type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

 
func main() {

	// Создание флагов для командной строки 
	addr := flag.String("addr",":4000","Сетевой адрес ХХТП")
	dsn := flag.String("dsn","web:pass@/snippetbox?parseTime=true", "Название MySQL источника данных")

	flag.Parse()

	// Создааем логи для откладки 
	infoLog := log.New(os.Stdout,"INFO \t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stderr,"ERROR \t", log.Ldate | log.Ltime | log.Lshortfile)

	// PULL Conection to DB
	 db, err := openDB(*dsn)
	 if err != nil{
		errorLog.Fatal(err)
	 }

	 defer db.Close()

	// Инициализируем новый кэш шаблона...

	templateCache, err :=newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	//Инициализируем новую структуру с зависимостями приложения.
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB:db},
		templateCache: templateCache,
	}


	// Инициализация структуру сервера
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

 
	//запуск сервера 
	infoLog.Printf("Запуск сервера на %s", *addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
}
 
func openDB(dsn string)(*sql.DB, error){
	db , err := sql.Open("mysql",dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil{
		return nil, err
	}
	return db, nil
}