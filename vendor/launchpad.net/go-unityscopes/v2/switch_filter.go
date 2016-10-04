package scopes

// SwitchFilter is a simple on/off switch filter.
type SwitchFilter struct {
	filterBase
	Label string
}

// NewSwitchFilter creates a new switch filter.
func NewSwitchFilter(id, label string) *SwitchFilter {
	return &SwitchFilter{
		filterBase: filterBase{
			Id:           id,
			DisplayHints: FilterDisplayDefault,
			FilterType:   "switch",
		},
		Label: label,
	}
}

func (f *SwitchFilter) IsOn(state FilterState) bool {
	value, ok := state[f.Id]
	if ok {
		return value.(bool)
	} else {
		return false
	}
	return true
}

// UpdateState updates the value of the filter to on/off
func (f *SwitchFilter) UpdateState(state FilterState, value bool) {
	state[f.Id] = value
}

func (f *SwitchFilter) serializeFilter() map[string]interface{} {
	v := f.filterBase.serializeFilter()
	v["label"] = f.Label
	return v
}
