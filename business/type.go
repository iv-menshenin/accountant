package business

import "github.com/iv-menshenin/accountant/store"

type (
	App struct {
		accounts *store.AccountCollection
	}
)

func New(
	accounts *store.AccountCollection,
) *App {
	return &App{
		accounts: accounts,
	}
}
