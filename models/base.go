package models

//
import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv" // this package is used for loading the .env file
)

var databaseInstance *gorm.DB // Globally declaring the database instancegit remote add origin git@github.com:soumyabrata09/googleCloudBuildSolutions.git
//retirn a handle to the database object
func getDB() *gorm.DB {
	return databaseInstance
}

//init() automatically get called by Go,
//the code retrieve connection information from .env file then build a connection string and use it to connect to the database.
//Follwoing method is opening database connection
func init() {
	e := godotenv.Load() // Loading the .env file
	if e != nil {
		fmt.Print(e)
	}

	userName := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=false password=%", dbHost, userName, dbName, password)
	fmt.Println(dbUri) //for testing purposes

	//connecting the postgres database
	connectionInstance, conn_err := gorm.Open("postgres", dbUri)
	if conn_err != nil {
		fmt.Print(conn_err)
	}
	databaseInstance = connectionInstance
	databaseInstance.Debug().AutoMigrate(&Account{})
}
