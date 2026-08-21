package expensecmd

import (
	icmd "github.com/mohit/finance-go/internal/common/cqrs/command"
	irepository "github.com/mohit/finance-go/internal/common/repository"
)

// DeleteHandler processes a DeleteCommand to remove an expense.
type DeleteHandler struct {
	expenseRepo irepository.IExpenseRepository
}

// Ensure DeleteHandler implements IHandler
var _ icmd.IHandler[*DeleteCommand, any] = &DeleteHandler{}

// NewDeleteHandler creates a new DeleteHandler.
func NewDeleteHandler(expenseRepo irepository.IExpenseRepository) *DeleteHandler {
	return &DeleteHandler{
		expenseRepo: expenseRepo,
	}
}

// Handle executes the deletion of an expense.
func (h *DeleteHandler) Handle(cmd *DeleteCommand) (any, error) {
	err := h.expenseRepo.Delete(cmd.Id, cmd.UserId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
