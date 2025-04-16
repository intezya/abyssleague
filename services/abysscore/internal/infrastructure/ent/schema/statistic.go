package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

const globalStatisticPeriod = 0

type Statistic struct {
	ent.Schema
}

func (Statistic) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.Int("user_id").Immutable(),

		field.Enum("type").Values("global").Default("global"),
		field.Int("period").Default(globalStatisticPeriod).Positive(),

		field.Int("xp").Default(0).Positive(),

		field.Int("match_count").Default(0).Positive(),
		field.Int("wins_count").Default(0).Positive(),
		field.Int("loses_count").Default(0).Positive(),
		field.Int("draws_count").Default(0).Positive(),

		field.Int("result_time").Default(0).Positive(),
		field.Int("retry_time").Default(0).Positive(),
		field.Int("retry_count").Default(0).Positive(),

		field.Int("best_result_time").Default(0).Positive(),
		field.Int("best_retry_count").Default(0).Positive(),
		field.Int("best_match_time").Default(0).Positive(),

		field.Int("worst_result_time").Default(0).Positive(),
		field.Int("worst_retry_count").Default(0).Positive(),
		field.Int("worst_match_time").Default(0).Positive(),

		field.Int("max_win_streak").Default(0).Positive(),
		field.Int("max_lose_streak").Default(0).Positive(),

		field.Int("max_login_streak").Default(0).Positive(),

		field.Int("search_score").Default(0).Positive(),

		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Statistic) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("statistics").
			Unique().
			Required().
			Immutable().
			Field("user_id"),
	}
}
