package phoneNumber

import (
	"database/sql"
	"fmt"

	"github.com/mikespinks0401/number-normalizer/db"
	"github.com/mikespinks0401/number-normalizer/model"
)

type SQLiteRepository struct {
	db *sql.DB
}

func (r *SQLiteRepository)NewSQLiteRepository(path string) (error) {
	db, err := db.OpenDB(path)
	if err != nil {
		return err
	}

	r.db = db 
	_, err = r.db.Exec("CREATE TABLE IF NOT EXISTS phone_numbers (id INTEGER PRIMARY KEY AUTOINCREMENT, number TEXT)")
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteRepository) Get(id int)(PhoneNumber.PhoneNumber, error){
	var phoneNumber PhoneNumber.PhoneNumber
	err := r.db.QueryRow("SELECT id FROM phone_numbers WHERE id=?", id).Scan(&phoneNumber)
	if err != nil {
		return PhoneNumber.PhoneNumber{}, err
	}
	return phoneNumber, nil
}

func (r *SQLiteRepository) GetAll() ([]PhoneNumber.PhoneNumber, error){
	rows,err := r.db.Query("SELECT id, number FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var phoneNumbers []PhoneNumber.PhoneNumber
	for rows.Next(){
		var number PhoneNumber.PhoneNumber
		err := rows.Scan(&number.Id, &number.Number)
		if err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, number)
	}
	return phoneNumbers, nil
}

func (r *SQLiteRepository)Insert( number string) error{
	stmt, err := r.db.Prepare("INSERT INTO phone_numbers(number) VALUES(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(number)
	if err != nil {
		return err
	}	
	return nil
}

func (r *SQLiteRepository)UpdateEach( f func(string)(string,error))error{
	rows, err := r.db.Query("SELECT id, number FROM phone_numbers")
	if err != nil {
		return err	
	}
	defer rows.Close()
	var numbers []PhoneNumber.PhoneNumber
	for rows.Next(){
		var number PhoneNumber.PhoneNumber
		if err := rows.Scan(&number.Id, &number.Number); err != nil {
			return err
		}
		numbers = append(numbers,number)
	}
	for _,val := range numbers{
		stmt, err := r.db.Prepare("UPDATE phone_numbers SET number=? WHERE id=?")
		if err != nil{
			return err
		}
		defer stmt.Close()
		newNumber, err := f(val.Number)
		res, stmtErr := stmt.Exec(newNumber, val.Id)
		if stmtErr != nil {
			return err
		}
		if int,_ := res.RowsAffected(); int != 1 {
			return fmt.Errorf("error updating the row")
		}
	}
	return nil
}

func (d *SQLiteRepository)RemoveDuplicates()error {
	list, err := d.GetAll()
	if err != nil {
		return err
	}
	numberMap := make(map[string]int)
	for _,val := range list{
		if _, ok := numberMap[val.Number]; !ok{
			fmt.Println("Doesn't Exist")
			numberMap[val.Number] = val.Id
		}else{
			
			oldval :=numberMap[val.Number]
			numberMap[val.Number] = val.Id
			stmt, err := d.db.Prepare("DELETE FROM phone_numbers WHERE id=?")
			if err != nil{
				return err
			}
			_, stmtErr := stmt.Exec(oldval)
			if stmtErr != nil {
				return stmtErr
			}
		}
	}
	return nil
}

func checkIfExist(og, new PhoneNumber.PhoneNumber)bool{
	return og.Number == new.Number
}