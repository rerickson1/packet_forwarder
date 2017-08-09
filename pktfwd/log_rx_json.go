package pktfwd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/TheThingsNetwork/packet_forwarder/wrapper"
)

type rxPacket struct {
	Tmst uint32  `json:"tmst"`
	Ch   int     `json:"chan"`
	Rfch int     `json:"rfch"`
	Freq float64 `json:"freq"`
	Stat int     `json:"stat"`
	Modu string  `json:"modu"`
	Datr string  `json:"datr"`
	Codr string  `json:"codr"`
	Lsnr float32 `json:"lsnr"`
	Rssi float32 `json:"rssi"`
	Size uint32  `json:"size"`
	Data string  `json:"data"`
}

type rxPackets struct {
	Rxpk []rxPacket `json:"rxpk"`
}

// LogRxPacketAsJSON prints an RX packet in JSON format
func LogRxPacketAsJSON(packet wrapper.Packet) error {
	upMsg := new(rxPackets)
	encodedPayload := base64.StdEncoding.EncodeToString(packet.Payload)
	bw, _ := packet.BandwidthString()
	dr, _ := packet.DatarateString()
	cr, _ := packet.CoderateString()
	var status int
	switch packet.Status {
	case wrapper.StatusCRCOK:
		status = 1
	case wrapper.StatusCRCBAD:
		status = -1
	case wrapper.StatusNOCRC:
		status = 0
	}

	modulation := "UNKNOWN"
	switch packet.Modulation {
	case wrapper.ModulationLoRa:
		modulation = "LORA"
	case wrapper.ModulationFSK:
		modulation = "FSK"
	}

	jsonPkt := rxPacket{
		Tmst: packet.CountUS,
		Ch:   int(packet.IFChain),
		Rfch: int(packet.RFChain),
		Freq: float64(packet.Freq) / 1e6,
		Stat: status,
		Modu: modulation,
		Datr: dr + bw,
		Codr: cr,
		Lsnr: packet.SNR,
		Rssi: packet.RSSI,
		Size: packet.Size,
		Data: encodedPayload,
	}
	upMsg.Rxpk = append(upMsg.Rxpk, jsonPkt)

	jsonMsgBytes, err := json.Marshal(*upMsg)
	if err != nil {
		return err
	}
	jsonMsgString := string(jsonMsgBytes[:])
	fmt.Printf("JSON up: %v\n", jsonMsgString)

	return nil
}
