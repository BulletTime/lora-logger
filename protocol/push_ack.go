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

// PushACKPacket is used by the server to acknowledge immediately all the
// PUSH_DATA packets received
type PushAckPacket struct {
	Protocol    uint8
	RandomToken uint16
}

func handlePushAck(data []byte) (Packet, error) {
	var packet PushAckPacket

	err := packet.unmarshalData(data)
	if err != nil {
		return nil, errors.Wrap(err, "handle push ack packet failed")
	}

	return &packet, nil
}

func (p *PushAckPacket) Log(ctx log.Interface) {
	ctx.WithFields(log.Fields{
		"protocol":     p.Protocol,
		"random token": p.RandomToken,
	}).Info("PUSH_ACK")
}

func (p *PushAckPacket) unmarshalData(data []byte) error {
	_, err := isValidPushAckPacket(data)
	if err != nil {
		return errors.Wrap(err, "unmarshal push ack packet failed")
	}

	p.Protocol = data[0]
	p.RandomToken = binary.LittleEndian.Uint16(data[1:3])

	return nil
}

func isValidPushAckPacket(data []byte) (bool, error) {
	if len(data) != 4 {
		return false, errors.New("invalid packet: 4 bytes expected")
	}

	return true, nil
}
