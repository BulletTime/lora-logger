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

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// TXACKPacket is used by the gateway to send a feedback to the server
// to inform if a downlink request has been accepted or rejected by the
// gateway.
type TXAckPacket struct {
	Protocol    uint8
	RandomToken uint16
	GatewayMac  [8]byte
	Payload     TXAckPayload
}

// TXACKPayload contains the TXACKPacket payload.
type TXAckPayload struct {
	TXPKACK TXPKACK `json:"txpk_ack"`
}

// TXPKACK contains the status information of the associated PULL_RESP
// packet.
type TXPKACK struct {
	Error string `json:"error"`
}

func handleTXAck(data []byte) (Packet, error) {
	var packet PullAckPacket

	err := packet.unmarshalData(data)
	if err != nil {
		return nil, errors.Wrap(err, "handle tx ack packet failed")
	}

	return &packet, nil
}

func (p *TXAckPacket) Log(ctx log.Interface) {
	ctx.WithFields(log.Fields{
		"protocol":     p.Protocol,
		"random token": p.RandomToken,
		"gateway mac":  fmt.Sprintf("%X", p.GatewayMac),
		"error":        p.Payload.TXPKACK.Error,
	}).Info("TX_ACK")
}

func (p *TXAckPacket) unmarshalData(data []byte) error {
	_, err := isValidTXAckPacket(data)
	if err != nil {
		return errors.Wrap(err, "unmarshal tx ack packet failed")
	}

	p.Protocol = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])

	for i := 0; i < 8; i++ {
		p.GatewayMac[i] = data[4+i]
	}

	return json.Unmarshal(data[12:], p.Payload)
}

func isValidTXAckPacket(data []byte) (bool, error) {
	if len(data) < 12 {
		return false, errors.New("invalid packet: at least 12 bytes expected")
	}

	return true, nil
}
