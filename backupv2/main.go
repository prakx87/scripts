package main

// var (
// 	WarningLogger *log.Logger
// 	InfoLogger *log.Logger
// 	ErrorLogger *log.Logger
// )

// func init() {
// 	// If file does not exist, then create it or append to file
// 	file, err := os.OpenFile("/var/log/dxbackup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
// }

func main() {

	dbList := []string{"dorama_vbull"}

	for _, db := range dbList {
		takeDump(db)
	}

	// upload dump to google and yandex drive
	// cleanup older dumps
}
