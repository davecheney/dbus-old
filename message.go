package dbus

// marshal and unmarshal support

import (
	"bytes"
	"io"
	"reflect"
)

type messageReader struct {
	io.Reader
}

type messageWriter struct {
	io.Writer     // the underlying writer
	pos       int // used to calculate padding
}

func (r *messageReader) readUint32() (uint32, error) {
	return 0, nil
}

// shared pading buffer
var padding = make([]byte, 8)

// pad adds additional padding to the underlying writer
// base on the alignment requested.
func (w *messageWriter) pad(alignment int) error {
	n, err := w.Write(padding[:w.pos%alignment])
	if err != nil {
		return err
	}
	w.pos += n
	return nil
}

func (w *messageWriter) writeBool(b bool) error {
	var v uint32 = 0
	if b {
		v = 1
	}
	return w.writeUint32(v)
}

func (w *messageWriter) writeByte(b byte) error {
	n, err := w.Write([]byte{b})
	if err != nil {
		return err
	}
	w.pos += n
	return nil
}

func (w *messageWriter) writeInt16(i int16) error {
	return w.writeUint16(uint16(i))
}

func (w *messageWriter) writeUint16(u uint16) error {
	if err := w.pad(2); err != nil {
		return err
	}
	b := make([]byte, 2)
	endian.PutUint16(b, u)
	n, err := w.Write(b)
	if err != nil {
		return err
	}
	w.pos += n
	return nil
}

func (w *messageWriter) writeInt32(i int32) error {
	return w.writeUint32(uint32(i))
}

func (w *messageWriter) writeUint32(u uint32) error {
	if err := w.pad(4); err != nil {
		return err
	}
	b := make([]byte, 4)
	endian.PutUint32(b, u)
	n, err := w.Write(b)
	if err != nil {
		return err
	}
	w.pos += n
	return nil
}

func (w *messageWriter) writeInt64(i int64) error {
	return w.writeUint64(uint64(i))
}

func (w *messageWriter) writeUint64(u uint64) error {
	if err := w.pad(8); err != nil {
		return err
	}
	b := make([]byte, 8)
	endian.PutUint64(b, u)
	n, err := w.Write(b)
	if err != nil {
		return err
	}
	w.pos += n
	return nil
}

func (w *messageWriter) writeString(s string) error {
	if err := w.writeUint32(uint32(len(s))); err != nil {
		return err
	}
	n, err := w.Write([]byte(s))
	if err != nil {
		return err
	}
	w.pos += n
	return w.writeByte(0)
}

func marshal(msg interface{}) []byte {
	var b bytes.Buffer
	w := messageWriter{&b, 0}
	v := reflect.ValueOf(msg)
	for i, n := 0, v.NumField(); i < n; i++ {
		field := v.Field(i)
		switch t := field.Type(); t.Kind() {
		case reflect.Bool:
			w.writeBool(field.Bool())
		case reflect.Uint8:
			w.writeByte(byte(field.Uint()))
		case reflect.Int16:
			w.writeInt16(int16(field.Int()))
		case reflect.Uint16:
			w.writeUint16(uint16(field.Uint()))
		case reflect.Int32:
			w.writeInt32(int32(field.Int()))
		case reflect.Uint32:
			w.writeUint32(uint32(field.Uint()))
		case reflect.Int64:
			w.writeInt64(int64(field.Int()))
		case reflect.Uint64:
			w.writeUint64(uint64(field.Uint()))
		case reflect.String:
			w.writeString(field.String())
		default:
			panic("unknown type")
		}
	}
	return b.Bytes()
}

func unmarshal(out interface{}, b []byte) error {
	r := messageReader{bytes.NewBuffer(b)}
	v := reflect.ValueOf(out).Elem()
	//structType := v.Type()
	//var ok bool
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		t := field.Type()
		switch t.Kind() {
		case reflect.Bool:
			v, err := r.readUint32()
			if err != nil {
				return err
			}
			field.SetBool(v == 1)

		}
	}
	return nil
}
