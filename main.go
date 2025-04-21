// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

func iniciarMovimentoInimigos(jogo *Jogo) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			inimigoMover(jogo)
		}
	}()
}

func iniciarCronometro(jogo *Jogo, encerrar chan bool) {
	go func() {
		timer := time.NewTimer(60 * time.Second)
		<-timer.C
		if jogo.Vida > 0 {
			jogo.StatusMsg = "Tempo esgotado! Você não conseguiu escapar!"
			jogo.Vida = 0
		}
		encerrar <- true
	}()
}


func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Canal para encerrar o jogo
	encerrar := make(chan bool)
	iniciarCronometro(&jogo, encerrar)


	// Loop principal de entrada
	rodando := true
	for rodando {
		select {
		case <-encerrar:  // Se o canal encerrar for acionado (tempo esgotado)
			rodando = false  // Definimos rodando como false para sair do loop
		default:
			evento := interfaceLerEventoTeclado()
			if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
				rodando = false  // Caso o jogador saia ou o evento de fim de jogo aconteça, rodando também se torna false
			}
			interfaceDesenharJogo(&jogo)
		}
	}
	

	iniciarMovimentoInimigos(&jogo)
}