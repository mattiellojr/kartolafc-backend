package kartolafc

import (
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sort"
)

type AtletaRanking struct {
	AtletaId int `json:"atleta_id"`
}

type AtletasRanking struct {
	Pontuacao float32
	Atletas []AtletaRanking `bson:"atletas"`
	TimeCompleto struct{
		TimeId int
	} `bson:"timecompleto"`
}

type TimesRanking []AtletasRanking

type TimeRankingFormated struct {
	TimeId int `json:"time_id"`
	Pontuacao float32 `json:"pontuacao"`
	Atletas []AtletaRanking `bson:"atletas" json:"atletas,omitempty"`
}

// para melhor apresentacao no endpoint
type TimesRankingFormated []TimeRankingFormated

// para apresentar a posicao atraves do ID/Slug
type TimeIdRanking struct {
	TimeId int `json:"time_id"`
	Pontuacao float32 `json:"pontuacao"`
	Posicao int `json:"posicao"`
}

func (a TimesRankingFormated) Len() int {
	return len(a)
}

func (a TimesRankingFormated) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a TimesRankingFormated) Less(i, j int) bool {
	a[i].Pontuacao = SomaPontuacao(a[i])
	a[j].Pontuacao = SomaPontuacao(a[j])
	return a[i].Pontuacao > a[j].Pontuacao
}

func SomaPontuacao(atletasTime TimeRankingFormated) float32 {
	var soma float32
	for _, a := range atletasTime.Atletas {

		// se o jogador nao existe no map pode retornar um erro
		if _, ok := CachePontuados.Atletas[a.AtletaId]; ok {
			soma+= CachePontuados.Atletas[a.AtletaId].Pontuacao
		}
	}
	return soma
}

func LoadInMemory(collection *mgo.Collection) {
	inicio := time.Now()

	var atl TimesRanking
	err := collection.Find(bson.M{}).All(&atl)

	if err != nil {
		panic(err)
	}
	log.Println("pegar todos times da collection", time.Since(inicio))

	// formatando os dados e criando array de tamamnho especifico
	atletasFormatado := make([]TimeRankingFormated, len(atl))
	for k, a := range atl {
		timeTemp := TimeRankingFormated{}
		timeTemp.TimeId = a.TimeCompleto.TimeId
		timeTemp.Pontuacao = a.Pontuacao
		timeTemp.Atletas = a.Atletas
		atletasFormatado[k] = timeTemp
	}
	log.Println("atualizado array de times com a posicao no indice")

	go UpdatePontuados()
	// Aguarda o cache dos pontuados serem carregados da api
	time.Sleep(5 * time.Second)
	go SortPontuados(atletasFormatado)
}

func SortPontuados(times TimesRankingFormated) {
	CacheRankingTimeIdPontuados = make([]TimeIdRanking, 15000000)
	for {
		inicio := time.Now()
		sort.Sort(times)
		CacheRankingPontuados = times
		log.Println("sort", time.Since(inicio))

		// outro array de times, porem a chave e o time_id
		for k, t := range times {
			timeTemp := TimeIdRanking{}
			timeTemp.Pontuacao = t.Pontuacao
			timeTemp.Posicao = (k+1)
			timeTemp.TimeId = t.TimeId
			CacheRankingTimeIdPontuados[t.TimeId] = timeTemp
		}
		log.Println("atualizado array de times com time_id no indice")

		melhores()
		//aguarda 2 minutos para fazer o sort novamente
		time.Sleep(120 * time.Second)
	}
}

func melhores() {
	melhores := make([]TimeRankingFormated, 100)
	if len(CacheRankingPontuados) >= 100 {
		for k, temp := range CacheRankingPontuados[:100] {
			melhores[k] = temp
			melhores[k].Atletas = nil
		}
		CacheRankingPontuadosMelhores = melhores
	}
}