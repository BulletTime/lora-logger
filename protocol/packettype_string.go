// Code generated by "stringer -type=PacketType"; DO NOT EDIT.

package protocol

import "strconv"

const _PacketType_name = "PushDataPushAckPullDataPullRespPullAckTXAck"

var _PacketType_index = [...]uint8{0, 8, 15, 23, 31, 38, 43}

func (i PacketType) String() string {
	if i >= PacketType(len(_PacketType_index)-1) {
		return "PacketType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PacketType_name[_PacketType_index[i]:_PacketType_index[i+1]]
}
