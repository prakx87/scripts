package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/bigkevmcd/go-configparser"
	"github.com/go-sql-driver/mysql"
	"github.com/jamf/go-mysqldump"
)

type dumpDetails struct {
	dbIp       string
	dbCredFile string
	dumpPath   string
}

func takeDump(db string) {
	dumpInf := dumpDetails{
		dbIp:       "127.0.0.1",
		dbCredFile: "/etc/my.cnf.d/backup.cnf",
		dumpPath:   "/root/backups",
	}

	config := dumpInf.createDbConfig(db)
	dbconn := openDbConn(config)
	// dumpName := getDumpFileName(db)
	startDump(dbconn, dumpInf.dumpPath)
}

func getConfigCreds(filePath string) map[string]string {
	CnfFileData, err := configparser.NewConfigParserFromFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	creds, err := CnfFileData.Items("client")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return creds
}

func (d dumpDetails) createDbConfig(dbName string) *mysql.Config {
	// Get MySQL credentials
	dbCreds := getConfigCreds(d.dbCredFile)

	// setup mysql connection object
	fmt.Println("Create MySQL connection")
	config := mysql.NewConfig()
	config.User = dbCreds["user"]
	config.Passwd = dbCreds["password"]
	config.DBName = dbName
	config.Net = "tcp"
	config.Addr = d.dbIp
	config.ParseTime = true

	return config
}

func getDumpFileName(dbName string) string {
	t := time.Now()
	const layoutDUMP = "20060102150405"
	return dbName + "_" + t.Format(layoutDUMP)
}

func openDbConn(config *mysql.Config) *sql.DB {
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Println("Error opening database: ", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to MySQL database")

	return db
}

func startDump(dbconn *sql.DB, dumpDir string) {
	// use mysql object to take dump
	const layoutDUMP = "_20060102150405"
	dumper, err := mysqldump.Register(dbconn, dumpDir, layoutDUMP)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		os.Exit(1)
	}

	fmt.Println("Registered database for MySQL dump")

	// Dump database to file
	outErr := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", outErr)
		os.Exit(1)
	}

	fmt.Printf("DB taken successfully and saved at %s/%s", dumpDir, dumpFilename)
	dumper.Close()
	dbconn.Close()
}
