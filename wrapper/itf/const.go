package itf

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCanceled  JobStatus = "canceled"
	JobCompleted JobStatus = "completed"
	JobError     JobStatus = "error"
	JobOK        JobStatus = "ok"
)

var JobStateMap = map[JobStatus]struct{}{
	JobPending:   {},
	JobRunning:   {},
	JobCanceled:  {},
	JobCompleted: {},
	JobError:     {},
	JobOK:        {},
}
