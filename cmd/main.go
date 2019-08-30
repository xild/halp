package main

import (
	"fmt"
	"github.com/kataras/tablewriter"
	"github.com/manifoldco/promptui"
	"github.com/xild/youn26b/cmd/bank"
	"github.com/xild/youn26b/cmd/ynab"
	"os"
	"sort"
	"time"
)

type Executor struct {
	ynab *ynab.YNAB
	bank *bank.CSV
}

func main() {
	e := Executor{
		ynab: ynab.New(),
		bank: bank.New(),
	}

	budget, err := e.fetchBudget()
	if err != nil {
		panic(err)
	}

	acc, err := e.fetchAccount(budget.ID)
	if err != nil {
		panic(err)
	}

	categories, err := e.ynab.GetCategories(budget.ID)
	if err != nil {
		// it's ok to panic
		panic(err)
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	fmt.Println(firstOfMonth)
	p := promptui.Prompt{
		Label:    "Get transaction since date (YYYY-MM-DD)",
		Default:  firstOfMonth.Format("2006-01-02"),
		Validate: nil,
	}
	sinceDate, err := p.Run()
	if err != nil {
		return
	}

	ynabTXs, err := e.fetchYNABTransaction(budget, sinceDate)
	if err != nil {
		panic(err)
	}

	bankTXs, err := e.fetchBankTransaction(sinceDate)
	if err != nil {
		panic(err)
	}

	finalTXs := e.leftAntiJoin(e.bankTXtoYNABTX(bankTXs), ynabTXs.Data.Transactions)

	if len(finalTXs) == 0 {
		fmt.Println("Nothing to be created, seems you are up to date :)")
	} else {
		fmt.Println("#######################################################")
		fmt.Printf("#### You have %d to be created ####", len(finalTXs))
		fmt.Println("#######################################################")
		sort.Slice(finalTXs, func(i, j int) bool {
			return finalTXs[i].Date <= finalTXs[j].Date
		})

		for _, t := range finalTXs {
			t.AccountID = acc.ID
			e.suggestCreateTransaction(budget.ID, t, categories.Data.CategoryGroups)
		}
	}

}

func (e *Executor) suggestCreateTransaction(budgetID string, transaction ynab.Transaction, categoryGroups []ynab.CategoryGroups) {
	transaction.Memo = "cmdline #sowhat?"
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Payee", "Amount", "Date", "Memo"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor})
	table.Append([]string{transaction.PayeeName, fmt.Sprintf("%.2f", float64(transaction.Amount)*0.0001), transaction.Date, transaction.Memo})

	table.Render()

	wantCategoryPrompt := promptui.Prompt{
		Label: "Wanna specify one category? [y/N]",
	}

	result, err := wantCategoryPrompt.Run()
	if err != nil {
		panic(err)
		return
	}

	if result == "y" || result == "Y" {
		prompt := promptui.Select{
			Label: "Select a category group",
			Items: categoryGroups,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U0001F525 {{ .Name | bold }} ({{ .ID | red | italic }})",
				Inactive: "   {{ .Name | bold }} ({{ .ID | red | italic }})",
				Selected: "\U0001F525 Category group {{ .Name | red | bold }}",
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			return
		}

		categories := categoryGroups[i].Categories
		prompt = promptui.Select{
			Label: "Select a sub category",
			Items: categories,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U0001F525 {{ .Name | bold }} ({{ .ID | red | italic }})",
				Inactive: "   {{ .Name | bold }} ({{ .ID | red | italic }})",
				Selected: "\U0001F525 Category  {{ .Name | red | bold }}",
			},
		}

		i, _, err = prompt.Run()
		if err != nil {
			return
		}
		category := categories[i]

		transaction.CategoryName = category.Name
		transaction.CategoryID = category.ID
	}
	err = e.ynab.CreateTransaction(budgetID, transaction)
	if err != nil {
		panic(err)
	}

}

func (e *Executor) fetchBudget() (budget ynab.Budget, err error) {
	fmt.Println("1# Fetching budget")
	budgets, err := e.ynab.GetBudget()
	if err != nil {
		panic(err)
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F525 {{ .Name | bold }} ({{ .ID | red | italic }})",
		Inactive: "   {{ .Name | bold }} ({{ .ID | red | italic }})",
		Selected: "\U0001F525 {{ .Name | red | bold }}",
	}
	prompt := promptui.Select{
		Label:     "Select a Budget",
		Items:     budgets.Data.Budgets,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return
	}
	budget = budgets.Data.Budgets[i]
	return

}

func (e *Executor) fetchAccount(budgetID string) (account ynab.Account, err error) {
	fmt.Println("1# Fetching accounts")
	ac, err := e.ynab.GetAccounts(budgetID)
	if err != nil {
		panic(err)
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F525 {{ .Name | bold }} ({{ .ID | red | italic }})",
		Inactive: "   {{ .Name | bold }} ({{ .ID | red | italic }})",
		Selected: "\U0001F525 {{ .Name | red | bold }}",
	}
	prompt := promptui.Select{
		Label:     "Select an Account",
		Items:     ac.Data.Accounts,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return
	}

	return ac.Data.Accounts[i], nil

}

func (e *Executor) fetchYNABTransaction(budget ynab.Budget, sinceDate string) (ynab.Transactions, error) {
	return e.ynab.GetTransaction(budget.ID, sinceDate)
}

func (e *Executor) fetchBankTransaction(sinceDate string) (txs []bank.Transaction, err error) {
	fmt.Println("#2 Fetching bank transaction")
	p := promptui.Prompt{
		Label: "File transaction path",
	}
	fileName, err := p.Run()
	if err != nil {
		return
	}
	return e.bank.Fetch(fileName, sinceDate)
}

func (e *Executor) bankTXtoYNABTX(transactions []bank.Transaction) (txs []ynab.Transaction) {
	for _, btx := range transactions {
		amount := btx.Price

		txs = append(txs, ynab.Transaction{
			Date:         btx.Date,
			Amount:       amount,
			Deleted:      false,
			PayeeName:    btx.Payee,
			CategoryName: "",
		})
	}
	return
}

func (e *Executor) leftAntiJoin(left []ynab.Transaction, right []ynab.Transaction) (result []ynab.Transaction) {
	// maybe sort and search
	for _, l := range left {
		if !contains(right, l) {
			result = append(result, l)
		}
	}
	return
}

func contains(a []ynab.Transaction, b ynab.Transaction) bool {
	for _, n := range a {
		if b.Date == n.Date && b.Amount == n.Amount {
			return true
		}
	}
	return false
}
