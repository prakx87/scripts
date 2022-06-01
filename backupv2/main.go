package main

import(
	"fmt"
	"log"
	"os"
	"database/sql"

	"github.com/jamf/go-mysqldump"
	"github.com/go-sql-driver/mysql"
	// "golang.org/x/net/context"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/drive/v3"
)

var (
	WarningLogger *log.Logger
	InfoLogger *log.Logger
	ErrorLogger *log.Logger
)


func init() {
	// If file does not exist, then create it or append to file
	file, err := os.OpenFile("/var/log/dxbackup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// setup mysql connection object
	fmt.Println("Create MySQL connection")
	config := mysql.NewConfig()
	config.User = "backup_user"
	config.Passwd = "NQNRQw4tTbvDu7n8"
	config.DBName = "dorama_vbull"
	config.Net = "tcp"
	config.Addr = "127.0.0.1"

	dumpDir := "/root/backups"
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", config.DBName)

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully connected to MySQL database")

	// use mysql object to take dump
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("DB taken successfully and saved at %s", dumpFilenameFormat)
	dumper.Close()
	
	// upload dump to google and yandex drive
	// cleanup older dumps
}