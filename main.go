package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"pacman/utils"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User User構造体
type User struct {
	ID        int
	FirstName string
	LastName  string
}

type Door struct {
	Door string `json:"door"`
}
type Mode struct {
	Mode string `json:"mode"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	// router.HandleFunc("/users", findAllUsers)
	// router.HandleFunc("/users/{id}", findByID)
	router.HandleFunc("/door", findAllDoors).Methods("GET")
	router.HandleFunc("/mode", findAllModes).Methods("GET")

	router.HandleFunc("/door", createDoor).Methods("POST")
	router.HandleFunc("/mode", createMode).Methods("POST")

	log.Fatal(http.ListenAndServe(":80", router))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func findAllDoors(w http.ResponseWriter, r *http.Request) {
	// DB接続
	db := utils.GetConnection()
	defer db.Close()

	var doorList []Door
	var door_res Door

	// db.Table("door").Find(&doorList)
	db.Raw(`select door from door where created > current_timestamp + interval - 30 second`).Scan(&doorList)

	fmt.Println(len(doorList))
	if len(doorList) == 0 {
		rand.Seed(time.Now().UnixNano())
		doorsChoice := []string{"red", "orange", "yellow", "green", "cyan", "blue", "purple"}
		door_res.Door = doorsChoice[rand.Intn(7)]

	} else {
		m := map[string]int{"red": 0, "orange": 0, "yellow": 0, "green": 0, "cyan": 0, "blue": 0, "purple": 0}
		for i := 0; i < len(doorList); i++ {
			m[doorList[i].Door]++
		}
		var max_key string
		var max_value int
		max_value = -1

		for k, v := range m {
			if max_value < v {
				max_value = v
				max_key = k
			}

			fmt.Printf("key: %s, value: %d\n", k, v)
		}
		door_res.Door = max_key

	}

	// 共通化した処理を使う
	utils.RespondWithJSON(w, http.StatusOK, door_res)
}

func findAllModes(w http.ResponseWriter, r *http.Request) {
	// DB接続
	db := utils.GetConnection()
	defer db.Close()

	var modeList []Mode
	var mode_res Mode

	// db.Table("door").Find(&doorList)
	db.Raw(`select mode from mode where created > current_timestamp + interval - 30 second`).Scan(&modeList)

	fmt.Println(len(modeList))
	if len(modeList) == 0 {
		rand.Seed(time.Now().UnixNano())
		modesChoice := []string{"light", "dark"}
		mode_res.Mode = modesChoice[rand.Intn(2)]

	} else {
		m := map[string]int{"light": 0, "dark": 0}
		for i := 0; i < len(modeList); i++ {
			m[modeList[i].Mode]++
		}
		var max_key string
		var max_value int
		max_value = -1

		for k, v := range m {
			if max_value < v {
				max_value = v
				max_key = k
			}

			fmt.Printf("key: %s, value: %d\n", k, v)
		}
		mode_res.Mode = max_key

	}

	// 共通化した処理を使う
	utils.RespondWithJSON(w, http.StatusOK, mode_res)
}

func createDoor(w http.ResponseWriter, r *http.Request) {
	// リクエストボディ取得
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var door Door
	// 読み込んだJSONを構造体に変換
	if err := json.Unmarshal(body, &door); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "JSON Unmarshaling failed .")
		return
	}

	// DB接続
	db := utils.GetConnection()
	defer db.Close()

	// DBにINSERTする
	db.Table("door").Create(&door)

	utils.RespondWithJSON(w, http.StatusOK, door)

}

func createMode(w http.ResponseWriter, r *http.Request) {
	// リクエストボディ取得
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var mode Mode
	// 読み込んだJSONを構造体に変換
	if err := json.Unmarshal(body, &mode); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "JSON Unmarshaling failed .")
		return
	}

	// DB接続
	db := utils.GetConnection()
	defer db.Close()

	// DBにINSERTする
	db.Table("mode").Create(&mode)

	utils.RespondWithJSON(w, http.StatusOK, mode)

}

// /////////////////////SAMPLE////////////////////////
func findAllUsers(w http.ResponseWriter, r *http.Request) {
	// DB接続
	db := utils.GetConnection()
	defer db.Close()

	var userList []User
	db.Find(&userList)

	// 共通化した処理を使う
	utils.RespondWithJSON(w, http.StatusOK, userList)
}

func findByID(w http.ResponseWriter, r *http.Request) {

	id, err := utils.GetID(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid parameter")
		return
	}

	// DB接続
	db := utils.GetConnection()
	defer db.Close()

	var user User
	db.Where("id = ?", id).Find(&user)

	// 共通化した処理を使う
	utils.RespondWithJSON(w, http.StatusOK, user)
}
