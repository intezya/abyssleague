package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type MatchResult struct {
	ent.Schema
}

func (MatchResult) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.Int("match_id").Immutable().Positive(),
		field.Int("player_id").Immutable().Positive(),

		field.Int("value").Range(0, 600),

		field.Bool("is_retry").Default(false),

		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (MatchResult) Edges() []ent.Edge {
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
