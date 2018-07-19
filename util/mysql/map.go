package mysql

import "strconv"

type MapModel map[string]interface{}

const DefaultSTRVAL = ""
const DefaultFLTVAL = 0

func (this MapModel) GetAttrString(k string) string {
	if val, ok := this[k];ok {
		if vStr, ok := val.(string); ok {
			return vStr
		}
	}

	return DefaultSTRVAL
}

func (this MapModel) GetAttrFloat(k string) float64 {
	str := this.GetAttrString(k)

	flt, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return DefaultFLTVAL
	}

	return flt
}

func (this MapModel) GetAttrInt(k string) int {
	str := this.GetAttrString(k)

	i, err := strconv.Atoi(str)
	if err != nil {
		return DefaultFLTVAL
	}

	return i
}
