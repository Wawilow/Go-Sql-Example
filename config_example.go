package main

type Config struct {
	login        string //login
	password     string //password
	addr         string //server address
	port         string //server port
	databaseName string // database name
}

func config() Config {
	login := ""        //login
	password := ""     //password
	addr := ""         //addr
	port := ""         //port
	databaseName := "" //db name
	result := Config{login, password, addr, port, databaseName}
	return result
}
