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

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// PullACKPacket is used by the server to confirm that the network route is
// open and that the server can send PULL_RESP packets at any time.
type PullAckPacket struct {
	Protocol    uint8
	RandomToken uint16
}

func handlePullAck(data []byte) (Packet, error) {
	var packet PullAckPacket

	err := packet.unmarshalData(data)
	if err != nil {
		return nil, errors.Wrap(err, "handle pull ack packet failed")
	}

	return &packet, nil
}

func (p *PullAckPacket) Log(ctx log.Interface) {
	ctx.WithFields(log.Fields{
		"protocol":     p.Protocol,
		"random token": p.RandomToken,
	}).Info("PULL_ACK")
}

func (p *PullAckPacket) unmarshalData(data []byte) error {
	_, err := isValidPullAckPacket(data)
	if err != nil {
		return errors.Wrap(err, "unmarshal pull ack packet failed")
	}

	p.Protocol = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])

	return nil
}

func isValidPullAckPacket(data []byte) (bool, error) {
	if len(data) != 4 {
		return false, errors.New("invalid packet: 4 bytes expected")
	}

	return true, nil
}
