package main

import (
	"../camada"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type Sensor struct {
	nome string
	id uint16
	valor int16
}

type Parametros struct {
	temperaturaMin  int16
	temperaturaMax  int16
	umidadeMin 	    uint16
	nivelCO2Min     uint16
}

type Atuador struct {
	nome string
	id uint16
	status uint8 //NO LUGAR DO BOOL - TRUE: 1, FALSE: 0
}

type Estufa struct {
	parametrosIni Parametros
	sensores [] Sensor
	atuadores [] Atuador
}

// DECLARAÇOES VARIAVEIS GLOBAIS DO SERVIDOR
var estufa Estufa
//-----------

func checkError(err error, msg string){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro em " + msg, err.Error())
		os.Exit(1)
	}
}

func main()  {
	presetSensores()
	tcpAddrSensor, err := net.ResolveTCPAddr("tcp", ":3300")
	checkError(err, "ResolveTCPAddr")

	fmt.Println("-------- ESTUFA INICIADA ---------")
	fmt.Println("As leituras estão sendo enviadas...")

	go connRetornaSensoresInfo(tcpAddrSensor)

	for {
		continue
	}
}

func connDisparoDeAtuador(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		result, err := ioutil.ReadAll(conn)
		fmt.Println(string(result))
		conn.Close()
	}
}

func presetSensores() {
	//DEFININDO VALORES INICIAIS PARA A STRUCT SENSORES
	var temperatura Sensor
	temperatura.nome = "Temperatura"
	temperatura.id = 1
	temperatura.valor = -15

	var umidade Sensor
	umidade.nome = "Umidade do Solo"
	umidade.id = 2
	umidade.valor = 400

	var nivelCO2 Sensor
	nivelCO2.nome = "Nível de CO2"
	nivelCO2.id = 3
	nivelCO2.valor = 300

	//ADICIONANDO SENSORES A ESTUFA
	estufa.sensores = append(estufa.sensores, temperatura)
	estufa.sensores = append(estufa.sensores, umidade)
	estufa.sensores = append(estufa.sensores, nivelCO2)
	//----------------
}

func presetAtuadores() {
	//DEFININDO VALORES INICIAIS PARA A STRUCT ATUADORES
	var aquecedor Atuador
	aquecedor.nome = "Aquecedor"
	aquecedor.id = 10
	aquecedor.status = 0

	var resfriador Atuador
	resfriador.nome = "Resfriador"
	resfriador.id = 20
	resfriador.status = 0

	var irrigador Atuador
	irrigador.nome = "Irrigador"
	irrigador.id = 30
	irrigador.status = 0

	var injetorCO2 Atuador
	injetorCO2.nome = "injetorCO2"
	injetorCO2.id = 40
	injetorCO2.status = 0

	//ADICIONANDO SENSORES A ESTUFA
	estufa.atuadores = append(estufa.atuadores, aquecedor)
	estufa.atuadores = append(estufa.atuadores, resfriador)
	estufa.atuadores = append(estufa.atuadores, irrigador)
	estufa.atuadores = append(estufa.atuadores, injetorCO2)
}

func connRetornaSensoresInfo(addr *net.TCPAddr) {
	for {
		conn, err := net.DialTCP("tcp", nil, addr)
		checkError(err, "DialTCP")

		simulaMudancas()
		var buffer bytes.Buffer

		buffer = converteSensoresEmArrayDeBytes(estufa.sensores, buffer)

		var pacote = gopacket.NewPacket(
			buffer.Bytes(),
			camada.SensoresLayerType,
			gopacket.Default,
		)
		conn.Write(pacote.Data())

		msg := make([]byte, 40)
		conn.Read(msg[:])
		if msg[:] != nil {
			fmt.Println(string(msg))
		}


		conn.Close()
		time.Sleep(4 * time.Second)
	}
}

func simulaMudancas() {
	estufa.sensores[0].valor++
}

func converteSensoresEmArrayDeBytes(sensores []Sensor, buffer bytes.Buffer) bytes.Buffer {
	for i := 0; i < 3; i++ {
		var nomeBytes = make([]byte, 15)
		for i, j := range []byte(sensores[i].nome) {
			nomeBytes[i] = byte(j)
		}
		var idSensorBytes = make([]byte, 2)
		var valorBytes = make([]byte, 4)

		binary.BigEndian.PutUint16(idSensorBytes, sensores[i].id)
		binary.BigEndian.PutUint32(valorBytes, uint32(sensores[i].valor))

		buffer.Write(nomeBytes)
		buffer.Write(idSensorBytes)
		buffer.Write(valorBytes)
	}

	return buffer
}