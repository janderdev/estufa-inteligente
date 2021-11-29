package camada

import (
	"encoding/binary"
	"github.com/google/gopacket"
)

// ParametersLayer É uma camada do protocolo que trata a trasmissão dos parametros min e max dos sensores da estufa passados pelo cliente.
type ParametersLayer struct {
	TemperaturaMin  uint32
	TemperaturaMax  uint32
	UmidadeMin 	    uint16
	NivelCO2Min     uint16
	restoDosDados   [] byte
}

// SensorLayer É uma camada do protocolo que trata comunicação (codificação e decodificação) dos dados dos sensores na aplicação.
type SensorLayer struct {
	Nome string
	IDSensor uint16
	Valor uint32
	restoDosDados [] byte
}

// AtuadorLayer É uma camada do protocolo que trata comunicação (codificação e decodificação) dos dados dos sensores na aplicação.
type AtuadorLayer struct {
	Nome string
	IDAtuador uint16
	Status uint8
	restoDosDados [] byte
}

type SensoresLayer struct {
	Temperatura SensorLayer
	Umidade SensorLayer
	NivelDeCO2 SensorLayer
	restoDosDados [] byte
}


var ParametersLayerType = gopacket.RegisterLayerType(
	2002,
	gopacket.LayerTypeMetadata{
		"ParametersLayerType",
		gopacket.DecodeFunc(decodeParametersLayer),
	},
)

var SensorLayerType = gopacket.RegisterLayerType(
	2003,
	gopacket.LayerTypeMetadata{
		"SensorLayerType",
		gopacket.DecodeFunc(decodeSensorLayer),
	},
)

var SensoresLayerType = gopacket.RegisterLayerType(
	2004,
	gopacket.LayerTypeMetadata{
		"SensoresLayerType",
		gopacket.DecodeFunc(decodeSensoresLayer),
})

var AtuadorLayerType = gopacket.RegisterLayerType(
	2005,
	gopacket.LayerTypeMetadata{
		"AtuadorLayerType",
		gopacket.DecodeFunc(decodeAtuadorLayer),
})



/* ParametersLayer FUNCTIONS -------------------------------------- */
func decodeParametersLayer(data []byte, p gopacket.PacketBuilder) error {
	temperaturaMin := binary.BigEndian.Uint32(data[0:4])
	temperaturaMax := binary.BigEndian.Uint32(data[4:8])
	umidadeMin := binary.BigEndian.Uint16(data[8:10])
	nivelCO2Min := binary.BigEndian.Uint16(data[10:12])
	var restoDosDados []byte = nil

	if len(data) >= 12 {
		restoDosDados = data[12:]
	}

	p.AddLayer(&ParametersLayer {
		temperaturaMin,
		temperaturaMax,
		umidadeMin,
		nivelCO2Min,
		restoDosDados,
	})

	return p.NextDecoder(gopacket.LayerTypePayload)
}

func (p ParametersLayer) LayerType() gopacket.LayerType {
	return ParametersLayerType
}

func (p ParametersLayer) LayerPayload() []byte {
	return p.restoDosDados
}

func (p ParametersLayer) LayerContents() []byte {
	var tempMinBytes = make([]byte, 3)
	binary.BigEndian.PutUint32(tempMinBytes, p.TemperaturaMin)

	var tempMaxBytes = make([]byte, 3)
	binary.BigEndian.PutUint32(tempMaxBytes, p.TemperaturaMax)

	var umidadeMinBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(umidadeMinBytes, p.UmidadeMin)

	var nivelCO2MinBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(nivelCO2MinBytes, p.NivelCO2Min)

	var contents []byte
	contents = append(tempMinBytes)
	contents = append(tempMaxBytes)
	contents = append(umidadeMinBytes)
	contents = append(nivelCO2MinBytes)

	return contents
}


/* SensorLayer FUNCTIONS ----------------------------------------- */
func decodeSensorLayer(data []byte, p gopacket.PacketBuilder) error {
	nome := string(data[:15])
	idSensor := binary.BigEndian.Uint16(data[15:17])
	valor := binary.BigEndian.Uint32(data[17:21])
	var restoDosDados []byte = nil

	if len(data) >= 21 {
		restoDosDados = data[21:]
	}

	p.AddLayer(&SensorLayer{
		nome,
		idSensor,
		valor,
		restoDosDados,
	})

	return p.NextDecoder(gopacket.LayerTypePayload)
}

func (e SensorLayer) LayerType() gopacket.LayerType {
	return SensorLayerType
}

func (e SensorLayer) LayerPayload() []byte {
	return e.restoDosDados
}

func (e SensorLayer) LayerContents() []byte {
	var nomeBytes = []byte(e.Nome)
	var idSensorBytes = make([]byte, 2)
	var valorBytes = make([]byte, 4)

	binary.BigEndian.PutUint16(idSensorBytes, e.IDSensor)
	binary.BigEndian.PutUint32(valorBytes, e.Valor)

	var contents []byte
	contents = append(nomeBytes)
	contents = append(idSensorBytes)
	contents = append(valorBytes)

	return contents
}

/* SensoresLayer FUNCTIONS ----------------------------------------- */
func decodeSensoresLayer(data []byte, p gopacket.PacketBuilder) error {
	var temperatura SensorLayer
	var umidade SensorLayer
	var nivelCO2 SensorLayer

	temperatura.Nome = string(data[:15])
	temperatura.IDSensor = binary.BigEndian.Uint16(data[15:17])
	temperatura.Valor = binary.BigEndian.Uint32(data[17:21])

	umidade.Nome = string(data[21:36])
	umidade.IDSensor = binary.BigEndian.Uint16(data[36:38])
	umidade.Valor = binary.BigEndian.Uint32(data[38:42])

	nivelCO2.Nome = string(data[42:57])
	nivelCO2.IDSensor = binary.BigEndian.Uint16(data[57:59])
	nivelCO2.Valor = binary.BigEndian.Uint32(data[59:63])

	var restoDosDados []byte = nil

	if len(data) >= 63 {
		restoDosDados = data[63:]
	}


	p.AddLayer(&SensoresLayer{
		temperatura,
		umidade,
		nivelCO2,
		restoDosDados,
	})

	return p.NextDecoder(gopacket.LayerTypePayload)
}

func (s SensoresLayer) LayerType() gopacket.LayerType {
	return SensoresLayerType
}

func (s SensoresLayer) LayerPayload() []byte {
	return s.restoDosDados
}

func (s SensoresLayer) LayerContents() []byte {
	//CONVERTERNDO DADOS PARA BYTES DO SENSOR UMIDADE
	var nomeBytes = []byte(s.Temperatura.Nome)
	var idSensorBytes = make([]byte, 2)
	var valorBytes = make([]byte, 4)

	binary.BigEndian.PutUint16(idSensorBytes, s.Temperatura.IDSensor)
	binary.BigEndian.PutUint32(valorBytes, s.Temperatura.Valor)

	//PASSANDO PRO ARRAY DE BYTES
	var contents []byte
	contents = append(nomeBytes)
	contents = append(idSensorBytes)
	contents = append(valorBytes)

	//CONVERTERNDO DADOS PARA BYTES DO SENSOR UMIDADE
	nomeBytes = []byte(s.Umidade.Nome)
	idSensorBytes = make([]byte, 2)
	valorBytes = make([]byte, 4)

	binary.BigEndian.PutUint16(idSensorBytes, s.Umidade.IDSensor)
	binary.BigEndian.PutUint32(valorBytes, s.Umidade.Valor)

	//PASSANDO PRO ARRAY DE BYTES
	contents = append(nomeBytes)
	contents = append(idSensorBytes)
	contents = append(valorBytes)

	//CONVERTERNDO DADOS PARA BYTES DO SENSOR NIVELDECO2
	nomeBytes = []byte(s.NivelDeCO2.Nome)
	idSensorBytes = make([]byte, 2)
	valorBytes = make([]byte, 4)

	binary.BigEndian.PutUint16(idSensorBytes, s.NivelDeCO2.IDSensor)
	binary.BigEndian.PutUint32(valorBytes, s.NivelDeCO2.Valor)

	//PASSANDO PRO ARRAY DE BYTES
	contents = append(nomeBytes)
	contents = append(idSensorBytes)
	contents = append(valorBytes)

	return contents
}


/* AtuadorLayer FUNCTIONS -------------------------------------- */
func decodeAtuadorLayer(data []byte, p gopacket.PacketBuilder) error {
	nome := string(data[:15])
	idAtuador := binary.BigEndian.Uint16(data[15:17])
	status := data[17]

	var restoDosDados []byte = nil
	if len(data) >= 17 {
		restoDosDados = data[17:]
	}

	p.AddLayer(&AtuadorLayer{
		nome,
		idAtuador,
		status,
		restoDosDados,
	})

	return p.NextDecoder(gopacket.LayerTypePayload)
}

func (a AtuadorLayer) LayerType() gopacket.LayerType {
	return AtuadorLayerType
}

func (a AtuadorLayer) LayerPayload() []byte {
	return a.restoDosDados
}

func (a AtuadorLayer) LayerContents() []byte {
	var nomeBytes = []byte(a.Nome)
	var statusByte = []byte{a.Status}
	var idAtuadorBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(idAtuadorBytes, a.IDAtuador)

	var contents []byte
	contents = append(nomeBytes)
	contents = append(idAtuadorBytes)
	contents = append(statusByte)

	return contents
}


