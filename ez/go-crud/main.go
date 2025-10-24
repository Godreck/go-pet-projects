package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Game struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Genre       string     `json:"genre"`
	ReleaseDate time.Time  `json:"releaseDate"`
	DevStudio   *DevStudio `json:"devStudio"`
	Publisher   string     `json:"Publisher"`
	Budget      int        `json:"budget"`
}

type DevStudio struct {
	StudioName string    `json:"studioName"`
	Director   *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"LastName"`
}

func getGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range games {
		if strconv.Itoa(item.ID) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

}

func getGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content/type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func mkGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content/type", "application/json")
	var game Game
	_ = json.NewDecoder(r.Body).Decode(&game)
	game.ID = rand.Intn(1000000)
	games = append(games, game)
	json.NewEncoder(w).Encode(games)
}

func uptGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content/type", "application/json")
	params := mux.Vars(r)
	var game Game
	_ = json.NewDecoder(r.Body).Decode(&game)
	for idx, item := range games {
		if strconv.Itoa(item.ID) == params["id"] {
			cur_game := &games[idx]
			cur_game.Name = game.Name
			cur_game.Genre = game.Genre
			cur_game.Budget = game.Budget
			json.NewEncoder(w).Encode(games[idx])
			return
		}
	}

}

func delGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content/type", "application/json")
	params := mux.Vars(r)
	for i, item := range games {
		if strconv.Itoa(item.ID) == params["id"] {
			games = append(games[:i], games[i+1:]...)
			return
		}
	}
	json.NewEncoder(w).Encode(games)
}

var games []Game

func main() {
	games = append(games, Game{
		ID:    1,
		Name:  "Project: Ahra",
		Genre: "ActionRPG",
		ReleaseDate: time.Date(
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			time.Now().Hour(),
			time.Now().Minute(),
			time.Now().Second(),
			time.Now().Nanosecond(),
			time.Now().Location(),
		),
		DevStudio: &DevStudio{
			StudioName: "Reveil",
			Director: &Director{
				FirstName: "Aaron",
				LastName:  "Right",
			},
		},
		Budget: 100,
	})

	games = append(games, Game{
		ID:    1,
		Name:  "Project: Ahra - Part II",
		Genre: "ActionRPG",
		ReleaseDate: time.Date(
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			time.Now().Hour(),
			time.Now().Minute(),
			time.Now().Second(),
			time.Now().Nanosecond(),
			time.Now().Location(),
		),
		DevStudio: &DevStudio{
			StudioName: "Reveil",
			Director: &Director{
				FirstName: "Aaron",
				LastName:  "Right",
			},
		},
		Budget: 200,
	})

	router := mux.NewRouter()
	router.HandleFunc("/Games", getGames).Methods("GET")
	router.HandleFunc("/Games/{id}", getGame).Methods("GET")
	router.HandleFunc("/Games", mkGame).Methods("POST")
	router.HandleFunc("/Games/{id}", uptGame).Methods("PUT")
	router.HandleFunc("/Games/{id}", delGame).Methods("DELETE")
	router.HandleFunc("/", getGames)
	http.Handle("/", router)

	fmt.Println("Starting mux server at port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
