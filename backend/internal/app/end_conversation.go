package app

import (
	"log/slog"

	"backend/internal/review"
)

func (a *conversationApp) EndConversation(req *EndConversationRequest) error {
	if req.ReviewID == "" {
		return nil
	}
	// set the review status to closed
	rev, err := a.reviewRepo.Get(req.ReviewID)
	if err != nil {
		a.logger.Debug("failed to get review", slog.Any("error", err))
		return err
	}
	rev.Status = review.StatusClosed
	_, err = a.reviewRepo.Upsert(rev)
	if err != nil {
		a.logger.Debug("failed to update review", slog.Any("error", err))
		return err
	}
	return nil
}
