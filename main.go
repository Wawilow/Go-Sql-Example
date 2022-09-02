package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

type Event struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type Database struct {
	dbObj   *sql.DB
	history Dictionary
}

type Dictionary []Event

func databaseConnect(cnf Config) (Database, error) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cnf.login, cnf.password, cnf.addr, cnf.port, cnf.databaseName)
	db, err := sql.Open("mysql", addr)

	if err != nil {
		log.Printf("Error connecting to database %s:%s/%s", cnf.addr, cnf.port, cnf.databaseName)
		return Database{}, err
	} else {
		log.Printf("Connect to database %s:%s/%s", cnf.addr, cnf.port, cnf.databaseName)
		hstr := Dictionary{Event{"crate", "", fmt.Sprintf("%s", time.Now())}}
		return Database{db, hstr}, nil
	}
}

func (db Database) query(command string) (*sql.Rows, error) {
	rows, err := db.dbObj.Query(command)
	if err != nil {
		return nil, err
	}
	//defer rows.Close()

	return rows, nil
}

func (db Database) create(tableName string, columns string) (bool, error) {
	results, err := db.dbObj.Query(fmt.Sprintf("CREATE TABLE %s (%s)", tableName, columns))
	if err != nil {
		return false, err
	}
	if results.Err() == nil {
		return true, nil
	} else {
		return false, nil
	}
}

func (db Database) insert(tableName string, columns []string) (bool, error) {
	command := fmt.Sprintf("INSERT INTO %s VALUES %s", tableName, strings.Join(columns, ", "))
	results, err := db.dbObj.Query(command)
	log.Print(command)
	if err != nil {
		return false, err
	}
	if results.Err() == nil {
		return true, nil
	} else {
		return false, nil
	}
}

func (db Database) get(tableName string, columns []string) (*sql.Rows, error) {
	var col string
	if len(columns) == 0 {
		col = "*"
	} else {
		col = strings.Join(columns, ", ")
	}
	command := fmt.Sprintf("SELECT %s FROM %s", col, tableName)
	rows, err := db.dbObj.Query(command)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db Database) close() (bool, error) {
	err := db.dbObj.Close()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (db Database) getProcessingExample() bool {
	log.Println("GO sql example")

	var (
		id   int
		name string
	)

	res, err := db.get("users", []string{})
	for res.Next() {
		err = res.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}

	db.close()
	log.Println("End!")
	if err != nil {
		return false
	} else {
		return true
	}
}

func main() {
	cnf := config()
	db, err := databaseConnect(cnf)
	if err != nil {
		panic(err.Error())
	}

	log.Print(db.getProcessingExample())

}

