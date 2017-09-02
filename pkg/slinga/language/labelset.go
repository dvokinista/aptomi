package language

import "reflect"

// LabelSet defines the set of labels that will be manipulated
type LabelSet struct {
	Labels map[string]string
}

// NewLabelSetEmpty creates a new empty LabelSet
func NewLabelSetEmpty() LabelSet {
	return LabelSet{Labels: make(map[string]string)}
}

// NewLabelSet creates a new LabelSet from a given set of labels
func NewLabelSet(labels map[string]string) LabelSet {
	return LabelSet{Labels: labels}
}

// NewLabelSet creates a new LabelSet from a given set of labels
func (src LabelSet) MakeCopy() LabelSet {
	return NewLabelSetEmpty().addToCurrent(src)
}

// AddLabels merges two sets of labels and returns a new merged set
func (src LabelSet) AddLabels(addSet LabelSet) LabelSet {
	return src.MakeCopy().addToCurrent(addSet)
}

// ApplyTransform applies set of transformations to labels
func (src LabelSet) ApplyTransform(ops LabelOperations) LabelSet {
	result := src.MakeCopy()

	if ops != nil {
		// set labels
		for k, v := range ops["set"] {
			result.Labels[k] = v
		}

		// remove labels
		for k := range ops["remove"] {
			delete(result.Labels, k)
		}
	}

	return result
}

// Adds labels to the current set and returns it (without creating a copy)
func (src LabelSet) addToCurrent(addSet LabelSet) LabelSet {
	for k, v := range addSet.Labels {
		src.Labels[k] = v
	}
	return src
}

// Equal compares two labels sets. If one is nil and another one is empty, it will return true as well
// This method ignores IsSecret for now
func (src LabelSet) Equal(dst LabelSet) bool {
	if len(src.Labels) == 0 && len(dst.Labels) == 0 {
		return true
	}
	return reflect.DeepEqual(src.Labels, dst.Labels)
}
