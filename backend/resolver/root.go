package resolver

import (
	"goask/core/entity"

	"github.com/graph-gophers/graphql-go"
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

func (q Query) Action(args struct{ UserID graphql.ID }) QueryAction {
	return QueryAction{
		stdResolver: q.stdResolver,
		userSession: UserSession{UserID: entity.NewIDString(string(args.UserID))},
	}
}

type Mutation struct {
	stdResolver
}

func NewMutation(stdResolver stdResolver) Mutation {
	return Mutation{stdResolver}
}
