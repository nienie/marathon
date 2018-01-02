package stats

//DataCollector An object that collects new values incrementally.
type DataCollector interface {
	//NoteValue Adds a value to the collected data.
	//This must run very quickly, and so can safely be called in time-critical code.
	NoteValue(val float64)
}
