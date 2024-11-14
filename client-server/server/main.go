package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type CotacaoResponse struct {
	USDBRL Cotacao
}

type Cotacao struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/", handlerCotar)
	http.ListenAndServe(":8080", nil)
}

func handlerCotar(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	fmt.Println("Request Iiniciated")

	select {
	case <-time.After(200 * time.Millisecond):

		_, err := GetCotacoes(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case <-ctx.Done():
		fmt.Println("Request cancelled for client or timeout...")
		return
	}
}

func GetCotacoes(w http.ResponseWriter, r *http.Request) (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("Error Reqquest")

		return nil, err
	}

	//executo request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// ler o response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error LER response")
		panic(err)
	}

	// decodifica o json
	var cotacao CotacaoResponse
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		fmt.Println("Error DEcodificar response")
		panic(err)
	}

	//retorna cota json para client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao.USDBRL.Bid)

	// persistir BD as infos 10ms
	newCotacao := NewCotacao(cotacao.USDBRL)
	err = addProduct(*newCotacao)
	if err != nil {
		return nil, err
	}
	return &cotacao.USDBRL, nil

}

func NewCotacao(cotacao Cotacao) *Cotacao {
	return &Cotacao{
		ID:        int(uuid.New().ID()),
		Code:      cotacao.Code,
		Codein:    cotacao.Codein,
		Name:      cotacao.Name,
		High:      cotacao.High,
		Low:       cotacao.Low,
		VarBid:    cotacao.VarBid,
		PctChange: cotacao.PctChange,
		Bid:       cotacao.Bid,
		Ask:       cotacao.Ask,
		Timestamp: cotacao.Timestamp,
	}
}

func addProduct(cotacao Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	//inicia conexao com BD
	db, err := sql.Open("sqlite3", "./t1.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Criar a tabela se ela não existir
	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacoes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        code TEXT,
        codein TEXT,
        name TEXT,
		high TEXT,
        low TEXT,
        varBid TEXT,
        pctChange TEXT,
        bid TEXT,
		ask TEXT,
        timestamp TEXT,
        create_date TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.PrepareContext(ctx, "insert into cotacoes(id, code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) values(?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("ERROR:", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cotacao.ID, cotacao.Code, cotacao.Codein, cotacao.Name, cotacao.High, cotacao.Low, cotacao.VarBid, cotacao.PctChange, cotacao.Bid, cotacao.Ask, cotacao.Timestamp, time.Now())
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	fmt.Println("INSERIU DADOS NO BD SUCCESS - SERVER: ", cotacao)

	return nil
}
