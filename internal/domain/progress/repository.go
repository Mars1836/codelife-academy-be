package progress

import "context"

type Repository interface {
	ListByUser(context.Context, string) ([]LearningProgress, error)
	FindByUserAndDocument(context.Context, string, string) (LearningProgress, error)
	Upsert(context.Context, string, LearningProgress) (LearningProgress, error)
}
