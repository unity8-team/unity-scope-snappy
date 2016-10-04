package scopes

import (
	"errors"
	"fmt"
)

// RangeInputFilter is a range filter which allows a start and end value to be entered by user, and any of them is optional.
type RangeInputFilter struct {
	filterBase
	DefaultStartValue interface{}
	DefaultEndValue   interface{}
	StartPrefixLabel  string
	StartPostfixLabel string
	EndPrefixLabel    string
	EndPostfixLabel   string
	CentralLabel      string
}

func checkRangeValidType(value interface{}) bool {
	switch value.(type) {
	case int, float64, nil:
		return true
	default:
		return false
	}
}

// NewRangeInputFilter creates a new range input filter.
func NewRangeInputFilter(id string, defaultStartValue, defaultEndValue interface{}, startPrefixLabel, startPostfixLabel, endPrefixLabel, endPostfixLabel, centralLabel string) *RangeInputFilter {
	if !checkRangeValidType(defaultStartValue) {
		panic("bad type for defaultStartValue")
	}
	if !checkRangeValidType(defaultEndValue) {
		panic("bad type for defaultEndValue")
	}
	return &RangeInputFilter{
		filterBase: filterBase{
			Id:           id,
			DisplayHints: FilterDisplayDefault,
			FilterType:   "range_input",
		},
		DefaultStartValue: defaultStartValue,
		DefaultEndValue:   defaultEndValue,
		StartPrefixLabel:  startPrefixLabel,
		StartPostfixLabel: startPostfixLabel,
		EndPrefixLabel:    endPrefixLabel,
		EndPostfixLabel:   endPostfixLabel,
		CentralLabel:      centralLabel,
	}
}

// StartValue gets the start value of this filter from filter state object.
// If the value is not set for the filter it returns false as the second return statement,
// it returns true otherwise
func (f *RangeInputFilter) StartValue(state FilterState) (float64, bool) {
	var start float64
	var ok bool
	slice_interface, ok := state[f.Id].([]interface{})
	if ok {
		if len(slice_interface) != 2 {
			// something went really bad.
			// we should have just 2 values
			panic("RangeInputFilter:StartValue unexpected number of values found.")
		}

		switch v := slice_interface[0].(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case nil:
			return 0, false
		default:
			panic("RangeInputFilter:StartValue Unknown value type")
		}
	} else {
		switch v := f.DefaultStartValue.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case nil:
			return 0, false
		}
	}
	return start, ok
}

// EndValue gets the end value of this filter from filter state object.
// If the value is not set for the filter it returns false as the second return statement,
// it returns true otherwise
func (f *RangeInputFilter) EndValue(state FilterState) (float64, bool) {
	var end float64
	var ok bool
	slice_interface, ok := state[f.Id].([]interface{})
	if ok {
		if len(slice_interface) != 2 {
			// something went really bad.
			// we should have just 2 values
			panic("RangeInputFilter:EndValue unexpected number of values found.")
		}

		switch v := slice_interface[1].(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case nil:
			return 0, false
		default:
			panic("RangeInputFilter:EndValue Unknown value type")
		}
	} else {
		switch v := f.DefaultEndValue.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case nil:
			return 0, false
		}
	}
	return end, ok
}

func convertToFloat(value interface{}) float64 {
	if value != nil {
		fVal, ok := value.(float64)
		if !ok {
			iVal, ok := value.(int)
			if !ok {
				panic(fmt.Sprint("RangeInputFilter:convertToFloat unexpected type for given value %v", value))
			}
			return float64(iVal)
		}
		return fVal
	} else {
		panic("RangeInputFilter:convertToFloat nil values are not accepted")
	}
}

// UpdateState updates the value of the filter
func (f *RangeInputFilter) UpdateState(state FilterState, start, end interface{}) error {
	if !checkRangeValidType(start) {
		return errors.New("RangeInputFilter:UpdateState: Bad type for start value. Valid types are int float64 and nil")
	}
	if !checkRangeValidType(end) {
		return errors.New("RangeInputFilter:UpdateState: Bad type for end value. Valid types are int float64 and nil")
	}

	if start == nil && end == nil {
		// remove the state
		delete(state, f.Id)
		return nil
	}
	if start != nil && end != nil {
		fStart := convertToFloat(start)
		fEnd := convertToFloat(end)
		if fStart >= fEnd {
			return errors.New(fmt.Sprintf("RangeInputFilter::UpdateState(): start_value %v is greater or equal to end_value %v for filter %s", start, end, f.Id))
		}
	}
	state[f.Id] = []interface{}{start, end}
	return nil
}

func (f *RangeInputFilter) serializeFilter() map[string]interface{} {
	v := f.filterBase.serializeFilter()
	v["default_start_value"] = f.DefaultStartValue
	v["default_end_value"] = f.DefaultEndValue
	v["start_prefix_label"] = f.StartPrefixLabel
	v["start_postfix_label"] = f.StartPostfixLabel
	v["end_prefix_label"] = f.EndPrefixLabel
	v["end_postfix_label"] = f.EndPostfixLabel
	v["central_label"] = f.CentralLabel
	return v
}
