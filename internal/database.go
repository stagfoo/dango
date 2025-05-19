package internal

import (
	"fmt"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
)

type MyDB struct {
	Version   int
	Items     []string
	Clipboard string
}

var usr, errorPath = user.Current()
var Path = usr.HomeDir + "/.config/dango/dango.toml"

func ViewDB(path string) MyDB {
	// view the database
	doc, readErr := os.ReadFile(Path)
	if readErr != nil {
		panic(readErr)
	}
	var db MyDB
	err := toml.Unmarshal([]byte(doc), &db)
	if err != nil {
		panic(err)
	}
	return db
}

func SaveDb(db MyDB) bool {
	// save the database
	fmt.Print(db.Clipboard)
	if db.Clipboard == "" {
		db.Clipboard = "pbcopy"
	}
	if db.Version == 0 {
		db.Version = 1
	}
	b, err := toml.Marshal(db)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(Path, b, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return false
	}
	return true
}

func AddToDB(filepath string) bool {
	// add a name and fileapth to the database
	db := ViewDB(Path)
	// If duplicate
	for _, item := range db.Items {
		if item == filepath {
			return true
		}
	}
	db.Items = append(db.Items, filepath)
	return SaveDb(db)
}

func RemoveFromDB(filepath string) error {
	// remove a name and filepath from the database
	db := ViewDB(Path)
	for i, item := range db.Items {
		if item == filepath {
			db.Items = append(db.Items[:i], db.Items[i+1:]...)
			SaveDb(db)
			return nil
		}
	}
	return fmt.Errorf("item not found in database")
}

func FindInDB(filepath string) string {
	// TODO replace with bubble tea fuzzy search
	db := ViewDB(Path)
	for _, item := range db.Items {
		if item == filepath {
			return item
		}
	}
	return ""
}
