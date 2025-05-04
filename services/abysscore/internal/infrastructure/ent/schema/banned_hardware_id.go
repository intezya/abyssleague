package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

type BannedHardwareID struct {
	ent.Schema
}

func (BannedHardwareID) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.String("hardware_id").Unique().Immutable().Sensitive(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.String("ban_reason").Nillable().Optional(),
	}
}

func (BannedHardwareID) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("hardware_id").Unique(),
	}
}
