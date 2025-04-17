package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

const systemIssuerId = 0

type UserItem struct {
	ent.Schema
}

func (UserItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.Int("user_id").Immutable(),
		field.Int("item_id").Immutable(),

		field.Int("received_from_id").Default(systemIssuerId).Positive(), // todo: relationship with trades

		field.Time("obtained_at").Default(time.Now),
	}
}

func (UserItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("items").
			Field("user_id").
			Unique().
			Required().
			Immutable(),

		edge.From("item", GameItem.Type).
			Ref("user_items").
			Field("item_id").
			Unique().
			Required().
			Immutable(),
	}
}

func (UserItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("item_id").Unique(),
		index.Fields("user_id").Unique(),
	}
}
