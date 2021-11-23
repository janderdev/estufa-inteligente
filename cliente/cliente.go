package main

import (
	"../camada"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"net"
	"os"
)

func checkError(err error, msg string){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro em " + msg, err.Error())
		os.Exit(1)
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":3200")
	checkError(err, "ResolveTCPAddr")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "DialTCP")

	parametros := inputParametrosIniciais()

	var buffer bytes.Buffer
	buffer = converteDadosEmArrayDeBytes(parametros, buffer)

	var pacote = gopacket.NewPacket(
		buffer.Bytes(),
		camada.ParametersLayerType,
		gopacket.Default,
	)

	conn.Write(pacote.Data())
	conn.Close()


	idSensor := reqIDSensorCliente()
	var idSensorBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(idSensorBytes, idSensor)
	conn.Write(idSensorBytes)

}

func reqIDSensorCliente() uint16 {
	var idSensor uint16
	fmt.Println("---- BUSCANDO INFORMAÇÕES DO SENSOR ----")
	fmt.Print("Digite o ID do Sensor: ")
	fmt.Scanln(&idSensor)
	return idSensor
}

func converteDadosEmArrayDeBytes(parametros struct {
	tempMin     int16
	tempMax     int16
	umidadeMin  uint16
	nivelCO2Min uint16
}, buffer bytes.Buffer) bytes.Buffer {

	var tempMinBytes = make([]byte, 4)
	var tempMaxBytes = make([]byte, 4)
	var umidadeMinBytes = make([]byte, 2)
	var nivelCO2MinBytes = make([]byte, 2)

	binary.BigEndian.PutUint32(tempMinBytes, uint32(parametros.tempMin))
	binary.BigEndian.PutUint32(tempMaxBytes, uint32(parametros.tempMax))
	binary.BigEndian.PutUint16(umidadeMinBytes, parametros.umidadeMin)
	binary.BigEndian.PutUint16(nivelCO2MinBytes, parametros.nivelCO2Min)

	buffer.Write(tempMinBytes)
	buffer.Write(tempMaxBytes)
	buffer.Write(umidadeMinBytes)
	buffer.Write(nivelCO2MinBytes)

	return buffer
}

func inputParametrosIniciais() struct {
	tempMin     int16
	tempMax     int16
	umidadeMin  uint16
	nivelCO2Min uint16
} {
	t := struct {
		tempMin     int16
		tempMax     int16
		umidadeMin  uint16
		nivelCO2Min uint16
	}{}

	fmt.Println("---- DEFINA OS PARAMETROS INICIAIS DA ESTUFA ----")
	fmt.Print("Temperatura Mínima: ")
	fmt.Scanln(&t.tempMin)
	fmt.Print("Temperatura Máxima: ")
	fmt.Scanln(&t.tempMax)
	fmt.Print("Umidade Mínima: ")
	fmt.Scanln(&t.umidadeMin)
	fmt.Print("Nível de CO2 Mínimo: ")
	fmt.Scanln(&t.nivelCO2Min)
	fmt.Println("-------------------------------------------------")

	return t
}



