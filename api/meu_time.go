package api

import (
	"log"
	"encoding/json"
)

type MeuTime struct {
	Atletas []Atleta `json:"atletas"`
	TimeCompleto Time `json:"time"`
	Patrimonio int `json:"patrimonio"`
	EsquemaId int `json:"esquema_id"`
	ValorTime float32 `json:"valor_time"`
	Mensagem string `json:"mensagem"`
	RodadaAtual int `json:"rodada_atual"`
	Pontos float32 `json:"pontos"`
	Clubes map[int]Clube `json:"clubes"`
	Posicoes map[int]Posicoes `json:"posicoes"`
	Status map[int] struct{
		Id int `json:"id"`
		Nome string `json:"nome"`
	} `json:"status"`
}

const URL_MEU_TIME = "/auth/time"

func (m *MeuTime) Get(token string) {
	request := Request{}
	resp, err := request.GetToken(URL_MEU_TIME, 10, token)

	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(resp.Body(), &m)
}