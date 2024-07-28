package stream

import "github.com/orbit-w/meteor/bases/packet"

/*
   @Author: orbit-w
   @File: streamer
   @2024 7月 周日 16:22
*/

type Codec struct{}

type StreamingMsg struct {
	Type     int8
	End      bool
	StreamId int64
	Data     packet.IPacket
}

func (f *Codec) Encode(msg *StreamingMsg) packet.IPacket {
	w := packet.Writer()
	w.WriteInt8(msg.Type)
	w.WriteBool(msg.End)
	w.WriteInt64(msg.StreamId)
	if data := msg.Data; data != nil {
		msg.Data = nil
		w.Write(data.Remain())
		data.Return()
	}
	return w
}

func (f *Codec) Decode(data packet.IPacket) (StreamingMsg, error) {
	defer data.Return()
	msg := StreamingMsg{}
	ft, err := data.ReadInt8()
	if err != nil {
		return msg, err
	}

	end, err := data.ReadBool()
	if err != nil {
		return msg, err
	}

	sId, err := data.ReadUint64()
	if err != nil {
		return msg, err
	}

	msg.StreamId = int64(sId)
	msg.Type = ft
	msg.End = end
	if len(data.Remain()) > 0 {
		reader := packet.Reader(data.Remain())
		msg.Data = reader
	}
	return msg, nil
}
