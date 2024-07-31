package mux

import "github.com/orbit-w/meteor/bases/packet"

/*
   @Author: orbit-w
   @File: streamer
   @2024 7月 周日 16:22
*/

type Codec struct{}

type Msg struct {
	Type int8
	End  bool
	Id   int64
	Data []byte
}

func (f *Codec) Encode(msg *Msg) packet.IPacket {
	w := packet.Writer()
	w.WriteInt8(msg.Type)
	w.WriteBool(msg.End)
	w.WriteInt64(msg.Id)
	if data := msg.Data; data != nil || len(data) > 0 {
		msg.Data = nil
		w.Write(data)
	}
	return w
}

func (f *Codec) Decode(data []byte) (Msg, error) {
	reader := packet.Reader(data)
	defer reader.Return()
	msg := Msg{}
	ft, err := reader.ReadInt8()
	if err != nil {
		return msg, err
	}

	end, err := reader.ReadBool()
	if err != nil {
		return msg, err
	}

	sId, err := reader.ReadUint64()
	if err != nil {
		return msg, err
	}

	msg.Id = int64(sId)
	msg.Type = ft
	msg.End = end
	if len(reader.Remain()) > 0 {
		msg.Data = reader.CopyRemain()
	}
	return msg, nil
}
