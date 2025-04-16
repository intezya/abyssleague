package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type FriendRequest struct {
	ent.Schema
}

func (FriendRequest) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.Int("from_user_id").Immutable(),
		field.Int("to_user_id").Immutable(),

		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (FriendRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("from_user", User.Type).
			Ref("sent_friend_requests").
			Unique().
			Required().
			Immutable().
			Field("from_user_id"),

		edge.From("to_user", User.Type).
			Ref("received_friend_requests").
			Unique().
			Required().
			Immutable().
			Field("to_user_id"),
	}
}
