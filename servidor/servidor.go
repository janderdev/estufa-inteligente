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

func main() {
	presetSensores()
	presetAtuadores()

	fmt.Println("------- SERVIDOR INICIADO ---------")
	fmt.Println("Aguardando CLIENTE definir parms iniciais...")

	// ELEMENTOS INICIAIS DA CONEXÃO DOS SOCKETS CLIENTE & SERVIDOR
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":3200")
	checkError(err, "ResolveTCPAddr")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")

	// ELEMENTOS INICIAIS DA CONEXÃO DOS SOCKETS ESTUFA & SERVIDOR (SENSORES)
	addrEstufaSensor, err := net.ResolveTCPAddr("tcp", ":3300")
	checkError(err, "ResolveTCPAddr")
	listenerEstufaSensor, err := net.ListenTCP("tcp", addrEstufaSensor)
	checkError(err, "ListenTCP")

	// ELEMENTOS INICIAIS DA CONEXÃO DOS SOCKETS SERVIDOR & CLIENTE (SENSORES)
	//addrEstufaAtuador, err := net.ResolveTCPAddr("tcp", ":3500")
	//checkError(err, "ResolveTCPAddr")
	//
	//connEstufaAtuador, err := net.DialTCP("tcp", nil, addrEstufaAtuador)
	//checkError(err, "DialTCP Atuador")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		connGetParametrosDoCliente(conn)
		break
	}
	fmt.Println("Aguardando estufa.go ser iniciada...")

	go loopGetSensorCliente(listener)
	go loopGetSensoresEstufa(listenerEstufaSensor)

	for {
		continue
	}
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

func presetSensores() {
	var temperatura Sensor
	var umidade Sensor
	var nivelCO2 Sensor

	//ADICIONANDO SENSORES A ESTUFA
	estufa.sensores = append(estufa.sensores, temperatura)
	estufa.sensores = append(estufa.sensores, umidade)
	estufa.sensores = append(estufa.sensores, nivelCO2)
	//----------------
}


func VerificaOsLimitesMinMax(conn net.Conn) {
	tempMin := estufa.parametrosIni.temperaturaMin
	tempMax := estufa.parametrosIni.temperaturaMax
	umidadeMin := estufa.parametrosIni.umidadeMin
	nivelC02min := estufa.parametrosIni.nivelCO2Min
	var update uint8 = 0

	// VERIFICA SENSOR 1 _ ATUADOR 1 - AQUECEDOR
	if estufa.sensores[0].valor < tempMin {
		if estufa.atuadores[0].status == 0 {
			estufa.atuadores[0].status = 1
			connEnviaInfoDoAtuador(conn, estufa.atuadores[0])
			update = 1
		}
	} else {
		if estufa.atuadores[0].status == 1 {
			estufa.atuadores[0].status = 0
			connEnviaInfoDoAtuador(conn, estufa.atuadores[0])
			update = 1
		}
	}

	// VERIFICA SENSOR 1 _ ATUADOR 2 - RESFRIADOR
	if estufa.sensores[0].valor > tempMax {
		if estufa.atuadores[1].status == 0 {
			estufa.atuadores[1].status = 1
			connEnviaInfoDoAtuador(conn, estufa.atuadores[1])
			update = 1
		}
	} else {
		if estufa.atuadores[1].status == 1 {
			estufa.atuadores[1].status = 0
			connEnviaInfoDoAtuador(conn, estufa.atuadores[1])
			update = 1
		}
	}

	// VERIFICA SENSOR 2 _ ATUADOR 3 - IRRIGADOR
	if uint16(estufa.sensores[1].valor) < umidadeMin {
		if estufa.atuadores[2].status == 0 {
			estufa.atuadores[2].status = 1
			connEnviaInfoDoAtuador(conn, estufa.atuadores[2])
			update = 1
		}
	} else {
		if estufa.atuadores[2].status == 1 {
			estufa.atuadores[2].status = 0
			connEnviaInfoDoAtuador(conn, estufa.atuadores[2])
			update = 1
		}
	}

	// VERIFICA SENSOR 2 _ ATUADOR 3 - IRRIGADOR
	if uint16(estufa.sensores[2].valor) < nivelC02min {
		if estufa.atuadores[3].status == 0 {
			estufa.atuadores[3].status = 1
			connEnviaInfoDoAtuador(conn, estufa.atuadores[3])
			update = 1
		}
	} else {
		if estufa.atuadores[3].status == 1 {
			estufa.atuadores[3].status = 0
			connEnviaInfoDoAtuador(conn, estufa.atuadores[3])
			update = 1
		}
	}

	if update == 0 {
		var pacote []byte
		pacote = nil
		conn.Write(pacote)
	}
}


func loopGetSensoresEstufa(listenerEstufa *net.TCPListener) {

	for {
		for {
			conn, err := listenerEstufa.Accept()
			if err != nil {
				return
			}

			connGetSensoresInfoDaEstufa(conn)
			break
		}
		time.Sleep(4 * time.Second)
	}
}

func loopGetSensorCliente(listener *net.TCPListener) {
	for {
		for {
			novaConn, err := listener.Accept()
			if err != nil {
				return
			}

			connRetornaSensorInfo(novaConn)
		}
	}
}

func connGetParametrosDoCliente(conn net.Conn) {
	result, err := ioutil.ReadAll(conn)
	checkError(err, "ReadAll")

	packet := gopacket.NewPacket(
		result,
		camada.ParametersLayerType,
		gopacket.Default,
	)

	decodePacket := packet.Layer(camada.ParametersLayerType)

	if decodePacket != nil {
		fmt.Println("------------ PARAMETROS DEFINIDOS --------------")
		content, _ := decodePacket.(*camada.ParametersLayer)
		fmt.Println("TemperaturaMin:", int16(content.TemperaturaMin))
		fmt.Println("TemperaturaMax:", int16(content.TemperaturaMax))
		fmt.Println("UmidadeMin:", content.UmidadeMin)
		fmt.Println("NivelCO2Min:", content.NivelCO2Min)
		fmt.Println("------------------------------------------------")

		//GUARDADOS PARAMENTROS NA MEMORIA LOCAL DO SEVIDOR
		estufa.parametrosIni.temperaturaMin = int16(content.TemperaturaMin)
		estufa.parametrosIni.temperaturaMax = int16(content.TemperaturaMax)
		estufa.parametrosIni.umidadeMin = content.UmidadeMin
		estufa.parametrosIni.nivelCO2Min = content.NivelCO2Min
	}
	conn.Close()
}

func connGetSensoresInfoDaEstufa(conn net.Conn) {
	result := make([]byte, 63)
	conn.Read(result[:])

	pacote := gopacket.NewPacket(
		result,
		camada.SensoresLayerType,
		gopacket.Default,
	)

	decodePacote := pacote.Layer(camada.SensoresLayerType)

	if decodePacote != nil {
		content, _ := decodePacote.(*camada.SensoresLayer)
		estufa.sensores[0].id = content.Temperatura.IDSensor
		estufa.sensores[0].nome = content.Temperatura.Nome
		estufa.sensores[0].valor = int16(content.Temperatura.Valor)

		estufa.sensores[1].id = content.Umidade.IDSensor
		estufa.sensores[1].nome = content.Umidade.Nome
		estufa.sensores[1].valor =  int16(content.Umidade.Valor)

		estufa.sensores[2].id = content.NivelDeCO2.IDSensor
		estufa.sensores[2].nome = content.NivelDeCO2.Nome
		estufa.sensores[2].valor =  int16(content.NivelDeCO2.Valor)
		fmt.Println("A leitura dos sensores chegaram!")
	}

	VerificaOsLimitesMinMax(conn)
	conn.Close()
}

func connRetornaSensorInfo(conn net.Conn) {
	result := make([]byte, 2)
	conn.Read(result[:])
	valor := binary.BigEndian.Uint16(result)

	var dadosSensor Sensor
	for _, sensor := range estufa.sensores {
		if sensor.id == valor {
			dadosSensor = sensor
		}
	}

	var buffer bytes.Buffer
	buffer = converteSensorEmArrayDeBytes(dadosSensor, buffer)

	pacote := gopacket.NewPacket(
		buffer.Bytes(),
		camada.SensorLayerType,
		gopacket.Default,
	)

	conn.Write(pacote.Data())
	conn.Close()
}

func connEnviaInfoDoAtuador(conn net.Conn, atuador Atuador) {
	msg := []byte("Um atuador foi alterado:" + atuador.nome)
	conn.Write(msg)
}

func converteSensorEmArrayDeBytes(sensor Sensor, buffer bytes.Buffer) bytes.Buffer {

	var nomeBytes = make([]byte, 15)
	for i, j := range []byte(sensor.nome) {
		nomeBytes[i] = byte(j)
	}

	var idBytes = make([]byte, 2)
	var valorBytes = make([]byte, 4)

	binary.BigEndian.PutUint16(idBytes, sensor.id)
	binary.BigEndian.PutUint32(valorBytes, uint32(sensor.valor))

	buffer.Write(nomeBytes)
	buffer.Write(idBytes)
	buffer.Write(valorBytes)

	return buffer
}

func checkError(err error, msg string){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro em " + msg, err.Error())
		os.Exit(1)
	}
}