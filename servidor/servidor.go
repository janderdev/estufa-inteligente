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

type Estufa struct {
	sensores [] Sensor
	parametrosIni Parametros
}

func checkError(err error, msg string){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro em " + msg, err.Error())
		os.Exit(1)
	}
}

// DECLARAÇOES VARIAVEIS GLOBAIS DO SERVIDOR
var estufa Estufa
//-----------

func main() {
	fmt.Println("------- SERVIDOR INICIADO ---------")
	fmt.Println("Aguardando CLIENTE definir parms iniciais...")

	//DEFININDO VALORES PARA A STRUCT SENSORES
	var temperatura Sensor
	temperatura.nome = "Temperatura"
	temperatura.id = 1
	temperatura.valor = -36

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
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":3200")
	checkError(err, "ResolveTCPAddr")

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}


		connGetParametrosDoCliente(conn)
		break
	}

	//-----------------------------------------------
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
		camada.RequestLayerType,
		gopacket.Default,
	)

	conn.Write(pacote.Data())
	conn.Close()
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

func converteSensorEmArrayDeBytes(sensor struct {
	nome string
	id uint16
	valor int16
}, buffer bytes.Buffer) bytes.Buffer {

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