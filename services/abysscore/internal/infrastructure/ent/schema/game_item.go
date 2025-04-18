package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type GameItem struct {
	ent.Schema
}

func (GameItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.String("name").NotEmpty(),
		field.String("collection").NotEmpty(),
		field.Int("type").Positive(),
		field.Int("rarity").Positive(),
	}
}

func (GameItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_items", UserItem.Type),
	}
}
