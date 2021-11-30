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
	status uint16
}

type Estufa struct {
	parametrosIni Parametros
	sensores [] Sensor
	atuadores [] Atuador
}

// DECLARAÇÃO STRUCT GLOBAL DO SERVIDOR - ARMAZENAS OS DADOS DO Parametrs Mín e Max, Sensores e Atuadores
var estufa Estufa
//-----------

func main() {
	presetSensores()
	presetAtuadores()

	fmt.Println("------------ SERVIDOR INICIADO -----------")
	fmt.Println("Aguardando CLIENTE definir os Parametros de limite...")

	// ELEMENTOS INICIAIS DA CONEXÃO DO SOCKET CLIENTE & SERVIDOR
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":3200")
	checkError(err, "ResolveTCPAddr")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")

	// ELEMENTOS INICIAIS DA CONEXÃO DO SOCKET SERVIDOR & SERVIDOR
	addrEstufa, err := net.ResolveTCPAddr("tcp", ":3300")
	checkError(err, "ResolveTCPAddr")
	listenerEstufa, err := net.ListenTCP("tcp", addrEstufa)
	checkError(err, "ListenTCP")

	// PASSAGEM DE PARAMETROS DE LIMITES
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
	go loopGetSensoresEstufa(listenerEstufa)

	for {
		continue
	}
}

func presetAtuadores() {
	//DEFININDO VALORES INICIAIS PARA OS ATUADORES DENTRO DA ESTUFA
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
	injetorCO2.nome = "InjetorCO2"
	injetorCO2.id = 40
	injetorCO2.status = 0

	//ADICIONANDO ATUADORES A ESTUFA
	estufa.atuadores = append(estufa.atuadores, aquecedor)
	estufa.atuadores = append(estufa.atuadores, resfriador)
	estufa.atuadores = append(estufa.atuadores, irrigador)
	estufa.atuadores = append(estufa.atuadores, injetorCO2)
}

func presetSensores() {
	var temperatura Sensor
	var umidade Sensor
	var nivelCO2 Sensor

	//ADICIONANDO SENSORES VAZIOS A ESTUFA
	estufa.sensores = append(estufa.sensores, temperatura)
	estufa.sensores = append(estufa.sensores, umidade)
	estufa.sensores = append(estufa.sensores, nivelCO2)
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
		time.Sleep(1 * time.Second)
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

		//GUARDANDOS PARAMENTROS NA MEMORIA LOCAL DO SEVIDOR
		estufa.parametrosIni.temperaturaMin = int16(content.TemperaturaMin)
		estufa.parametrosIni.temperaturaMax = int16(content.TemperaturaMax)
		estufa.parametrosIni.umidadeMin = content.UmidadeMin
		estufa.parametrosIni.nivelCO2Min = content.NivelCO2Min
	}
	_ = conn.Close()
}

func connGetSensoresInfoDaEstufa(conn net.Conn) {
	result := make([]byte, 63)
	_, err := conn.Read(result[:])
	checkError(err, "Conn.Read/GetSensores")

	// CRIA O PACOTE DOS DADOS RECEBIDOS DE ACORDO COM SensoresLayer EXIGE
	pacote := gopacket.NewPacket(
		result,
		camada.SensoresLayerType,
		gopacket.Default,
	)

	decodePacote := pacote.Layer(camada.SensoresLayerType)

	//GUARDANDOS SENSORES ATUALIZADOS NA MEMORIA LOCAL DO SEVIDOR
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

	//ANTES DE ENCERRAR A CONEXAO --
	verificaOsLimitesMinMax(conn)
	_ = conn.Close()
}

func verificaOsLimitesMinMax(conn net.Conn) {
	tempMin := estufa.parametrosIni.temperaturaMin
	tempMax := estufa.parametrosIni.temperaturaMax
	umidadeMin := estufa.parametrosIni.umidadeMin
	nivelC02min := estufa.parametrosIni.nivelCO2Min
	var atuadoresAlterados [] Atuador

	// VERIFICA SENSOR 1 _ ATUADOR 1 - AQUECEDOR
	if estufa.sensores[0].valor < tempMin {
		if estufa.atuadores[0].status == 0 {
			estufa.atuadores[0].status = 1
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[0])
		}
	} else {
		if estufa.atuadores[0].status == 1 {
			estufa.atuadores[0].status = 0
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[0])
		}
	}

	// VERIFICA SENSOR 1 _ ATUADOR 2 - RESFRIADOR
	if estufa.sensores[0].valor > tempMax {
		if estufa.atuadores[1].status == 0 {
			estufa.atuadores[1].status = 1
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[1])
		}
	} else {
		if estufa.atuadores[1].status == 1 {
			estufa.atuadores[1].status = 0
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[1])
		}
	}

	// VERIFICA SENSOR 2 _ ATUADOR 3 - IRRIGADOR
	if uint16(estufa.sensores[1].valor) < umidadeMin {
		if estufa.atuadores[2].status == 0 {
			estufa.atuadores[2].status = 1
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[2])
		}
	} else {
		if estufa.atuadores[2].status == 1 {
			estufa.atuadores[2].status = 0
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[2])
		}
	}

	// VERIFICA SENSOR 2 _ ATUADOR 3 - IRRIGADOR
	if uint16(estufa.sensores[2].valor) < nivelC02min {
		if estufa.atuadores[3].status == 0 {
			estufa.atuadores[3].status = 1
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[3])
		}
	} else {
		if estufa.atuadores[3].status == 1 {
			estufa.atuadores[3].status = 0
			atuadoresAlterados = append(atuadoresAlterados, estufa.atuadores[3])
		}
	}

	var buffer bytes.Buffer
	n := len(atuadoresAlterados)

	//COMO MAIS DE UM SENSOR PODE SER LIGADO OU DESLIGADO ENTAO n RECEBE A QUANTIDADE
	if n != 0 {
		buffer = convertAtuadoresEmArrayBytes(atuadoresAlterados, buffer, n)
	}

	//ENVIO PARA A ESTUFA OS SENSORES QUE FORAM ATIVADOS OU DESATIVADOS
	_, err := conn.Write(buffer.Bytes())
	checkError(err, "conn.Write/Sensores" )
}

func connRetornaSensorInfo(conn net.Conn) {
	result := make([]byte, 2)
	_, _ = conn.Read(result[:])
	valor := binary.BigEndian.Uint16(result)

	var dadosSensor Sensor
	// BUSCA QUAL SENSOR DEVE SER ENVIADO
	for _, sensor := range estufa.sensores {
		if sensor.id == valor {
			dadosSensor = sensor
		}
	}

	var buffer bytes.Buffer
	buffer = converteSensorEmArrayDeBytes(dadosSensor, buffer)

	// CODIFICA O PACOTE USANDO A SensorLayer DO PROTOCOLO
	pacote := gopacket.NewPacket(
		buffer.Bytes(),
		camada.SensorLayerType,
		gopacket.Default,
	)

	//ENVIO O PACOTE PRO CLIENTE
	_,_ = conn.Write(pacote.Data())
	_ = conn.Close()
}

func convertAtuadoresEmArrayBytes(atuadores []Atuador, buffer bytes.Buffer, n int) bytes.Buffer {
	var nBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(nBytes, uint16(n))
	buffer.Write(nBytes)

	for k := 0; k < n; k++ {
		var nomeBytes = make([]byte, 15)
		for i, j := range []byte(atuadores[k].nome) {
			nomeBytes[i] = byte(j)
		}

		var idBytes = make([]byte, 2)
		var statusBytes = make([]byte, 2)

		binary.BigEndian.PutUint16(idBytes, atuadores[k].id)
		binary.BigEndian.PutUint16(statusBytes, atuadores[k].status)

		buffer.Write(nomeBytes)
		buffer.Write(idBytes)
		buffer.Write(statusBytes)
	}

	return buffer
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
		_ :fmt.Fprintf(os.Stderr, "Erro em " + msg, err.Error())
		os.Exit(1)
	}
}