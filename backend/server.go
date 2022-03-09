package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/futurepaul/bakeshop/backend/bakedgood"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/rs/cors"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const fileName = "bakedgoods_sqlite.db"

type serverLnd struct {
	lndClient lnrpc.LightningClient
	db        *bakedgood.SQLiteRepository

	lndURL string
	lndTLS []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ServerConfig struct {
	LndClient lnrpc.LightningClient
	LndURL    string
	LndTLS    []byte
}

func createServer(cfg *ServerConfig) (*http.Server, error) {
	// TODO we should get rid of this but it's nice for clean testing
	// (it wipes the db every time we load up)
	os.Remove(fileName)

	sqldb, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	bakedgoodRepository := bakedgood.NewSQLiteRepository(sqldb)

	if err := bakedgoodRepository.Migrate(); err != nil {
		log.Fatal(err)
	}

	s := serverLnd{
		lndClient: cfg.LndClient,
		db:        bakedgoodRepository,
		lndURL:    cfg.LndURL,
		lndTLS:    cfg.LndTLS,
	}

	// Start the macaroon intercepting
	err = s.createGrpcInterceptor(cfg.LndClient)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.HandleFunc("/bake", s.Bake)
	r.HandleFunc("/cancel", s.Cancel)
	r.HandleFunc("/list", s.List)
	r.HandleFunc("/details/{id}", s.Details)

	handler := cors.Default().Handler(r)
	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: handler,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Println(err.Error())

			// close the db
			// TODO apparently sqlite doesn't need to be "closed"?
			// defer s.db.Close()
		}
	}()

	return httpServer, nil
}

type BakeReq struct {
	Name     string `json:"name"`
	Interval uint64 `json:"interval"`
	Amount   uint64 `json:"amount"`
	Times    uint64 `json:"times"`
}

type BakeResp struct {
	Id string `json:"id"`
}

func (s *serverLnd) Bake(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// Parse request
	var bakeReq BakeReq
	err := json.NewDecoder(r.Body).Decode(&bakeReq)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("Req: Bake")

	macId, macStr, err := bakeMacaroon(s.lndClient, bakeReq)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bg := bakedgood.BakedGood{
		Id:       macId,
		Name:     bakeReq.Name,
		Macaroon: macStr,
		Status:   "active",
	}

	_, err = s.db.CreateBakedGood(bg)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&BakeResp{
		Id: macId,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

type CancelReq struct {
	Id string `json:"id"`
}

func (s *serverLnd) Cancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// Parse request
	var cancelReq CancelReq
	err := json.NewDecoder(r.Body).Decode(&cancelReq)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if cancelReq.Id == "" {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Printf("Req: Cancel  - id: %s\n", cancelReq.Id)

	// Get the macaroon first
	// TODO production - this should probably be in a lock
	macaroon, err := s.db.GetBakedGoodByUuid(cancelReq.Id)
	if err != nil {
		fmt.Println("macaroon not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO do this with sql
	fmt.Printf("found macaroon to cancel: %s\n", macaroon)

	macaroon.Status = "cancelled"

	if _, err := s.db.UpdateBakedGood(macaroon.Id, *macaroon); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Canceled macaroon")

	w.WriteHeader(http.StatusOK)
}

type macaroonKV struct {
	Id    []byte
	Value []byte
}

type MacaroonDetails struct {
	Name      string    `json:"name"`
	Id        string    `json:"id"`
	Macaroon  string    `json:"macaroon"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type MacaroonListResponse struct {
	Items []MacaroonDetails `json:"items"`
}

func (s *serverLnd) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	fmt.Println("Req: List")

	bgList, err := s.db.AllBakedGoods()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	macaroonResponse := &MacaroonListResponse{
		Items: make([]MacaroonDetails, 0),
	}
	for _, bg := range bgList {
		macaroonResponse.Items = append(macaroonResponse.Items, MacaroonDetails{
			Id:        bg.Id,
			Name:      bg.Name,
			Macaroon:  bg.Macaroon,
			Status:    bg.Status,
			CreatedAt: bg.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(macaroonResponse)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type DetailsResponse struct {
	Name       string    `json:"name"`
	Id         string    `json:"id"`
	Macaroon   string    `json:"macaroon"`
	Status     string    `json:"status"`
	LndConnect string    `json:"lndConnect"`
	CreatedAt  time.Time `json:"created_at"`
	// TODO created
	// TODO payments
}

func (s *serverLnd) Details(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// Parse request
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Printf("Req: Details - id: %s\n", id)

	// Get the  macaroon detail
	macaroon, err := s.db.GetBakedGoodByUuid(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Return response
	resp := DetailsResponse{
		Name:       macaroon.Name,
		Id:         id,
		Macaroon:   macaroon.Macaroon,
		Status:     macaroon.Status, // TODO
		LndConnect: fmt.Sprintf("lndconnect://%s?cert=%s&macaroon=%s", s.lndURL, base64.StdEncoding.EncodeToString(s.lndTLS), macaroon.Macaroon),
		CreatedAt:  macaroon.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
