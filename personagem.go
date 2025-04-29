// personagem.go - Funções para movimentação e ações do personagem
package main

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1 // Move para cima
	case 'a': dx = -1 // Move para a esquerda
	case 's': dy = 1  // Move para baixo
	case 'd': dx = 1  // Move para a direita
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
		jogo.MovimentosPersonagem++
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	if jogo.Vida <= 0 {
		jogo.StatusMsg = "Você não tem mais forças para interagir..."
		return
	}

	// Verifica os 4 blocos ao redor
	direcoes := []struct{ dx, dy int }{
		{0, -1}, // cima
		{0, 1},  // baixo
		{-1, 0}, // esquerda
		{1, 0},  // direita
	}

	for _, d := range direcoes {
		x, y := jogo.PosX+d.dx, jogo.PosY+d.dy

		if y >= 0 && y < len(jogo.Mapa) && x >= 0 && x < len(jogo.Mapa[y]) {
			elem := jogo.Mapa[y][x]

			switch elem {
			case Chave:
				if jogo.MissaoAdquirida {
					jogo.TemChave = true
					jogo.Mapa[y][x] = Vazio
					jogo.StatusMsg = "Você coletou a chave! Vá até o portal antes que o inimigo o alcance!"
					go func() {
						jogo.CanalChave <- true // envia mensagem ao portal
					}()
				} else {
					jogo.StatusMsg = "Você encontrou uma chave, mas não sabe para que serve."
				}
			
			case Portal:
				if jogo.TemChave {
					jogo.PortalAtivo = false;
					jogo.StatusMsg = "Parabéns! Você abriu o portal e conseguiu escapar em segurança!"
				} else {
					jogo.StatusMsg = "Você precisa da chave para abrir o portal!"
				}
			case NPC:
				jogo.MissaoAdquirida = true
				jogo.StatusMsg = "Olá, jogador! Para você escapar é necessário encontrar a chave para liberar o portal! O inimigo escondeu a chave em meio à vegetação, mas se você olhar com olhos atentos, você conseguirá identificar. Boa sorte!"
			}
		}
	}
}



// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
		inimigoMover(jogo)
	}
	return true // Continua o jogo
}
