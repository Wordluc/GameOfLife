package core

type Job int8

const (
	FARMER Job = iota
	MINER
)

type Person struct {
	id       int
	job      Job
	position *BaseCell
}

func newPerson(job Job) Person {
	p := Person{
		id:  ID_PEOPLE,
		job: job,
	}
	ID_PEOPLE++
	return p
}
