package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type PlayerMatchResult struct {
	ent.Schema
}

const (
	MatchResultMinScore = 0
	MatchResultMaxScore = 600
)

func (PlayerMatchResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.Int("match_id").Immutable().Positive(),
		field.Int("player_id").Immutable().Positive(),

		field.Int("score").
			Range(MatchResultMinScore, MatchResultMaxScore).
			Comment("lower is better: 0 is best, 600 is worst"),

		field.Bool("is_retried").Default(false),

		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (PlayerMatchResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("match", Match.Type).
			Ref("results").
			Unique().
			Field("match_id").
			Required().
			Immutable(),

		edge.To("user", User.Type).
			Field("player_id").
			Unique().
			Immutable().
			Required(),
	}
}
