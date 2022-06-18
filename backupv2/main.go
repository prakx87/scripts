package main

func main() {

	dbList := []string{"dorama_vbull"}

	for _, db := range dbList {
		takeDump(db)
	}

	// upload dump to google and yandex drive
	// cleanup older dumps
}
