package pktfwd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/TheThingsNetwork/ttn/api/protocol/lorawan"
	"github.com/TheThingsNetwork/ttn/api/router"
)

type txPacketInfo struct {
	Imme bool    `json:"imme"`
	Tmst uint32  `json:"tmst"`
	Freq float64 `json:"freq"`
	Rfch uint32  `json:"rfch"`
	Powe int32   `json:"powe"`
	Modu string  `json:"modu"`
	Datr string  `json:"datr"`
	Codr string  `json:"codr"`
	Ipol bool    `json:"ipol"`
	Size int     `json:"size"`
	Data string  `json:"data"`
}

// JSON down: {"txpk":{"imme":false,"tmst":380885892,"freq":925.1,"rfch":0,"powe":20,"modu":"LORA","datr":"SF10BW500","codr":"4/5","ipol":true,"size":17,"data":"IPYHUhBOVUNx1q7l9CqUi8c="}}
type txPacket struct {
	Txpk txPacketInfo `json:"txpk"`
}

// LogTxPacketAsJSON prints a TX packet (downlink) in JSON format
func LogTxPacketAsJSON(downlink *router.DownlinkMessage) error {
	time := downlink.GetGatewayConfiguration().GetTimestamp()
	freq := downlink.GetGatewayConfiguration().GetFrequency()
	rfChain := downlink.GetGatewayConfiguration().GetRfChain()
	power := downlink.GetGatewayConfiguration().GetPower()
	modulation := downlink.GetProtocolConfiguration().GetLorawan().GetModulation()
	modulationStr := "UNKNOWN"
	switch modulation {
	case lorawan.Modulation_LORA:
		modulationStr = "LORA"
	case lorawan.Modulation_FSK:
		modulationStr = "FSK"
	}
	dr := downlink.GetProtocolConfiguration().GetLorawan().GetDataRate()
	cr := downlink.GetProtocolConfiguration().GetLorawan().GetCodingRate()
	ipol := downlink.GetGatewayConfiguration().GetPolarizationInversion()
	size := len(downlink.GetPayload())
	encodedPayload := base64.StdEncoding.EncodeToString(downlink.GetPayload())
	pktInfo := txPacketInfo{
		Imme: false,
		Tmst: time,
		Freq: float64(freq) / 1e6,
		Rfch: rfChain,
		Powe: power,
		Modu: modulationStr,
		Datr: dr,
		Codr: cr,
		Ipol: ipol,
		Size: size,
		Data: encodedPayload,
	}
	downMsg := txPacket{
		Txpk: pktInfo,
	}

	jsonMsgBytes, err := json.Marshal(downMsg)
	if err != nil {
		return err
	}
	jsonMsgString := string(jsonMsgBytes[:])
	fmt.Printf("JSON down: %v\n", jsonMsgString)

	return nil
}
