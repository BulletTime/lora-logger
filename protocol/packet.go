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

//go:generate stringer -type=PacketType

package protocol

import (
	"fmt"

	"github.com/apex/log"
	"github.com/pkg/errors"
)

// PacketType defines the packet type.
type PacketType byte

// Available packet types
const (
	PushData PacketType = iota
	PushAck
)

const (
	ProtoVersion1 uint8 = 0x01
	ProtoVersion2 uint8 = 0x02
)

// HandlePacket will check the packet and try to handle it accordingly.
// If the packet is recognized, it will extract the data and write it to the log.
// If there was an error of some sort, the raw data with the error will be logged.
func HandlePacket(data []byte) error {
	_, err := isValidPacket(data)
	if err != nil {
		return errors.Wrap(err, "handle packet failed")
	}

	pType := PacketType(data[3])

	switch pType {
	case PushData:
		return handlePushData(data)
	case PushAck:
		return handlePushAck(data)
	default:
		return errors.New(fmt.Sprintf("unknown packet type: %s", pType))
	}
}

func isProtocolSupported(protocol uint8) bool {
	if protocol == ProtoVersion1 || protocol == ProtoVersion2 {
		log.WithField("protocol", protocol).Debug("supported protocol")
		return true
	}
	return false
}

func isValidPacket(data []byte) (bool, error) {
	if len(data) < 4 {
		return false, errors.New("invalid packet: less than 4 bytes")
	}

	if !isProtocolSupported(data[0]) {
		return false, errors.New("invalid protocol")
	}

	log.WithField("data", data).Debug("valid packet")

	return true, nil
}
