package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type BannedHardwareID struct {
	ent.Schema
}

func (BannedHardwareID) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.String("hardware_id").Unique().Immutable(),
		field.Time("created_at").Immutable(),
		field.String("ban_reason").Nillable().Optional(),
	}
}

func (BannedHardwareID) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("hardware_id").Unique(),
	}
}
