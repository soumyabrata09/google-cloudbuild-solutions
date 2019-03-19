package models

import (
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	u "../utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

//Folllwoing is the jwt standard claims struct
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//DEFINITION of Validate()
//  function to validate the imcoming request to the api
//very basic val;idatation is being carried out in this method
//followoing is a pointer receiver of Account struct type
func (accountStruct *Account) Validate() (map[string]interface{}, bool) {
	if strings.Contains(accountStruct.Email, "@") {
		return u.Message(false, "Email id is required"), false
	}
	if len(accountStruct.Password) < 8 {
		return u.Message(false, "Password is mismatched"), false
	}
	temp := &Account{}
	//checking duplicate email ids
	accountError := GetDB().Table("accounts").Where("email = ?", accountStruct.Email).First(temp).Error
	if accountError != nil && accountError != gorm.ErrRecordNotFound {
		return u.Message(false, "Interrupeted Network Connection, try again"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email Id is already taken,choose another"), false
	}

	return u.Message(false, "Data Validate Successfully "), true
}

//DEFINITION of CreateUser()
//followoing is a pointer receiver of Account struct type
func (accountStruct *Account) CreateUser() map[string]interface{} {
	//here we are getting a respone after validating it ,hence we are calling validate method receiver of type Account
	if resp, ok := accountStruct.Validate(); !ok {
		return resp
	}
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(accountStruct.Password), bcrypt.DefaultCost)
	accountStruct.Password = string(hashedPwd) // converting the hashed password to a type of string
	//creating the db of type accountstruct
	GetDB().Create(accountStruct)
	//check for error during account reation process
	if accountStruct.ID <= 0 {
		return u.Message(false, "Failed to create account, may be due to connection error")
	}

	//Jwt token creation
	tokenData := &Token{
		UserId: accountStruct.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenData) //using HS256 hashing technique
	//then we are applying this hashing technique to the password data loaded from the .env file
	tokenStringData, _ := token.SignedString([]byte(os.Getenv("token_password")))
	//Adding the token to the Account struct ( as it is having token filed defined within it)
	accountStruct.Token = tokenStringData
	accountStruct.Password = "" // deleting the password itself
	resposnse := u.Message(true, "Account Successfully Created!!!")
	resposnse["account"] = accountStruct
	return resposnse
}

//DEFINITION of GetUser()
//by using the user id we will be fetching the user's respoective account information, which is why we are returning the Account struct type here
func GetUser(uid uint) *Account {
	accountStruct := &Account{}
	//getting user from the the database by checking the uid passed
	GetDB().Table("accounts").Where("id=?", uid).First(accountStruct)
	//checking for the authenicity of the retrieved data if any
	if accountStruct.Email == "" {
		return nil // as email id mismatched
	}
	accountStruct.Password = ""
	return accountStruct
}

//followoing is a function with a return type map
//DEFINITION of Login() method
func Login(email, password string) map[string]interface{} {
	// getting pointer of Account struct
	accountStruct := &Account{}
	//checking the autheticity of the email
	err := GetDB().Table("accounts").Where("email=?", email).First(accountStruct).Error
	if err != nil {
		//gorm specific exception thrown
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address is wrong,try again")
		}
		return u.Message(false, "Connection Error!!!")
	}
	//compare the hash and loaded password
	err = bcrypt.CompareHashAndPassword([]byte(accountStruct.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrHashTooShort {
			return u.Message(false, "Invalid Credential[HashTooShort]")
		}
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return u.Message(false, "Invalid Credential[PasswordMismatched]")
		}
	}

	//successfully logged in
	//deleting the password data to prevent any backdoor loophole
	accountStruct.Password = ""
	//create jwt token
	tokenData := &Token{UserId: accountStruct.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenData) // decide which hashing technique is goin to be used
	//assign the token technique to the password loaded from the .env file
	tokenStringData, _ := token.SignedString([]byte(os.Getenv("token_password")))
	accountStruct.Token = tokenStringData // storing the token in the response
	response := u.Message(true, "You are logged in")
	response["account"] = accountStruct
	return response
}
