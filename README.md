# Jogo de Terminal em Go

Este projeto é um pequeno jogo desenvolvido em Go que roda no terminal usando a biblioteca [termbox-go](https://github.com/nsf/termbox-go). O jogador controla um personagem que pode se mover por um mapa carregado de um arquivo de texto.

## Como funciona

- O mapa é carregado de um arquivo `.txt` contendo caracteres que representam diferentes elementos do jogo.
- O personagem se move com as teclas **W**, **A**, **S**, **D**.
- Pressione **E** para interagir com o ambiente.
- Pressione **ESC** para sair do jogo.

### Controles

| Tecla | Ação              |
|-------|-------------------|
| W     | Mover para cima   |
| A     | Mover para esquerda |
| S     | Mover para baixo  |
| D     | Mover para direita |
| E     | Interagir         |
| ESC   | Sair do jogo      |

## Como compilar

1. Instale o Go e clone este repositório.
2. Inicialize um novo módulo "jogo":

```bash
go mod init jogo
go get -u github.com/nsf/termbox-go
```

3. Compile o programa:

Linux:

```bash
go build -o jogo
```

Windows:

```bash
go build -o jogo.exe
```

Também é possivel compilar o projeto usando o comando `make` no Linux ou o script `build.bat` no Windows.

## Como executar

1. Certifique-se de ter o arquivo `mapa.txt` com um mapa válido.
2. Execute o programa no termimal:

```bash
./jogo
```

## Estrutura do projeto

- main.go — Ponto de entrada e loop principal
- interface.go — Entrada, saída e renderização com termbox
- jogo.go — Estruturas e lógica do estado do jogo
- personagem.go — Ações do jogador

## Implementação

*Cor personalizada para a chave*
A chave que antes usava uma cor padrão agora foi alterada para a cor amarela, tornando sua identificação mais intuitiva para o jogador.

*Leitura de mapa por arquivo .txt*
O jogo utiliza um sistema de leitura de mapas a partir de arquivos de texto (mapa.txt e maze.txt), permitindo fácil customização do ambiente do jogo.

*Sistema de interação*
O jogador pode interagir com elementos do mapa pressionando a tecla E, como por exemplo coletar a chave.

*Detecção de vitória*
Quando o jogador coleta a chave e alcança a saída, o jogo exibe uma mensagem de vitória, encerrando a partida.

*Elementos Adicionais*
'⚉' -> NPC
'✷' -> Portal
'⚵' -> Chave