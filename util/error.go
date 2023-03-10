package util

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/99designs/gqlgen/graphql"
)

func GqlError(ctx context.Context, message string, code string) *gqlerror.Error {
	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: message,
		Extensions: map[string]interface{}{
			"code": code,
		}}
}

/*
**The GraphQL operation string contains a syntax errors.
**/
const GraphqlParseFailed string = ("GRAPHQL_PARSE_FAILED")

/*
**The GraphQL operation is not valid against the server's schema.
**/
const GraphqlValidationFailed string = ("GRAPHQL_VALIDATION_FAILED")

/*
**The GraphQL operation includes an invalid value for a field argument.
**/
const BadUserInput string = ("BAD_USER_INPUT")

/*
**A client sent the hash of a query string to execute via automatic persisted queries, but the query was not in the APQ cache.
**/
const PersistedQueryNotFound string = ("PERSISTED_QUERY_NOT_FOUND")

/*
**A client sent the hash of a query string to execute via automatic persisted queries, but the server has disabled APQ.
**/
const PersistedQueryNotSupported string = ("PERSISTED_QUERY_NOT_SUPPORTED")

/*
**The request was parsed successfully and is valid against the server's schema, but the server couldn't resolve which operation to run.
gqlerror
**This occurs when a request containing multiple named operations doesn't specify which operation to run (i.e.,operationName), or if the named operation isn't included in the request.
*
*/
const OperationResolutionFailure string = ("OPERATION_RESOLUTION_FAILURE")

/*
**An errors occurred before your server could attempt to parse the given GraphQL operation.
**/
const BadRequest string = ("BAD_REQUEST")

/*
**An unspecified errors occurred.
*
**When Apollo Server formats an errors in a response, it sets the code extension to this value if no other code is set.
**/
const InternalServerError string = ("INTERNAL_SERVER_ERROR")
