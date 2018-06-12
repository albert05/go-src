package model

import "strconv"

type MapModel map[string]interface{}

const DEFAULT_STRVAL = ""
const DEFAULT_FLTVAL = 0

func (this MapModel) GetAttrString(k string) string {
	if val, ok := this[k];ok {
		if vStr, ok := val.(string); ok {
			return vStr
		}
	}

	return DEFAULT_STRVAL
}

func (this MapModel) GetAttrFloat(k string) float64 {
	str := this.GetAttrString(k)

	flt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return DEFAULT_FLTVAL
	}

	return flt
}

func (this MapModel) GetAttrInt(k string) int {
	str := this.GetAttrString(k)

	i, err := strconv.Atoi(str)
	if err != nil {
		return DEFAULT_FLTVAL
	}

	return i
}
