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

	var tempMin = -10
	var tempMax = 5
	var umidadeMin = 100
	var nivelCO2Min	= 400

	var buffer bytes.Buffer

	var tempMinBytes = make([]byte, 4)
	var tempMaxBytes = make([]byte, 4)
	var umidadeMinBytes = make([]byte, 2)
	var nivelCO2MinBytes = make([]byte, 2)

	binary.BigEndian.PutUint32(tempMinBytes, uint32(tempMin))
	binary.BigEndian.PutUint32(tempMaxBytes, uint32(tempMax))
	binary.BigEndian.PutUint16(umidadeMinBytes, uint16(umidadeMin))
	binary.BigEndian.PutUint16(nivelCO2MinBytes, uint16(nivelCO2Min))

	buffer.Write(tempMinBytes)
	buffer.Write(tempMaxBytes)
	buffer.Write(umidadeMinBytes)
	buffer.Write(nivelCO2MinBytes)

	var pacote = gopacket.NewPacket(
		buffer.Bytes(),
		camada.ParametersLayerType,
		gopacket.Default,

	)

	conn.Write(pacote.Data())
	conn.Close()

	os.Exit(0)
}


