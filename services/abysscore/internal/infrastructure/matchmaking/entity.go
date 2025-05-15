package matchmaking

import "time"

// TODO: add link on tech about

type Entity struct {
	id        int
	name      string
	baseScore int // 0 <= this <= 1_000
	startedAt time.Time
}

func NewEntity(id int, name string, baseScore int) *Entity {
	return &Entity{
		id:        id,
		name:      name,
		baseScore: baseScore,
		startedAt: time.Now(),
	}
}

// TODO: add link on tech about

func (e *Entity) additionalScoreForWaiting() int {
	// TODO: implement me
	panic("implement me")
}
