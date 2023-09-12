//go:build ignore

package codec

type Encoder interface {
	EncodeStruct(name string) StructEncoder
}

type StructEncoder interface {
	EncodeField(name string, v any) error
	End() error
}

type FieldsEncoder interface{}

type ElementsEncoder interface{}

type Widget struct {
	Name string
}

func (w *Widget) Encode(enc Encoder) error {
	s := enc.EncodeStruct("widget")
	s.EncodeField("name", &w.Name)
	return s.End()
}

func (w *Widget) DecodeFields(dec FieldsDecoder) error {
	for dec.More() {
		f, err := dec.Next()
		if err != nil {
			return err
		}
		switch f.Name() {
		case "name":
			dec.Decode(&w.Name)
		}
	}
	return dec.End()
}
