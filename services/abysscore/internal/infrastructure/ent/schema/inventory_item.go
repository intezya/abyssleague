package schema

import (
	"entgo.io/ent/schema/index"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const systemIssuerId = 0

type InventoryItem struct {
	ent.Schema
}

func (InventoryItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.Int("user_id").Immutable(),
		field.Int("item_id").Immutable(),

		field.Int("received_from_id").
			Default(systemIssuerId).
			Immutable().
			Nillable(), // nil if from trade
		// TODO: relationship with trades

		field.Time("obtained_at").Default(time.Now),
	}
}

func (InventoryItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("items").
			Field("user_id").
			Unique().
			Required().
			Immutable(),

		edge.From("item", GameItem.Type).
			Ref("inventory_items").
			Field("item_id").
			Unique().
			Required().
			Immutable(),
	}
}

func (InventoryItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "item_id"),
	}
}
