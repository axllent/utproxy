package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	// MyLSQL driver
	_ "github.com/go-sql-driver/mysql"
)

// Check returns a test
func Check(v *viper.Viper) error {

	ep := v.GetString("Endpoint")
	user := v.GetString("User")
	pass := v.GetString("Pass")

	if ep == "" {
		return errors.New("No endpoint set")
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s)/", user, pass, ep)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return err
	}

	defer db.Close()

	// make sure connection is available
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}
