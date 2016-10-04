package scopes

import (
	"errors"
	"fmt"
)

// ValueSliderFilter is a value slider filter that allows for selecting a value within a given range.
type ValueSliderFilter struct {
	filterBase
	DefaultValue float64
	Min          float64
	Max          float64
	Labels       ValueSliderLabels
}

type ValueSliderLabels struct {
	MinLabel    string
	MaxLabel    string
	ExtraLabels []ValueSliderExtraLabel
}

type ValueSliderExtraLabel struct {
	Value float64
	Label string
}

// NewValueSliderFilter creates a new value slider filter.
func NewValueSliderFilter(id string, min, max, defaultValue float64, labels ValueSliderLabels) *ValueSliderFilter {
	filter := &ValueSliderFilter{
		filterBase: filterBase{
			Id:           id,
			DisplayHints: FilterDisplayDefault,
			FilterType:   "value_slider",
		},
		Min:          min,
		Max:          max,
		DefaultValue: defaultValue,
		Labels:       labels,
	}
	filter.validate()
	return filter
}

func (f *ValueSliderFilter) validate() {
	if f.Min >= f.Max {
		panic("Invalid range for value slider filter")
	}
	last := f.Min
	labels := map[string]bool{
		f.Labels.MinLabel: true,
		f.Labels.MaxLabel: true,
	}
	// check that values of extra labels grow, i.e. v1 < v2 < v3 ...
	// and labels are unique and not empty
	for _, l := range f.Labels.ExtraLabels {
		if l.Value <= last {
			panic("Extra label for value slider filter out of sequence")
		}
		last = l.Value
		if l.Label == "" {
			panic("Extra labels cannot be empty")
		}
		if labels[l.Label] {
			panic("Multiple definitions for extra label")
		}
		labels[l.Label] = true
	}
	if f.Max <= last {
		panic("Last extra label value greater than maximum value")
	}
}

// Value gets value of this filter from filter state object.
// If the value is not set for the filter it returns false as the second return statement,
// it returns true otherwise
func (f *ValueSliderFilter) Value(state FilterState) float64 {
	if value, ok := state[f.Id].(float64); ok {
		return value
	}
	return f.DefaultValue
}

// UpdateState updates the value of the filter to the given value
func (f *ValueSliderFilter) UpdateState(state FilterState, value float64) error {
	if value < f.Min || value > f.Max {
		return errors.New(fmt.Sprintf("ValueSliderFilter:UpdateState: value %f outside of allowed range (%f,%f)", value, f.Min, f.Max))
	}
	state[f.Id] = value
	return nil
}

func (f *ValueSliderFilter) serializeFilter() map[string]interface{} {
	v := f.filterBase.serializeFilter()
	v["min"] = marshalFloat(f.Min)
	v["max"] = marshalFloat(f.Max)
	v["default"] = marshalFloat(f.DefaultValue)
	extra := make([]interface{}, 0, 2*len(f.Labels.ExtraLabels))
	for _, l := range f.Labels.ExtraLabels {
		extra = append(extra, marshalFloat(l.Value), l.Label)
	}
	v["labels"] = map[string]interface{}{
		"min_label":    f.Labels.MinLabel,
		"max_label":    f.Labels.MaxLabel,
		"extra_labels": extra,
	}
	return v
}
