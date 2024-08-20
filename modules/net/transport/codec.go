package transport

import (
	packet2 "github.com/orbit-w/meteor/bases/net/packet"
)

/*
   @Author: orbit-w
   @File: codec
   @2023 12月 周六 20:41
*/

func packHeadByte(data []byte, mt int8) packet2.IPacket {
	writer := packet2.WriterP(1 + len(data))
	writer.WriteInt8(mt)
	if data != nil && len(data) > 0 {
		writer.Write(data)
	}
	return writer
}

func packHeadByteP(pack packet2.IPacket, mt int8) packet2.IPacket {
	data := pack.Remain()
	writer := packet2.WriterP(1 + len(data))
	writer.WriteInt8(mt)
	if pack != nil {
		if len(data) > 0 {
			writer.Write(data)
		}
	}
	return writer
}

func unpackHeadByte(pack packet2.IPacket, handle func(h int8, data []byte)) error {
	defer packet2.Return(pack)
	head, err := pack.ReadInt8()
	if err != nil {
		return err
	}

	data := pack.CopyRemain()
	handle(head, data)
	return nil
}
