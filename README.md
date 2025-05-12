# 🌿 Estufa Inteligente com Controle Automatizado

Este projeto consiste em uma **estufa inteligente desenvolvida em Go**, com automação completa de sensores, aquecedores e resfriadores, possibilitando o **monitoramento dinâmico e o ajuste automático do ambiente**.

## 📦 Download da Imagem Docker

Faça o download da imagem remota para sua máquina com o comando abaixo:

```
docker pull douglasmoraiis/estufainteligente:latest
```

## 🚀 Instalação

Execute a imagem baixada em um container:

```
docker run --name estufa -it douglasmoraiis/estufainteligente
```

## 🧠 Execução da Aplicação

### 1. Inicie o servidor

Após o terminal do container ser iniciado, execute o servidor:

```
go run servidor/servidor.go
```

### 2. Inicie o cliente

Abra um novo terminal em seu computador e execute:

```
docker container exec -it estufa bash
```

Dentro do container, rode o cliente:

```
go run cliente/cliente.go
```

Durante a execução, o cliente solicitará os **parâmetros que definem os limites da estufa** (ex: temperatura mínima e máxima).

> ⚠️ Após o preenchimento dos limites, um novo diálogo será exibido. Nesse ponto, os valores ainda estarão vazios, pois a conexão entre o servidor e a estufa ainda não foi estabelecida.

### 3. Inicie a estufa

Em outro terminal, execute novamente o comando:

```
docker container exec -it estufa bash
```

E depois:

```
go run estufa/estufa.go
```

✅ **Pronto!** A aplicação está em funcionamento e a comunicação entre **Cliente**, **Servidor** e **Estufa** está ativa.

O terminal do cliente agora exibirá o valor do sensor atualizado sempre que requisitado.

## 🛠️ Pré-requisitos

Antes de iniciar o projeto, certifique-se de ter os seguintes itens instalados:

- [Go](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.
