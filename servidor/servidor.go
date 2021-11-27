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
	valor uint32
}

type Estufa struct {
	sensores [] Sensor
}

func checkError(err error, msg string){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro em " + msg, err.Error())
		os.Exit(1)
	}
}

func main() {
	//STRUCT SENSORES
	var estufa Estufa
	var temperatura Sensor
	temperatura.nome = "0123456789"
	temperatura.id = 1
	temperatura.valor = 36

	//PESQUISAR SOBRE APPEND
	estufa.sensores = append(estufa.sensores, temperatura)

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

			connRetornaSensorInfo(novaConn, estufa.sensores)
		}
	}

}

func connRetornaSensorInfo(conn net.Conn, sensores []Sensor, ) {
	result := make([]byte, 2)
	conn.Read(result[:])
	valor := binary.BigEndian.Uint16(result)

	var dadosSensor Sensor
	for _, sensor := range sensores {
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
		fmt.Println("--- PARAMETROS DA ESTUFA PASSADOS PELO CLIENTE ---")
		content, _ := decodePacket.(*camada.ParametersLayer)
		fmt.Println("TemperaturaMin:", int32(content.TemperaturaMin))
		fmt.Println("TemperaturaMax:", int32(content.TemperaturaMax))
		fmt.Println("UmidadeMin:", content.UmidadeMin)
		fmt.Println("NivelCO2Min:", content.NivelCO2Min)
		fmt.Println("---------------------------------------------------")
	}
	conn.Close()
}

func converteSensorEmArrayDeBytes(sensor struct {
	nome string
	id uint16
	valor uint32
}, buffer bytes.Buffer) bytes.Buffer {

	//var nomeBytes = make([]byte, 15)
	var nomeBytes = make([]byte, 15)
	for i, j := range []byte(sensor.nome) {
		nomeBytes[i] = byte(j)
	}

	var idBytes = make([]byte, 2)
	var valorBytes = make([]byte, 4)

	binary.BigEndian.PutUint16(idBytes, sensor.id)
	binary.BigEndian.PutUint32(valorBytes, uint32(sensor.valor))

	fmt.Println(nomeBytes)
	buffer.Write(nomeBytes)
	buffer.Write(idBytes)
	buffer.Write(valorBytes)

	return buffer
}