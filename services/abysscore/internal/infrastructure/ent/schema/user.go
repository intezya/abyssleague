package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/schema/access_level"
	"time"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),

		field.String("username").NotEmpty().Unique(),
		field.String("lower_username").NotEmpty().Unique(),
		field.String("email").Nillable().Optional().Unique(),
		field.String("password").NotEmpty().Sensitive(),
		field.String("hardware_id").Nillable().Optional().Unique().Sensitive(),

		field.String("access_level").
			GoType(access_level.AccessLevel(0)).
			DefaultFunc(
				func() access_level.AccessLevel {
					return access_level.User
				},
			),

		field.String("genshin_uid").Optional().Nillable().Unique(),
		field.String("hoyolab_login").Optional().Nillable().Unique(),

		field.Int("current_match_id").Optional().Nillable(),
		field.Int("current_item_in_profile_id").Optional().Nillable().Unique(),

		field.String("avatar_url").Optional().Nillable(),

		field.Bool("invites_enabled").Default(false),

		field.Time("login_at").Default(
			func() time.Time {
				return time.Now()
			},
		),
		field.Int("login_streak").Default(0),

		field.Time("created_at").Default(time.Now).Immutable(),

		field.Time("search_blocked_until").Optional().Nillable(),
		field.String("search_block_reason").Optional().Nillable(),
		field.Int("search_blocked_level").Default(0).Min(0),

		field.Time("account_blocked_until").Optional().Nillable(),
		field.String("account_block_reason").Optional().Nillable(),
		field.Int("account_blocked_level").Default(0).Min(0),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("statistics", Statistic.Type),

		edge.To("friends", User.Type),

		edge.To("sent_friend_requests", FriendRequest.Type),
		edge.To("received_friend_requests", FriendRequest.Type),

		edge.To("items", InventoryItem.Type),

		edge.To("current_item", InventoryItem.Type).Field("current_item_in_profile_id").Unique(),

		edge.To("current_match", Match.Type).
			Unique().
			Field("current_match_id"),

		edge.To("balance", UserBalance.Type).
			Unique(),
	}
}
