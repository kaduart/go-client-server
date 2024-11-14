package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

type Cota struct {
	Dolar string `json:"dolar"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// preparo o request para fazer chamada no WS
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/", nil)
	if err != nil {
		log.Println("Error Reqquest")
	}

	// request ao server
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("decodificou response - CLIENT", string(body))

	//criar arquivo
	a, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal(err)
	}

	objCotacaoDolar := Cota{
		Dolar: string(body),
	}

	//converter struct para json
	data, err := json.Marshal(objCotacaoDolar)
	if err != nil {
		log.Fatal(err)
	}

	//inserindo dados no arquivo
	_, err = a.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	a.Close()

	arquivo, err := os.ReadFile("cotacao.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Escreveu %d bytes no arquivo\n", string(arquivo))

	io.Copy(os.Stdout, res.Body)
}
