package enum

type IConvertor interface {
	Has(any) bool
	ToUint32(string) uint32
	ToString(uint32) string
}

type Convertor struct {
	names   map[string]int32
	numbers map[int32]string
}

func WrapConvertor(n map[string]int32, i map[int32]string) *Convertor {
	return &Convertor{
		names:   n,
		numbers: i,
	}
}

func (d *Convertor) Has(val any) bool {
	var ok bool
	switch vv := val.(type) {
	case uint32:
		_, ok = d.numbers[int32(vv)]
	case int32:
		_, ok = d.numbers[vv]
	case string:
		_, ok = d.names[vv]
	}
	return ok
}

func (d *Convertor) ToUint32(s string) uint32 {
	return uint32(d.names[s])
}

func (d *Convertor) ToString(i uint32) string {
	return d.numbers[int32(i)]
}
