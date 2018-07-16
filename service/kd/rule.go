package kd

import (
	"reflect"
	"strings"
)

const RuleFuncPrefix  = "CheckI"

type Rule struct {
	Fee float64 `json:"fee"`
	Rate float64 `json:"rate"`
	Restdays int `json:"restdays"`
}

type Rules struct {
	R []*Rule
}

func (rules *Rules) Check(item TransferItem) bool {
	for _, rule := range rules.R {
		if rule.Check(item) {
			return true
		}
	}

	return false
}

func InitRule() *Rule {
	return &Rule{}
}

func (rule *Rule) SetFee(fee float64) {
	rule.Fee = fee
}

func (rule *Rule) SetRate(rate float64) {
	rule.Rate = rate
}

func (rule *Rule) SetRestdays(restdays int) {
	rule.Restdays = restdays
}

func (rule *Rule) Check(item TransferItem) bool {
	t := reflect.TypeOf(rule)
	f := reflect.ValueOf(rule)

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(item)

	i := 0
	for i < t.NumMethod() {
		if strings.Contains(t.Method(i).Name, RuleFuncPrefix) {
			res := f.Method(i).Call(params)
			if !res[0].Bool() {
				return false
			}
		}
		i++
	}

	return true
}

func (rule *Rule) CheckIFee(item TransferItem) bool {
	if item.GetFee() > rule.Fee {
		return false
	}
	return true
}

func (rule *Rule) CheckIRate(item TransferItem) bool {
	if item.GetRate() < rule.Rate {
		return false
	}
	return true
}

func (rule *Rule) CheckIRestDays(item TransferItem) bool {
	if item.RestDays > rule.Restdays {
		return false
	}
	return true
}
