package main

import (
	"time"
)

func ativarPortal(jogo *Jogo) {
	select {
	case <-jogo.CanalChave:
		jogo.PortalAtivo = true
		for y := range jogo.Mapa {
			for x := range jogo.Mapa[y] {
				if jogo.Mapa[y][x] == Portal {
					go func(x, y int) {
						for jogo.PortalAtivo {
							jogo.Mapa[y][x].cor = CorCinzaEscuro
							interfaceDesenharJogo(jogo)
							sleep()
							jogo.Mapa[y][x].cor = CorAzul
							interfaceDesenharJogo(jogo)
							sleep()
						}
						jogo.Mapa[y][x].cor = CorCinzaEscuro
					}(x, y)
				}
			}
		}
	}
}

func sleep() {
	// Pequeno delay para efeito visual
	time.Sleep(300 * time.Millisecond)
}