package main

type stat struct {
	numberOfAppearances int
	numberOfEdits       int
	name                string
	words               int
}

type statByNumberOfAppearances []stat

func (a statByNumberOfAppearances) Len() int      { return len(a) }
func (a statByNumberOfAppearances) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a statByNumberOfAppearances) Less(i, j int) bool {
	return a[i].numberOfAppearances < a[j].numberOfAppearances
}
