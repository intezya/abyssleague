package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type Match struct {
	ent.Schema
}

func (Match) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.Int("player1_id").Immutable(),
		field.Int("player2_id").Immutable(),

		field.Int("player1_penalty_time").Default(0).Positive(),
		field.Int("player2_penalty_time").Default(0).Positive(),

		field.Enum("status").Values(
			"characters_reveal",
			"waiting_for_ready",
			"drafting",
			"matching",
			"finished",
		).Default("characters_reveal"),

		field.Enum("result").Values(
			"player1_win",
			"player2_win",
			"draw",
		).Optional().Nillable(),

		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("changed_to_current_status_at").Default(time.Now).Immutable(),
	}
}

func (Match) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("player1", User.Type).
			Unique().
			Field("player1_id").
			Required().
			Immutable(),

		edge.To("player2", User.Type).
			Unique().
			Field("player2_id").
			Required().
			Immutable(),

		edge.To("results", PlayerMatchResult.Type),
	}
}
