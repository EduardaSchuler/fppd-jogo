// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"sync"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo   rune
	cor       Cor
	corFundo  Cor
	tangivel  bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            	 [][]Elemento // grade 2D representando o mapa
	PosX, PosY           int          // posição atual do personagem
	UltimoVisitado	 	 Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg    	     string       // mensagem para a barra de status

	Vida 			     int		  // vida inicial do personagem
	MovimentosPersonagem int		  // contador de movimentos do personagem para limitar os movimentos do inimigo
	TemChave			 bool		  // verifica se o personagem pegou a chave
	PortalAtivo			 bool		  // variavel que verifica se o personagem pegou a chave
	MissaoAdquirida		 bool		  // verifica se o personagem sabe que deve encontrar a chave

	mu 				 	 sync.RWMutex //Mutex
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}

	NPC		   = Elemento{'⚉', CorCinzaEscuro, CorPadrao, true}
	Portal	   = Elemento{'✷', CorAzul, CorPadrao, true}
	Chave	   = Elemento{'⚵', CorVerde ,CorPadrao, true}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	j:= Jogo{
		UltimoVisitado: Vazio,
		Vida: 5,
		MovimentosPersonagem: 0,
		TemChave: false,
		PortalAtivo: false,
		MissaoAdquirida: false,
	}
	j.StatusMsg = fmt.Sprintf("Você começou o jogo com %d de vida!", j.Vida)
	return j
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case NPC.simbolo:
				e = NPC
			case Portal.simbolo:
				e = Portal
			case Chave.simbolo:
				e = Chave
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.mu.Lock()
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
				jogo.mu.Unlock()
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.mu.Lock()
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		jogo.mu.Unlock()
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	jogo.mu.Lock()
	defer jogo.mu.Unlock()

	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	jogo.mu.Lock()
	defer jogo.mu.Unlock()

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado     // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]   // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento              // move o elemento
}

// Move inimigo em direção ao personagem
func inimigoMover(jogo *Jogo) {	
	// Só faz o movimento do inimigo a cada 3 movimentos do personagem
	if jogo.MovimentosPersonagem % 2 == 0 {
		for y := range jogo.Mapa {
			for x := range jogo.Mapa[y] {
				if jogo.Mapa[y][x] == Inimigo {
					dx, dy := 0, 0

					if jogo.PosX > x {
						dx = 1
					} else if jogo.PosX < x {
						dx = -1
					}
					if jogo.PosY > y {
						dy = 1
					} else if jogo.PosY < y {
						dy = -1
					}

					nx, ny := x+dx, y+dy

					// Se nova posição for o personagem causa dano
					if nx == jogo.PosX && ny == jogo.PosY {
						jogo.Vida--
						if jogo.Vida <= 0 {
							jogo.StatusMsg = "Game Over!"
						} else {
							jogo.StatusMsg = fmt.Sprintf("Você foi atingido! Vida restante: %d", jogo.Vida)
						}
						return
					}

					// Move o inimigo se possível
					if jogoPodeMoverPara(jogo, nx, ny) {
						jogoMoverElemento(jogo, x, y, dx, dy)
					}
					return
				}
			}
		}
	}
}

func ativarPortal(jogo *Jogo) {
	for y := range jogo.Mapa {
		for x := range jogo.Mapa[y] {
			if jogo.Mapa[y][x] == Portal {
				for jogo.PortalAtivo {
					jogo.Mapa[y][x].cor = CorCinzaEscuro
					interfaceDesenharJogo(jogo)
					sleep()
					jogo.Mapa[y][x].cor = CorAzul
					interfaceDesenharJogo(jogo)
					sleep()

					if !jogo.PortalAtivo {
						jogo.Mapa[y][x].cor = CorCinzaEscuro
					}
				}
			}
		}
	}
}

func sleep() {
	// Pequeno delay para efeito visual
	time.Sleep(300 * time.Millisecond)
}
