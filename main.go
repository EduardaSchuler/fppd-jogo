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
		// Atualiza a mensagem de status com a mensagem de "tempo esgotado"
		if jogo.Vida > 0 {
			jogo.StatusMsg = "Tempo esgotado! Você perdeu sua chance de escapar."
			jogo.Vida = 0
		}
		// Envia sinal de encerramento para o canal, mas o jogo continuará rodando
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

	// Cria o canal de encerramento
	encerrar := make(chan bool)

	// Inicia o cronômetro e o movimento dos inimigos
	iniciarCronometro(&jogo, encerrar)
	iniciarMovimentoInimigos(&jogo)

	// Loop principal de entrada
	rodando := true
	for rodando {
		select {
		case <-encerrar:
			// Se o canal for sinalizado, encerra o jogo, mas só após exibir a mensagem
			rodando = false
		default:
			// Loop de interação com o teclado
			evento := interfaceLerEventoTeclado()
			if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
				rodando = false
			}
			interfaceDesenharJogo(&jogo)
		}
	}
}
