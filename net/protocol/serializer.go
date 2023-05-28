package protocol

import (
	"reflect"
)

type Packet interface {
	GetId() int
}

func Serialize(packet Packet) PacketBuffer {
	buffer := CreatePacketBuffer(packet.GetId())
	val := reflect.ValueOf(packet).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {

		case reflect.String:
			buffer.WriteString(field.String())
		case reflect.Int16:
			buffer.WriteInt16(int16(field.Int()))
		case reflect.Int:
		case reflect.Int32:
			buffer.WriteInt32(int32(field.Int()))
		case reflect.Int64:
			buffer.WriteInt64(field.Int())
		case reflect.Bool:
			buffer.WriteBool(field.Bool())
		case reflect.Slice:
			buffer.WriteBytes(field.Bytes())
		case reflect.Float64:
			buffer.WriteDouble(field.Float())
		case reflect.Array:
			switch field.Type().Elem().Kind() {
			case reflect.Int:
				buffer.WriteInt32Array(field.Interface().([]int))
			case reflect.Int64:
				buffer.WriteInt64Array(field.Interface().([]int64))
			case reflect.String:
				buffer.WriteStringArray(field.Interface().([]string))
			case reflect.Bool:
				buffer.WriteBoolArray(field.Interface().([]bool))
			case reflect.Float64:
				buffer.WriteDoubleArray(field.Interface().([]float64))
			}
		default:
			panic("unsupported type")
		}
	}

	buffer.Finalize()

	return buffer
}

func Deserialize(buffer PacketBuffer, packet Packet) {
	val := reflect.ValueOf(packet).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(buffer.ReadString())
			case reflect.Int16:
				field.SetInt(int64(buffer.ReadInt16()))
			case reflect.Int:
			case reflect.Int32:
				field.SetInt(int64(buffer.ReadInt32()))
			case reflect.Int64:
				field.SetInt(buffer.ReadInt64())
			case reflect.Bool:
				field.SetBool(buffer.ReadBool())
			case reflect.Slice:
				field.SetBytes(buffer.ReadBytes())
			case reflect.Float64:
				field.SetFloat(buffer.ReadDouble())
			case reflect.Array:
				switch field.Type().Elem().Kind() {
				case reflect.Int:
					field.Set(reflect.ValueOf(buffer.ReadInt32Array()))
				case reflect.Int64:
					field.Set(reflect.ValueOf(buffer.ReadInt64Array()))
				case reflect.String:
					field.Set(reflect.ValueOf(buffer.ReadStringArray()))
				case reflect.Bool:
					field.Set(reflect.ValueOf(buffer.ReadBoolArray()))
				case reflect.Float64:
					field.Set(reflect.ValueOf(buffer.ReadDoubleArray()))
				}

			default:
				panic("unsupported type")
			}
		}
	}
}
