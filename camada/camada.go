package camada //camada

import (
	"encoding/binary"
	"github.com/google/gopacket"
)

// ParametersLayer É uma camada do protocolo que trata a trasmissão dos parametros de entrada da estufa passados pelo cliente.
type ParametersLayer struct {
	TemperaturaMin  uint32
	TemperaturaMax  uint32
	UmidadeMin 	    uint16
	NivelCO2Min     uint16
	restoDosDados   [] byte
}

// RequestLayer É uma camada do protocolo que trata a requisição, feita pelo cliente, da leitura de um sensor específico na estufa.
type RequestLayer struct {
	Nome string
	IDSensor uint16
	Valor uint32
	restoDosDados [] byte
}

var ParametersLayerType = gopacket.RegisterLayerType(
	2002,
	gopacket.LayerTypeMetadata{
		"ParametersLayerType",
		gopacket.DecodeFunc(decodeParametersLayer),
	},
)

var RequestLayerType = gopacket.RegisterLayerType(
	2003,
	gopacket.LayerTypeMetadata{
		"RequestLayerType",
		gopacket.DecodeFunc(decodeRequestLayerType),
	},
)

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

func (e ParametersLayer) LayerType() gopacket.LayerType {
	return ParametersLayerType
}

func (e ParametersLayer) LayerPayload() []byte {
	return e.restoDosDados
}

func (e ParametersLayer) LayerContents() []byte {
	var tempMinBytes = make([]byte, 3)
	binary.BigEndian.PutUint32(tempMinBytes, e.TemperaturaMin)

	var tempMaxBytes = make([]byte, 3)
	binary.BigEndian.PutUint32(tempMaxBytes, e.TemperaturaMax)

	var umidadeMinBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(umidadeMinBytes, e.UmidadeMin)

	var nivelCO2MinBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(nivelCO2MinBytes, e.NivelCO2Min)

	var contents []byte
	contents = append(tempMinBytes)
	contents = append(tempMaxBytes)
	contents = append(umidadeMinBytes)
	contents = append(nivelCO2MinBytes)

	return contents
}
/* ---------------------------------------------------------------- */

/* RequestLayer FUNCTIONS ----------------------------------------- */
func decodeRequestLayerType(data []byte, p gopacket.PacketBuilder) error {
	nome := string(data[:15])
	idSensor := binary.BigEndian.Uint16(data[15:17])
	valor := binary.BigEndian.Uint32(data[17:20])
	var restoDosDados []byte = nil

	if len(data) >= 20 {
		restoDosDados = data[20:]
	}

	p.AddLayer(&RequestLayer {
		nome,
		idSensor,
		valor,
		restoDosDados,
	})

	return p.NextDecoder(gopacket.LayerTypePayload)
}

func (e RequestLayer) LayerType() gopacket.LayerType {
	return RequestLayerType
}

func (e RequestLayer) LayerPayload() []byte {
	return e.restoDosDados
}

func (e RequestLayer) LayerContents() []byte {
	var nomeBytes = []byte(e.Nome)
	var idSensorBytes = make([]byte, 2)
	var valorBytes = make([]byte, 3)

	binary.BigEndian.PutUint16(idSensorBytes, e.IDSensor)
	binary.BigEndian.PutUint32(valorBytes, e.Valor)

	var contents []byte
	contents = append(nomeBytes)
	contents = append(idSensorBytes)
	contents = append(valorBytes)

	return contents
}
/* ---------------------------------------------------------------- */




