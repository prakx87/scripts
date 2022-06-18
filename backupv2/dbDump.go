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
	dbList     []string
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
	dumpName := getDumpFileName(db)
	startDump(dbconn, dumpInf.dumpPath, dumpName)
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

	return config
}

func getDumpFileName(dbName string) string {
	t := time.Now()
	const layoutDUMP = "20060102150405"
	return dbName + "-" + t.Format(layoutDUMP) + ".sql"
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

func startDump(dbconn *sql.DB, dumpDir string, dumpFilename string) {
	// use mysql object to take dump
	dumper, err := mysqldump.Register(dbconn, dumpDir, dumpFilename)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	// Dump database to file
	resultFilename := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return
	}
	fmt.Printf("File is saved to %s", resultFilename)

	fmt.Printf("DB taken successfully and saved at %s/%s", dumpDir, dumpFilename)
	dumper.Close()
}
