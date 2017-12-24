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
	"fmt"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// PullDataPacket is used by the gateway to poll data from the server.
type PullDataPacket struct {
	Protocol    uint8
	RandomToken uint16
	GatewayMac  [8]byte
}

func handlePullData(data []byte) (Packet, error) {
	var packet PullDataPacket

	err := packet.unmarshalData(data)
	if err != nil {
		return nil, errors.Wrap(err, "handle pull data packet failed")
	}

	return &packet, nil
}

func (p *PullDataPacket) Log(ctx log.Interface) {
	ctx.WithFields(log.Fields{
		"protocol":     p.Protocol,
		"random token": p.RandomToken,
		"gateway mac":  fmt.Sprintf("%X", p.GatewayMac),
	}).Info("PULL_DATA")
}

func (p *PullDataPacket) unmarshalData(data []byte) error {
	_, err := isValidPullDataPacket(data)
	if err != nil {
		return errors.Wrap(err, "unmarshal pull data packet failed")
	}

	p.Protocol = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])

	for i := 0; i < 8; i++ {
		p.GatewayMac[i] = data[4+i]
	}

	return nil
}

func isValidPullDataPacket(data []byte) (bool, error) {
	if len(data) != 12 {
		return false, errors.New("invalid packet: 12 bytes expected")
	}

	return true, nil
}
