package main

import (
	"fmt"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup
const tempoCorrida = time.Duration(1)

func corredor (recebeBastao, entregaBastao chan time.Duration, idCorredor int){
	// Recebe por meio do canal a duração da corrida atual
	tempo := <- recebeBastao

	fmt.Println("Corredor", idCorredor, "com posse do bastão e correndo...")

	// Simula o período da corrida
	time.Sleep(tempo * time.Second)

	if idCorredor != 4 {
		fmt.Println("Corredor", idCorredor, "terminou a corrida e está entregando o bastão para o corredor", idCorredor+1)
	} else {
		fmt.Println("Corredor", idCorredor, "terminou a corrida")
	}

	// Passa o bastão para o próximo corredor ao escrever no canal o tempo de duração da sua corrida
	entregaBastao <- tempoCorrida
	waitGroup.Done()
}


func main () {
	// Canais para simular a passagem do bastão entre os corredores e repassar a duração de cada corrida
	largada, umParaDois, doisParaTres, tresParaQuatro, chegada := make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration)

	// Goroutine para permitir início da corrida do primeiro corredor
	go func() {
		largada <- tempoCorrida
	}()

	// Goroutine para permitir a finalização da corrida do último corredor
	go func () {
		<- chegada
	}()

	waitGroup.Add(4)
	go corredor(largada, umParaDois, 1)
	go corredor(umParaDois, doisParaTres, 2)
	go corredor(doisParaTres, tresParaQuatro, 3)
	go corredor(tresParaQuatro, chegada, 4)
	waitGroup.Wait()

	fmt.Println("Corrida finalizada!")

}