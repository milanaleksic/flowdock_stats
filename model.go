package main

type Stat struct {
	numberOfAppearances int
	numberOfEdits       int
	name                string
	words               int
}

type StatByNumberOfAppearances []Stat

func (a StatByNumberOfAppearances) Len() int      { return len(a) }
func (a StatByNumberOfAppearances) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a StatByNumberOfAppearances) Less(i, j int) bool {
	return a[i].numberOfAppearances < a[j].numberOfAppearances
}
