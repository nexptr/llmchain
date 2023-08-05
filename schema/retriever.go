package schema

import "context"

// Retriever is an interface that defines the behavior of a retriever.
type Retriever interface {
	//GetRelevantDocuments Get documents relevant for a query.

	//     Args:
	//         query: string to find relevant documents for

	//     Returns:
	//         List of relevant documents
	GetRelevantDocuments(ctx context.Context, query string) ([]Document, error)
}
