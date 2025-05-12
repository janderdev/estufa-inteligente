# 🌿 Estufa Inteligente com Controle Automatizado
Projeto de uma estufa inteligente desenvolvida em Go, com controle automatizado de sensores, aquecedores e resfriadores para monitoramento e ajuste dinâmico do ambiente.

## Download da imagem
Baixe a imagem remota para a sua máquina:
```
docker pull douglasmoraiis/estufainteligente:latest
```
## Instalação
Execute a imagem baixada em um container:
```
docker run --name estufa -it douglasmoraiis/estufainteligente
```
## Execução da aplicação
Quando o terminal da imagem iniciar, execute o arquivo `servidor.go`:
```
go run servidor/servidor.go
```
Agora abra um novo terminal no seu computador e execute o seguinte comando:
```
docker container exec -it estufa bash
```
Ele vai abrir um novo terminal da mesma imagem que já está em execução.

Execute o arquivo `cliente.go`:
```
go run cliente/cliente.go
```
Quando o cliente executar, preencha os Parâmetros que delimitam os limites da estufa.

Obs.: Depois que os limites forem definidos um novo dialogo é exibido. Por enquanto ele retorna apenas valores vazios, pois a conexão do Servidor com a Estufa ainda não foi estabelecida.

Então, abra um novo terminal e novamente execute o comando:
```
docker container exec -it estufa bash
```
E agora execute o arquivo `estufa.go`:
```
go run estufa/estufa.go
```
Pronto a aplicação está sendo executada e as informações estão sendo trocadas entre o
Cliente, Servidor e a Estufa.

O terminal referente ao Cliente agora retorna o valor do sensor atualizado, quando requisitado.
