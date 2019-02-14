package resolver

import (
	"goask/core/entity"
)

type Root struct {
	Query
	Mutation
}

type Query struct {
	stdResolver
}

func NewQuery(stdResolver stdResolver) Query {
	return Query{stdResolver: stdResolver}
}

func (q Query) Action(args struct{ UserID int32 }) QueryAction {
	return QueryAction{
		stdResolver: q.stdResolver,
		userSession: UserSession{UserID: entity.ID(args.UserID)},
	}
}

type Mutation struct {
	stdResolver
}

func NewMutation(stdResolver stdResolver) Mutation {
	return Mutation{stdResolver}
}
