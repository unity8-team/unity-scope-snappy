package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestValueSliderFilter(c *C) {
	labels := scopes.ValueSliderLabels{
		MinLabel: "min",
		MaxLabel: "max",
		ExtraLabels: []scopes.ValueSliderExtraLabel{
			{50, "middle"},
		},
	}
	filter1 := scopes.NewValueSliderFilter("f1", 10.0, 100.0, 50, labels)
	c.Check("f1", Equals, filter1.Id)
	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayDefault)
	c.Check(filter1.DefaultValue, Equals, 50.0)
	c.Check(filter1.Min, Equals, 10.0)
	c.Check(filter1.Max, Equals, 100.0)
	c.Check(filter1.Labels, DeepEquals, labels)

	fstate := make(scopes.FilterState)
	value := filter1.Value(fstate)
	c.Check(value, Equals, 50.0)

	err := filter1.UpdateState(fstate, 30.5)
	c.Check(err, IsNil)
	value = filter1.Value(fstate)
	c.Check(value, Equals, 30.5)

	err = filter1.UpdateState(fstate, 44.5)
	c.Check(err, IsNil)
	value = filter1.Value(fstate)
	c.Check(value, Equals, 44.5)

	err = filter1.UpdateState(fstate, 3545.33)
	c.Check(err, Not(Equals), nil)
	value = filter1.Value(fstate)
	c.Check(value, Equals, 44.5)
}
