package graph

import (
	icmd "github.com/mohit/finance-go/internal/common/cqrs/command"
	iquery "github.com/mohit/finance-go/internal/common/cqrs/query"
	expensemodel "github.com/mohit/finance-go/internal/expense/model"
	expensecmd "github.com/mohit/finance-go/internal/expense/service/command"
	expensqry "github.com/mohit/finance-go/internal/expense/service/query"
)

type Resolver struct {
	getExpenseHandler         iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
	getMultipleExpenseHandler iquery.IHandler[*expensqry.GetMultipleQuery, []*expensemodel.Expense]
	addExpenseHandler         icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	patchExpenseHandler       icmd.IHandler[*expensecmd.PatchCommand, *expensemodel.Expense]
}

type ResolverConfig struct {
	GetExpenseHandler         iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
	GetMultipleExpenseHandler iquery.IHandler[*expensqry.GetMultipleQuery, []*expensemodel.Expense]
	AddExpenseHandler         icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	PatchExpenseHandler       icmd.IHandler[*expensecmd.PatchCommand, *expensemodel.Expense]
}

func NewResolver(c ResolverConfig) *Resolver {
	return &Resolver{
		getExpenseHandler:         c.GetExpenseHandler,
		getMultipleExpenseHandler: c.GetMultipleExpenseHandler,
		addExpenseHandler:         c.AddExpenseHandler,
		patchExpenseHandler:       c.PatchExpenseHandler,
	}

}
