package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"slices"
)

type ID_PERSON int

var CURRENT_ID_PEOPLE ID_PERSON = 1

type bindingCellToPeople map[common.Vec[int32]]map[ID_NATION][]ID_PERSON

func (b bindingCellToPeople) checkConsitancy(idNation ID_NATION, pos common.Vec[int32]) {
	atPos := b[pos]
	if atPos == nil {
		b[pos] = map[ID_NATION][]ID_PERSON{}
	}
	inNation := b[pos][idNation]
	if inNation == nil {
		b[pos][idNation] = make([]ID_PERSON, 0)
	}
}
func (b bindingCellToPeople) AddNewPerson(p *Person, idNation ID_NATION, pos common.Vec[int32]) {
	b.checkConsitancy(idNation, pos)
	b[pos][idNation] = append(b[pos][idNation], p.id)
	p.pos = pos
}

func (b bindingCellToPeople) MovePerson(p *Person, idNation ID_NATION, from, to common.Vec[int32]) {
	b.checkConsitancy(idNation, from)
	b.checkConsitancy(idNation, to)
	b[from][idNation] = slices.DeleteFunc(b[from][idNation], func(id ID_PERSON) bool { return p.id == id })
	b[to][idNation] = append(b[to][idNation], p.id)
}

type WorldPopulation struct {
	*World
	people              []*Person
	toRunPathFinding    bool
	bindingCellToPeople bindingCellToPeople
}

func NewWorldPopulation(w *World) WorldPopulation {
	return WorldPopulation{
		World:               w,
		bindingCellToPeople: make(bindingCellToPeople),
	}
}

func (w *WorldPopulation) GetPerson(id ID_PERSON) (res *Person) {
	v := w.GetPeopleCustom(func(p Person) bool { return p.id == id })
	if len(v) == 0 {
		return nil
	}
	return v[0]
}

func (w *WorldPopulation) GetPeopleInsideCell(pos common.Vec[int32], idNation *ID_NATION) (res []ID_PERSON) {
	if idNation != nil {
		return w.bindingCellToPeople[pos][*idNation]
	}
	for _, nation := range w.World.idNations {
		res = append(res, w.bindingCellToPeople[pos][nation]...)
	}
	return res
}

func (w *WorldPopulation) GetPeopleCustom(check func(p Person) bool) (res []*Person) {
	if check == nil {
		check = func(p Person) bool { return true }
	}
	for i := range w.people {
		if check(*w.people[i]) {
			res = append(res, w.people[i])
		}
	}
	return res
}

func (w *WorldPopulation) newPerson(job Job, where common.Vec[int32], idNation ID_NATION) *Person {
	p := new(newPerson(job, idNation, where))
	w.people = append(w.people, p)
	w.toRunPathFinding = true
	w.bindingCellToPeople.AddNewPerson(p, idNation, where)
	return p
}

func (w *WorldPopulation) movePersonToGoal(person *Person) error {
	if person.paths == nil {
		return nil
	}
	from := person.paths.GetBack(1)
	if from != nil && !from.IsEqual(person.pos) {
		return errors.New("Error initial position")
	}
	to, end := person.paths.Denqueue()
	if end {
		person.Status = WORKING
		return nil
	}
	person.Status = MOVING
	w.bindingCellToPeople.MovePerson(person, person.idNation, *from, to)
	person.Touch()
	return nil
}

func (w *WorldPopulation) movePopulationToGoals(idNation ID_NATION) (err error) {
	for _, populations := range w.bindingCellToPeople {
		if people := populations[idNation]; people != nil {
			for _, id := range people {
				person := w.GetPerson(id)
				if person == nil {
					return fmt.Errorf("person not found %v", id)
				}
				err = w.movePersonToGoal(person)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
