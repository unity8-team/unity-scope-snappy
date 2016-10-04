package scopes

import (
	"sort"
)

// OptionSelectorFilter is used to implement single-select or multi-select filters.
type OptionSelectorFilter struct {
	filterWithOptions
	Label       string
	MultiSelect bool
}

// NewOptionSelectorFilter creates a new option filter.
func NewOptionSelectorFilter(id, label string, multiSelect bool) *OptionSelectorFilter {
	return &OptionSelectorFilter{
		filterWithOptions: filterWithOptions{
			filterBase: filterBase{
				Id:           id,
				DisplayHints: FilterDisplayDefault,
				FilterType:   "option_selector",
			},
		},
		Label:       label,
		MultiSelect: multiSelect,
	}
}

type optionSort struct {
	Options []interface{}
}

func (s optionSort) Len() int {
	return len(s.Options)
}

func (s optionSort) Less(i, j int) bool {
	return s.Options[i].(string) < s.Options[j].(string)
}

func (s optionSort) Swap(i, j int) {
	s.Options[i], s.Options[j] = s.Options[j], s.Options[i]
}

// UpdateState updates the value of a particular option in the filter state.
func (f *OptionSelectorFilter) UpdateState(state FilterState, optionId string, active bool) {
	if !f.isValidOption(optionId) {
		panic("invalid option ID")
	}
	// For single-select filters, clear the previous state when
	// setting a new active option.
	if active && !f.MultiSelect {
		delete(state, f.Id)
	}
	// If the state isn't in a form we expect, treat it as empty
	selected, _ := state[f.Id].([]interface{})
	sort.Sort(optionSort{selected})
	pos := sort.Search(len(selected), func(i int) bool { return selected[i].(string) >= optionId })
	if active {
		if pos == len(selected) {
			selected = append(selected, optionId)
		} else if pos < len(selected) && selected[pos] != optionId {
			selected = append(selected[:pos], append([]interface{}{optionId}, selected[pos:]...)...)
		}
	} else {
		if pos < len(selected) {
			selected = append(selected[:pos], selected[pos+1:]...)
		}
	}
	state[f.Id] = selected
}

func (f *OptionSelectorFilter) serializeFilter() map[string]interface{} {
	v := f.filterBase.serializeFilter()
	v["label"] = f.Label
	v["multi_select"] = f.MultiSelect
	v["options"] = f.Options
	return v
}
