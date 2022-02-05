package main

import (
	"fmt"
	"sync"
	"time"
)

var wgCorredores, wgEquipes sync.WaitGroup
var mutex sync.Mutex

const tempoCorrida = time.Duration(250)

func corredor (recebeBastao, entregaBastao chan time.Duration, idCorredor, idEquipe int, ordemChegada chan int){
	// Recebe por meio do canal a duração da corrida atual
	tempo := <- recebeBastao

	fmt.Printf("Equipe %d: corredor %d com posse do bastão e correndo...\n", idEquipe, idCorredor)

	// Simula o período da corrida
	time.Sleep(tempo * time.Millisecond)

	if idCorredor != 4 {
		fmt.Printf("Equipe %d: corredor %d terminou a corrida e está entregando o bastão para o corredor %d.\n", idEquipe, idCorredor, idCorredor+1)
	} else {
		// Mutex utilizado para que as goroutines escrevam na saída padrão na mesma ordem em que escreveram no canal
		mutex.Lock()
		ordemChegada <- idEquipe
		fmt.Printf("Equipe %d: corredor %d terminou a corrida.\n", idEquipe, idCorredor)
		mutex.Unlock()
	}

	// Passa o bastão para o próximo corredor ao escrever no canal o tempo de duração da sua corrida
	entregaBastao <- tempoCorrida
	wgCorredores.Done()
}

func equipe (idEquipe int, ordemChegada chan int) {
	// Canais para simular a passagem do bastão entre os corredores da equipe e repassar a duração da corrida
	largada, umParaDois, doisParaTres, tresParaQuatro, chegada := make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration), make(chan time.Duration)

	// Goroutine para permitir início da corrida do primeiro corredor
	go func() {
		largada <- tempoCorrida
	}()

	// Goroutine para permitir a finalização da corrida do último corredor
	go func () {
		<- chegada
	}()

	wgCorredores.Add(4)
	go corredor(largada, umParaDois, 1, idEquipe, ordemChegada)
	go corredor(umParaDois, doisParaTres, 2, idEquipe, ordemChegada)
	go corredor(doisParaTres, tresParaQuatro, 3, idEquipe, ordemChegada)
	go corredor(tresParaQuatro, chegada, 4, idEquipe, ordemChegada)
	wgCorredores.Wait()

	wgEquipes.Done()
}

func main(){
	var qtd int
	fmt.Print("Digite a quantidade de equipes participantes: ")
	fmt.Scan(&qtd)
	fmt.Println()

	// Aloca dinamicamente o canal bufferizado para armazenar a ordem de chegada
	ordemChegada := make(chan int, qtd)

	// Momento de largada
	for i := 1; i <= qtd; i++ {
		wgEquipes.Add(1)
		go equipe(i, ordemChegada)
	}
	wgEquipes.Wait()

	close(ordemChegada)

	fmt.Println("\nOrdem de chegada das equipes:")
	for i := 1; i <= qtd; i++ {
		fmt.Printf("#%d - Equipe %d\n", i, <- ordemChegada)
	}

}