package random

// WeightedItem represents an item value with an associated weight.
type WeightedItem struct {
	Value  interface{}
	Weight int
}

// WeightedItems represents a collection of weighted items with a
// pre-calculated total weight.
type WeightedItems struct {
	items       []WeightedItem
	totalWeight int
}

// MakeWeightedItems creates a slice of WeightedItems and calculates
// the total weight.
func MakeWeightedItems(items []WeightedItem) WeightedItems {
	wi := WeightedItems{
		items: items,
	}

	for _, item := range items {
		wi.totalWeight += item.Weight
	}

	return wi
}

// Choose selects a random value using the weights of each to ensure
// items with higher weights have more of a chance of being selected.
func (wi WeightedItems) Choose() interface{} {
	randomWeight := between64(1, int64(wi.totalWeight))
	for _, i := range wi.items {
		randomWeight -= int64(i.Weight)
		if randomWeight <= 0 {
			return i.Value
		}
	}

	panic("didn't select an item")
}
