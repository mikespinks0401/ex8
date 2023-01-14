package main

import (
	"github.com/mikespinks0401/number-normalizer/repositories/phoneNumber"
	"fmt"
)



func main(){
	numberRepo := phoneNumber.SQLiteRepository{}
	if err := numberRepo.NewSQLiteRepository("phoneNumbers.db"); err != nil {
		panic(err)
	}

	
	if err := numberRepo.RemoveDuplicates(); err != nil {
		panic(err)
	}
	list, err := numberRepo.GetAll()
	if err != nil {
		panic(err)
	}
	for _,val := range list {
		fmt.Println(val)
	}
}