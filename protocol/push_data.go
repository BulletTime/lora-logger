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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// PushDataPacket type is used by the gateway mainly to forward the RF packets
// received, and associated metadata, to the server.
type PushDataPacket struct {
	Protocol    uint8
	RandomToken uint16
	GatewayMac  [8]byte
	Payload     PushDataPayload
}

// PushDataPayload represents the upstream JSON data structure.
type PushDataPayload struct {
	RXPK []RXPK `json:"rxpk,omitempty"`
	Stat *Stat  `json:"stat,omitempty"`
}

// RXPK contains an RF packet and associated metadata.
type RXPK struct {
	Time CompactTime `json:"time"` // time | string | UTC time of pkt RX, us precision, ISO 8601 'compact' format
	TMMS int64       `json:"tmms"` // tmms | number | GPS time of pkt RX, number of milliseconds since 06.Jan.1980
	TMST uint32      `json:"tmst"` // tmst | number | Internal timestamp of "RX finished" event (32b unsigned)
	Freq float64     `json:"freq"` // freq | number | RX central frequency in MHz (unsigned float, Hz precision)
	Chan uint8       `json:"chan"` // chan | number | Concentrator "IF" channel used for RX (unsigned integer)
	RFCh uint8       `json:"rfch"` // rfch | number | Concentrator "RF chain" used for RX (unsigned integer)
	Stat int8        `json:"stat"` // stat | number | CRC status: 1 = OK, -1 = fail, 0 = no CRC
	Mod  string      `json:"modu"` // modu | string | Modulation identifier "LORA" or "FSK"
	DatR *DataRate   `json:"datr"` // datr | string | LoRa datarate identifier (eg. SF12BW500) || datr | number | FSK datarate (unsigned, in bits per second)
	CodR string      `json:"codr"` // codr | string | LoRa ECC coding rate identifier
	RSSI int16       `json:"rssi"` // rssi | number | RSSI in dBm (signed integer, 1 dB precision)
	SNR  float64     `json:"lsnr"` // lsnr | number | Lora SNR ratio in dB (signed float, 0.1 dB precision)
	Size uint16      `json:"size"` // size | number | RF packet payload size in bytes (unsigned integer)
	Data string      `json:"data"` // data | string | Base64 encoded RF packet payload, padded
}

// Stat contains the status of the gateway.
type Stat struct {
	Time ExpandedTime `json:"time"` // time | string | UTC 'system' time of the gateway, ISO 8601 'expanded' format
	Lati float64      `json:"lati"` // lati | number | GPS latitude of the gateway in degree (float, N is +)
	Long float64      `json:"long"` // long | number | GPS latitude of the gateway in degree (float, E is +)
	Alti int32        `json:"alti"` // alti | number | GPS altitude of the gateway in meter RX (integer)
	RXNb uint32       `json:"rxnb"` // rxnb | number | Number of radio packets received (unsigned integer)
	RXOK uint32       `json:"rxok"` // rxok | number | Number of radio packets received with a valid PHY CRC
	RXFW uint32       `json:"rxfw"` // rxfw | number | Number of radio packets forwarded (unsigned integer)
	ACKR float64      `json:"ackr"` // ackr | number | Percentage of upstream datagrams that were acknowledged
	DWNb uint32       `json:"dwnb"` // dwnb | number | Number of downlink datagrams received (unsigned integer)
	TXNb uint32       `json:"txnb"` // txnb | number | Number of packets emitted (unsigned integer)
}

// ExpandedTime implements time.Time but (un)marshals to and from
// ISO 8601 'expanded' format.
type ExpandedTime time.Time

// MarshalJSON implements the json.Marshaler interface for ExpandedTime.
func (t ExpandedTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).UTC().Format(`"2006-01-02 15:04:05 MST"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for ExpandedTime.
func (t *ExpandedTime) UnmarshalJSON(data []byte) error {
	t2, err := time.Parse(`"2006-01-02 15:04:05 MST"`, string(data))
	if err != nil {
		return err
	}
	*t = ExpandedTime(t2)
	return nil
}

// CompactTime implements time.Time but (un)marshals to and from
// ISO 8601 'compact' format.
type CompactTime time.Time

// MarshalJSON implements the json.Marshaler interface for CompactTime.
func (t CompactTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).UTC().Format(`"` + time.RFC3339Nano + `"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for CompactTime.
func (t *CompactTime) UnmarshalJSON(data []byte) error {
	t2, err := time.Parse(`"`+time.RFC3339Nano+`"`, string(data))
	if err != nil {
		return err
	}
	*t = CompactTime(t2)
	return nil
}

// DataRate implements the data rate which can be either a string (LoRa identifier)
// or an unsigned integer in case of FSK (bits per second).
type DataRate struct {
	LoRa string
	FSK  uint32
}

// String implements the stringer interface for DataRate.
func (d DataRate) String() string {
	if d.LoRa != "" {
		return d.LoRa
	}

	return string(d.FSK)
}

// MarshalJSON implements the json.Marshaler interface for DataRate.
func (d DataRate) MarshalJSON() ([]byte, error) {
	if d.LoRa != "" {
		return []byte(`"` + d.LoRa + `"`), nil
	}
	return []byte(strconv.FormatUint(uint64(d.FSK), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for DataRate.
func (d *DataRate) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		d.LoRa = strings.Trim(string(data), `"`)
		return nil
	}
	d.FSK = uint32(i)
	return nil
}

func handlePushData(data []byte) (Packet, error) {
	var pushDataPacket PushDataPacket

	err := pushDataPacket.unmarshalData(data)
	if err != nil {
		return nil, errors.Wrap(err, "handle push data packet failed")
	}

	return &pushDataPacket, nil
}

func (p *PushDataPacket) Log(ctx log.Interface) {
	ctx = ctx.WithFields(log.Fields{
		"protocol":     p.Protocol,
		"random token": p.RandomToken,
		"gateway mac":  fmt.Sprintf("%X", p.GatewayMac),
	})

	for _, rxpk := range p.Payload.RXPK {
		ctx.WithFields(log.Fields{
			"time":        time.Time(rxpk.Time),
			"frequency":   rxpk.Freq,
			"IF channel":  rxpk.Chan,
			"RF chain":    rxpk.RFCh,
			"crc":         rxpk.Stat,
			"modulation":  rxpk.Mod,
			"data rate":   rxpk.DatR,
			"coding rate": rxpk.CodR,
			"rssi":        rxpk.RSSI,
			"snr":         rxpk.SNR,
			"size":        rxpk.Size,
			"data":        rxpk.Data,
		}).Info("PUSH_DATA: RXPK")
	}

	if p.Payload.Stat != nil {
		ctx.WithFields(log.Fields{
			"time":                time.Time(p.Payload.Stat.Time),
			"rx received":         p.Payload.Stat.RXNb,
			"rx ok":               p.Payload.Stat.RXOK,
			"rx forwarded":        p.Payload.Stat.RXFW,
			"upstream ack (%)":    p.Payload.Stat.ACKR,
			"downstream received": p.Payload.Stat.DWNb,
			"tx ps":               p.Payload.Stat.TXNb,
		}).Info("PUSH_DATA: STAT")
	}
}

func (p *PushDataPacket) unmarshalData(data []byte) error {
	_, err := isValidPushDataPacket(data)
	if err != nil {
		return errors.Wrap(err, "unmarshal push data packet failed")
	}

	p.Protocol = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])

	for i := 0; i < 8; i++ {
		p.GatewayMac[i] = data[4+i]
	}

	return json.Unmarshal(data[12:], &p.Payload)
}

func isValidPushDataPacket(data []byte) (bool, error) {
	if len(data) < 12 {
		return false, errors.New("invalid packet: at least 12 bytes expected")
	}

	return true, nil
}
