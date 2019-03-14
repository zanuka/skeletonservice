package skelsvc

type ColVal struct {
	Columns []string
	Values  []interface{}
}

type ValueData struct {
	Value interface{}
}

func (vd *ValueData) ToString() string {
	return interfaceToString(vd.Value)
}

func (cv *ColVal) FindValueForCol(key string) *ValueData {

	var val interface{}

	for i, col := range cv.Columns {

		if key == col {
			val = cv.Values[i]

		}
	}

	vd := ValueData{Value: val}

	return &vd

}
