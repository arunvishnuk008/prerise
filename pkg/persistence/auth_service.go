package persistence

import (
	"bytes"
	"crimson-sunrise.site/pkg/common"
	"crimson-sunrise.site/pkg/db"
	"crimson-sunrise.site/pkg/model"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func Login(loginRequest model.LoginRequest) (model.LoginResponse, error) {
	log.Print("Finding User...")
	user, err := FindUser(db.DB,loginRequest.UserName, loginRequest.Password)
	if err != nil {
		return model.LoginResponse{}, err
	}
	log.Print("Encrypting token..")
	token := encryptToken(user)
	if token == "" {
		return model.LoginResponse{}, errors.New("unable to perform login")
	}
	log.Print("token encrypted.")
	err = saveToken(db.DB, token)
	if err != nil {
		return model.LoginResponse{}, err
	}
	log.Printf("Token saved to db.")
	return model.LoginResponse{
	AccessToken: 	token,
	}, nil
}


func saveToken(db *sql.DB, token string) error {
	const insertQuery = "INSERT INTO oauth_access_tokens(token, created_at, is_revoked) VALUES (?,?,?)"
	_, err := db.Exec(insertQuery,token, time.Now().Unix(),0)
	if err != nil {
		return err
	}
	return nil
}

func FindUser(db *sql.DB, userName string, password string) (model.User, error) {
	const query =
		"SELECT u.id as user_id," +
			"u.user_name as user_name," +
			"u.password as password," +
			"u.name as name," +
			"u.created_at as created_at," +
			"u.role as role from users u where u.user_name =?"
	log.Print("executing query to find user")
	result, err := db.Query(query, userName)
	log.Print("query successful")
	defer func(results *sql.Rows) {
		err := results.Close()
		if err != nil {
			log.Printf("Error while closing resultset, %v", err.Error())
		}
	}(result)
	if err != nil {
		log.Printf("Error while querying user :%v", err.Error())
		return model.User{}, err
	}

	if result.Next() {
		user := model.User{}
		err = result.Scan(&user.ID, &user.UserName, &user.Password, &user.Name, &user.CreatedAtInternal, &user.Role)
		if err != nil {
			log.Print("Error during result scanning")
			return model.User{}, err
		}
		loc, _ := time.LoadLocation("Asia/Kolkata")
		createdAt, err := common.EpochSecondToTime(strconv.Itoa(user.CreatedAtInternal))
		if err != nil {
			return model.User{}, err
		}
		createdAt = createdAt.In(loc)
		user.CreatedAt = createdAt.Format(time.RFC3339)
		// now check password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return model.User{}, err
		}
		return user, nil
	}
	log.Print("No result found")
	return model.User{}, errors.New("unable to authenticate user")
}

func FindToken(token string, db *sql.DB) (bool, error) {
	const query =
		" SELECT" +
		" oat.id as id, oat.token as token,  oat.created_at as created_at " +
		" from oauth_access_tokens oat where oat.token = ?"
	log.Printf("Finding token from db")
	result, err := db.Query(query,token)
	defer func(results *sql.Rows) {
		err := results.Close()
		if err != nil {
			log.Printf("Error while closing result query :%v", err.Error())
		}
	}(result)
	if err != nil {
		log.Printf("Error while querying account by id:%v", err.Error())
		return false, err
	}

	if result.Next() {
		log.Printf("Found token from table")
		return true, nil
	}

	return false, errors.New("unable to find token from database")
}


func VerifyToken(token string) (model.User, error) {
	ok, err := FindToken(token, db.DB)
	if err != nil {
		return model.User{}, err
	}
	if !ok {
		return model.User{}, errors.New("unable to find token in database")
	}
	user, err := decryptToken(token)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func encryptToken(user model.User) string {

	userJsonArray, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("failed to serialize user json : %s", err.Error())
	}
	encodedBytes := []byte(base64.StdEncoding.EncodeToString(userJsonArray))

	key := []byte(os.Getenv("TOKEN_SECRET"))
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("error creating block cipher from key :: %s\n", err.Error())
		return ""
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
	    fmt.Printf("error setting gcm mode from generated block :: %s", err.Error())
	    return ""
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	    fmt.Printf("error generating the nonce for encryption :: %s", err.Error())
	    return ""
	}

	encryptedBytes := gcm.Seal(nonce, nonce, encodedBytes, nil)

	tokenString := hex.EncodeToString(encryptedBytes)

	return tokenString
}

func decryptToken(token string) (model.User, error) {
	decodedCipherText, err := hex.DecodeString(token)
	if err != nil {
	    fmt.Printf("error decoding token from hex %s", err.Error())
	    return model.User{}, errors.New("error decoding token from hex")
	}

	key := []byte(os.Getenv("TOKEN_SECRET"))
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("error creating block cipher from key :: %s\n", err.Error())
		return model.User{}, errors.New("error creating block cipher from key")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
	    fmt.Printf("error setting gcm mode from generated block :: %s", err.Error())
		return model.User{}, errors.New("error setting gcm mode from generated block")
	}
	decryptedData, err := gcm.Open(nil, decodedCipherText[:gcm.NonceSize()], decodedCipherText[gcm.NonceSize():], nil)
	if err != nil {
	    fmt.Println("error decrypting data", err)
	    return model.User{}, errors.New("error decrypting token")
	}
	// decrypted data is base64 encoded. Decode to get json
	decodedJson := make([]byte, len(decryptedData))
	_, err = base64.StdEncoding.Decode(decodedJson, decryptedData)
	if err != nil {
		log.Printf("error in decoding json %s", err.Error())
		return model.User{}, err
	}

	var user model.User
	// need to do this below, because when we are stroing the bytes in db, it is adding a null character as well in the end.
	decodedJson = bytes.TrimRight(decodedJson, "\x00")
	err = json.Unmarshal(decodedJson, &user)
	if err != nil {
		log.Printf("error in converting json %s", err.Error())
		return model.User{}, err
	}
	return user,nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash),nil
}