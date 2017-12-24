// The MIT License (MIT)
//
// Copyright © 2017 Sven Agneessens <sven.agneessens@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package protocol

import (
	"encoding/binary"
	"encoding/json"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// PullRespPacket is used by the server to send RF packets and associated
// metadata that will have to be emitted by the gateway.
type PullRespPacket struct {
	Protocol    uint8
	RandomToken uint16
	Payload     PullRespPayload
}

// PullRespPayload represents the downstream JSON data structure.
type PullRespPayload struct {
	TXPK TXPK `json:"txpk"`
}

// TXPK contains a RF packet to be emitted and associated metadata.
type TXPK struct {
	Imme bool     `json:"imme"`           // imme | bool   | Send packet immediately (will ignore tmst & time)
	Tmst uint32   `json:"tmst,omitempty"` // tmst | number | Send packet on a certain timestamp value (will ignore time)
	Tmms int64    `json:"tmms,omitempty"` // tmms | number | Send packet at a certain GPS time (GPS synchronization required)
	Freq float64  `json:"freq"`           // freq | number | TX central frequency in MHz (unsigned float, Hz precision)
	RFCh uint8    `json:"rfch"`           // rfch | number | Concentrator "RF chain" used for TX (unsigned integer)
	Powe uint8    `json:"powe"`           // powe | number | TX output power in dBm (unsigned integer, dBm precision)
	Modu string   `json:"modu"`           // modu | string | Modulation identifier "LORA" or "FSK"
	DatR DataRate `json:"datr"`           // datr | string | LoRa datarate identifier (eg. SF12BW500) || datr | number | FSK datarate (unsigned, in bits per second)
	CodR string   `json:"codr,omitempty"` // codr | string | LoRa ECC coding rate identifier
	FDev uint16   `json:"fdev,omitempty"` // fdev | number | FSK frequency deviation (unsigned integer, in Hz)
	IPol bool     `json:"ipol"`           // ipol | bool   | Lora modulation polarization inversion
	Prea uint16   `json:"prea,omitempty"` // prea | number | RF preamble size (unsigned integer)
	Size uint16   `json:"size"`           // size | number | RF packet payload size in bytes (unsigned integer)
	Data string   `json:"data"`           // data | string | Base64 encoded RF packet payload, padding optional
	NCRC bool     `json:"ncrc,omitempty"` // ncrc | bool   | If true, disable the CRC of the physical layer (optional)
}

func handlePullResp(data []byte) (Packet, error) {
	var pullRespPacket PullRespPacket

	err := pullRespPacket.unmarshalData(data)
	if err != nil {
		return nil, errors.Wrap(err, "handle pull resp packet failed")
	}

	return &pullRespPacket, nil
}

func (p *PullRespPacket) Log(ctx log.Interface) {
	ctx.WithFields(log.Fields{
		"protocol":               p.Protocol,
		"random token":           p.RandomToken,
		"immediately":            p.Payload.TXPK.Imme,
		"timestamp":              p.Payload.TXPK.Tmst,
		"gps time":               p.Payload.TXPK.Tmms,
		"frequency":              p.Payload.TXPK.Freq,
		"RF chain":               p.Payload.TXPK.RFCh,
		"power":                  p.Payload.TXPK.Powe,
		"modulation":             p.Payload.TXPK.Modu,
		"data rate":              p.Payload.TXPK.DatR,
		"coding rate":            p.Payload.TXPK.CodR,
		"polarization inversion": p.Payload.TXPK.IPol,
		"preamble size":          p.Payload.TXPK.Prea,
		"no crc":                 p.Payload.TXPK.NCRC,
		"size":                   p.Payload.TXPK.Size,
		"data":                   p.Payload.TXPK.Data,
	}).Info("PULL_RESP")
}

func (p *PullRespPacket) unmarshalData(data []byte) error {
	_, err := isValidPullRespPacket(data)
	if err != nil {
		return errors.Wrap(err, "unmarshal pull resp packet failed")
	}

	p.Protocol = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])

	return json.Unmarshal(data[4:], &p.Payload)
}

func isValidPullRespPacket(data []byte) (bool, error) {
	if len(data) < 4 {
		return false, errors.New("invalid packet: at least 4 bytes expected")
	}

	return true, nil
}
