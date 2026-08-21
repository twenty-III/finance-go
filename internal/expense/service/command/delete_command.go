package expensecmd

import "github.com/google/uuid"

// DeleteCommand represents a command to delete an existing expense.
type DeleteCommand struct {
	Id     uuid.UUID // Unique identifier of the expense to be deleted
	UserId uuid.UUID // Identifier of the user who owns the expense
}
