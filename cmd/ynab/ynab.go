package ynab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"src/github.com/go-resty/resty"
	"time"
)

// transactions structs
type (
	SaveTransactionWrapper struct {
		SaveTransactionData
	}

	SaveTransactionData struct {
		Transactions []Transaction `json:"transactions,omitempty"` // only one should be provider
		Transaction  Transaction   `json:"transaction,omitempty"`  // only one should be provider
	}
	Transactions struct {
		Data DataTransaction `json:"data"`
	}

	Subtransactions struct {
		ID                string `json:"id"`
		TransactionID     string `json:"transaction_id"`
		Amount            int    `json:"amount"`
		Memo              string `json:"memo"`
		PayeeID           string `json:"payee_id"`
		CategoryID        string `json:"category_id"`
		TransferAccountID string `json:"transfer_account_id"`
		Deleted           bool   `json:"deleted"`
	}

	Transaction struct {
		ID                    string            `json:"id,omitempty"`
		Date                  string            `json:"date,omitempty"`
		Amount                int               `json:"amount,omitempty"`
		Memo                  string            `json:"memo,omitempty"`
		Cleared               string            `json:"cleared,omitempty"`
		Approved              bool              `json:"approved,omitempty"`
		FlagColor             string            `json:"flag_color,omitempty"`
		AccountID             string            `json:"account_id,omitempty"`
		PayeeID               string            `json:"payee_id,omitempty"`
		CategoryID            string            `json:"category_id,omitempty"`
		TransferAccountID     string            `json:"transfer_account_id,omitempty"`
		TransferTransactionID string            `json:"transfer_transaction_id,omitempty"`
		MatchedTransactionID  string            `json:"matched_transaction_id,omitempty"`
		ImportID              string            `json:"import_id,omitempty"`
		Deleted               bool              `json:"deleted,omitempty"`
		AccountName           string            `json:"account_name,omitempty"`
		PayeeName             string            `json:"payee_name,omitempty"`
		CategoryName          string            `json:"category_name,omitempty"`
		Subtransactions       []Subtransactions `json:"subtransactions,omitempty"`
	}
	DataTransaction struct {
		Transactions    []Transaction `json:"transactions"`
		ServerKnowledge int           `json:"server_knowledge"`
	}
)

// accounts
type (
	Accounts struct {
		Data struct {
			Accounts        []Account `json:"accounts"`
			ServerKnowledge int       `json:"server_knowledge"`
		} `json:"data"`
	}
	Account struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		Type             string `json:"type"`
		OnBudget         bool   `json:"on_budget"`
		Closed           bool   `json:"closed"`
		Note             string `json:"note"`
		Balance          int    `json:"balance"`
		ClearedBalance   int    `json:"cleared_balance"`
		UnclearedBalance int    `json:"uncleared_balance"`
		TransferPayeeID  string `json:"transfer_payee_id"`
		Deleted          bool   `json:"deleted"`
	}
)

// budgets structs
type (
	Budgets struct {
		Data Data `json:"data"`
	}

	Data struct {
		Budgets []Budget `json:"budgets"`
	}

	Budget struct {
		ID             string         `json:"id"`
		Name           string         `json:"name"`
		LastModifiedOn time.Time      `json:"last_modified_on"`
		FirstMonth     string         `json:"first_month"`
		LastMonth      string         `json:"last_month"`
		DateFormat     DateFormat     `json:"date_format"`
		CurrencyFormat CurrencyFormat `json:"currency_format"`
	}

	DateFormat struct {
		Format string `json:"format"`
	}
	CurrencyFormat struct {
		IsoCode          string `json:"iso_code"`
		ExampleFormat    string `json:"example_format"`
		DecimalDigits    int    `json:"decimal_digits"`
		DecimalSeparator string `json:"decimal_separator"`
		SymbolFirst      bool   `json:"symbol_first"`
		GroupSeparator   string `json:"group_separator"`
		CurrencySymbol   string `json:"currency_symbol"`
		DisplaySymbol    bool   `json:"display_symbol"`
	}
)

// Categories
type (
	Categories struct {
		Data CategoryData `json:"data"`
	}

	Category struct {
		ID                      string `json:"id"`
		CategoryGroupID         string `json:"category_group_id"`
		Name                    string `json:"name"`
		Hidden                  bool   `json:"hidden"`
		OriginalCategoryGroupID string `json:"original_category_group_id"`
		Note                    string `json:"note"`
		Budgeted                int    `json:"budgeted"`
		Activity                int    `json:"activity"`
		Balance                 int    `json:"balance"`
		GoalType                string `json:"goal_type"`
		GoalCreationMonth       string `json:"goal_creation_month"`
		GoalTarget              int    `json:"goal_target"`
		GoalTargetMonth         string `json:"goal_target_month"`
		GoalPercentageComplete  int    `json:"goal_percentage_complete"`
		Deleted                 bool   `json:"deleted"`
	}
	CategoryGroups struct {
		ID         string     `json:"id"`
		Name       string     `json:"name"`
		Hidden     bool       `json:"hidden"`
		Deleted    bool       `json:"deleted"`
		Categories []Category `json:"categories"`
	}

	CategoryData struct {
		CategoryGroups  []CategoryGroups `json:"category_groups"`
		ServerKnowledge int              `json:"server_knowledge"`
	}
)

type API interface {
	GetBudget() Budgets
	GetTransaction(budgetID, sinceDate string) (txs Transactions, err error)
	GetCategories(budgetID string) (cd CategoryData, err error)
	CreateTransaction(budgetID string, tx Transaction) (err error)
	GetAccounts(budgetID string, tx Transaction) (err error)
}

type YNAB struct {
	token string
	req   *resty.Request
}

// https://api.youneedabudget.com/v1#/Transactions/createTransaction
func (y YNAB) CreateTransaction(budgetID string, tx Transaction) (err error) {
	t := SaveTransactionWrapper{
		SaveTransactionData: SaveTransactionData{
			Transaction: tx,
		},
	}
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	response, err := y.post("/budgets/"+budgetID+"/transactions", b)
	if err != nil {
		return
	}

	if !response.IsSuccess() {
		return fmt.Errorf("HTTP status `code >= 200 and <= 299`, actual=%v, body=%v", response.StatusCode(), response.String())
	}
	return
}

// Get Categories return a list of budget under the YNAB account
// https://api.youneedabudget.com/v1#/Categories/getCategories
func (y YNAB) GetCategories(budgetID string) (categories Categories, err error) {
	response, err := y.get("/budgets/" + budgetID + "/categories")
	if err != nil {
		return
	}

	if !response.IsSuccess() {
		return categories, fmt.Errorf("HTTP status `code >= 200 and <= 299`, actual %v=", response.StatusCode())
	}

	err = json.Unmarshal(response.Body(), &categories)
	if err != nil {
		return
	}
	return
}

// Get Accounts return the accounts linked with a budget
// https://api.youneedabudget.com/v1#/Accounts
func (y YNAB) GetAccounts(budgetID string) (account Accounts, err error) {
	response, err := y.get("/budgets/" + budgetID + "/accounts")
	if err != nil {
		return
	}

	if !response.IsSuccess() {
		return account, fmt.Errorf("HTTP status `code >= 200 and <= 299`, actual %v=", response.StatusCode())
	}

	err = json.Unmarshal(response.Body(), &account)
	if err != nil {
		return
	}
	return
}

// Get Budget return a list of budget under the YNAB account
// https://api.youneedabudget.com/v1#/Budgets/getBudgets
func (y YNAB) GetBudget() (budgets Budgets, err error) {
	response, err := y.get("/budgets")
	if err != nil {
		return
	}

	if !response.IsSuccess() {
		return budgets, fmt.Errorf("HTTP status `code >= 200 and <= 299`, actual %v=", response.StatusCode())
	}

	err = json.Unmarshal(response.Body(), &budgets)
	if err != nil {
		return
	}
	return
}

// Given a budgetID return all transactions
// https://api.youneedabudget.com/v1#/Transactions/getTransactions
// TODO filter by date
func (y YNAB) GetTransaction(budgetID, sinceDate string) (txs Transactions, err error) {
	resp, err := y.get("/budgets/" + budgetID + "/transactions")
	if err != nil || resp.StatusCode() != http.StatusOK {
		panic("error")
	}
	err = json.Unmarshal(resp.Body(), &txs)
	if err != nil {
		return
	}
	return
}

func New() *YNAB {
	return &YNAB{
		req: resty.New().R().SetAuthToken(os.Getenv("YNAB_TOKEN")),
	}
}

func (y YNAB) get(path string) (*resty.Response, error) {
	return y.req.EnableTrace().Get("https://api.youneedabudget.com/v1/" + path)
}

func (y YNAB) post(path string, body []byte) (*resty.Response, error) {
	fmt.Println(string(body))
	return y.req.EnableTrace().SetHeader("Content-Type", "application/json").SetBody(body).Post("https://api.youneedabudget.com/v1/" + path)
}
