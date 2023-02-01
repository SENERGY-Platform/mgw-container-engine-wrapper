package itf

const (
	JobPending   JobState = "pending"
	JobRunning   JobState = "running"
	JobCanceled  JobState = "canceled"
	JobCompleted JobState = "completed"
	JobError     JobState = "error"
	JobOK        JobState = "ok"
)

var JobStateMap = map[JobState]struct{}{
	JobPending:   {},
	JobRunning:   {},
	JobCanceled:  {},
	JobCompleted: {},
	JobError:     {},
	JobOK:        {},
}
