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

type Parametros struct{
	tempMin     int16
	tempMax     int16
	umidadeMin  uint16
	nivelCO2Min uint16
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":3200")
	checkError(err, "ResolveTCPAddr")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "DialTCP")

	// PEGA OS DADOS PASSADOS PELO TERMINAL DO USUARIO E RETORNA PARA parametros
	parametros := inputParametrosIniciais()

	var buffer bytes.Buffer
	buffer = converteDadosEmArrayDeBytes(parametros, buffer)

	// CRIA O PACOTE DE ACORDO COM ParametersLayer DO PROTOCOLO
	var pacote = gopacket.NewPacket(
		buffer.Bytes(),
		camada.ParametersLayerType,
		gopacket.Default,
	)

	// ENVIA PARA O SERVIDOR
	_,_ = conn.Write(pacote.Data())
	_ = conn.Close()

	// MANTEM A CONEXAO EM UM LOOP INFINITO
	for {

		idSensor := reqIDSensorUsuario()
		var idSensorBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(idSensorBytes, idSensor)

		// QUANDO O USUÁRIO SOLICITA UM SENSOR PELO TERMINAL
		// UMA NOVA CONEXAO É ESTABELECIDA COM O SERVIDOR E OS DADOS SÃO RECEBIDOS
		novaConnGetSensorInfo(idSensorBytes, tcpAddr)
	}
}

func novaConnGetSensorInfo(codSensorBytes []byte, addr *net.TCPAddr) {
	novaConn, err := net.DialTCP("tcp", nil, addr)
	checkError(err, "DialTCP")
	novaConn.Write(codSensorBytes)

	result, err := ioutil.ReadAll(novaConn)
	checkError(err, "ReadAll")

	packet := gopacket.NewPacket(
		result,
		camada.SensorLayerType,
		gopacket.Default,
	)

	//PACOTE E DECODIFICADO DE ACORDO COM A SensorLayer EXIGE
	decodePacket := packet.Layer(camada.SensorLayerType)

	if decodePacket != nil {
		fmt.Println("####### DADOS DO SENSOR ########")
		content, _ := decodePacket.(*camada.SensorLayer)
		fmt.Println("ID:", content.IDSensor)
		fmt.Println("Nome:", content.Nome)
		fmt.Println("Valor:", int16(content.Valor))
		fmt.Println("--------------------------------")
	}

	_ = novaConn.Close()
}

func reqIDSensorUsuario() uint16 {
	var idSensor uint16
	fmt.Println("--------- SOLICITAR LEITURA DO SENSOR ---------")
	fmt.Print("Digite o ID do Sensor (Ex: 1, 2, 3): ")
	fmt.Scanln(&idSensor)
	return idSensor
}

func converteDadosEmArrayDeBytes(parametros Parametros, buffer bytes.Buffer) bytes.Buffer {

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

func inputParametrosIniciais() Parametros {
	var t Parametros

	fmt.Println("---- DEFINA OS PARAMETROS INICIAIS DA ESTUFA ----")
	fmt.Print("Temperatura Mínima: ")
	fmt.Scanln(&t.tempMin)
	fmt.Print("Temperatura Máxima: ")
	fmt.Scanln(&t.tempMax)
	fmt.Print("Umidade Mínima (>=0): ")
	fmt.Scanln(&t.umidadeMin)
	fmt.Print("Nível de CO2 Mínimo (>=0): ")
	fmt.Scanln(&t.nivelCO2Min)
	fmt.Println("-------------------------------------------------")

	return t
}

func checkError(err error, msg string){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro em " + msg, err.Error())
		os.Exit(1)
	}
}
