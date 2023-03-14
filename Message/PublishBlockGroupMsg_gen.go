package Message

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *PublishBlockGroupMsg) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Height":
			z.Height, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Height")
				return
			}
		case "MinicNum":
			z.MinicNum, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "MinicNum")
				return
			}
		case "BlockGroup":
			z.Data, err = dc.ReadBytes(z.Data)
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PublishBlockGroupMsg) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Height"
	err = en.Append(0x83, 0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Height)
	if err != nil {
		err = msgp.WrapError(err, "Height")
		return
	}
	// write "MinicNum"
	err = en.Append(0xa8, 0x4d, 0x69, 0x6e, 0x69, 0x63, 0x4e, 0x75, 0x6d)
	if err != nil {
		return
	}
	err = en.WriteInt(z.MinicNum)
	if err != nil {
		err = msgp.WrapError(err, "MinicNum")
		return
	}
	// write "BlockGroup"
	err = en.Append(0xaa, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x47, 0x72, 0x6f, 0x75, 0x70)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PublishBlockGroupMsg) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Height"
	o = append(o, 0x83, 0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
	o = msgp.AppendInt(o, z.Height)
	// string "MinicNum"
	o = append(o, 0xa8, 0x4d, 0x69, 0x6e, 0x69, 0x63, 0x4e, 0x75, 0x6d)
	o = msgp.AppendInt(o, z.MinicNum)
	// string "BlockGroup"
	o = append(o, 0xaa, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x47, 0x72, 0x6f, 0x75, 0x70)
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PublishBlockGroupMsg) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Height":
			z.Height, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Height")
				return
			}
		case "MinicNum":
			z.MinicNum, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MinicNum")
				return
			}
		case "BlockGroup":
			z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PublishBlockGroupMsg) Msgsize() (s int) {
	s = 1 + 7 + msgp.IntSize + 9 + msgp.IntSize + 11 + msgp.BytesPrefixSize + len(z.Data)
	return
}