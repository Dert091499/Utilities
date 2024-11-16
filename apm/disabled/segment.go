package disabled

import "github.com/Dert091499/Utilities/apm"

type (
	segment struct{}
)

func (s *segment) AddAttribute(key string, val interface{}) {}

func (s *segment) End() {}

func NewSegment() apm.Segment {
	return &segment{}
}
