# halp 

#### Interactive command line to import a csv of transaction into YNAB

Given a [csv](https://github.com/xild/halp/blob/master/README.md#how-to), import into YNAB what is not present, eg:

**bank** has transaction 1, transaction 2, and transaction 3

**ynab** has transaction 2 

**halp** will suggest the creation of the transaction 1 and 2, giving the opportunity to select a category. 

All transaction created will have a memo "sowhat? #cmdline" so it will be easier to remove or reject the transaction in case of fire.

# Run

Get ynab api token [here](https://api.youneedabudget.com/#getting-started)

`go get -u github.com/xild/halp` or clone.

`YNAB_TOKEN=$TOKEN go run main.go`

## How to


- n26 
  - [Export transactions n26](https://support.n26.com/en-eu/fixing-an-issue/payments-and-transfers/how-to-export-a-list-of-my-transactions)
 
 
 The csv should be delimited with a semilon `;` and must have the following structure, without a header:
 
 1. Column - The date in ISO format (e.g. 2016-12-01)
 
 2. Column -  The Payee name
 
 3. Column - The category (not used yet)
 
 4. Column - The transaction amount
 
 |  |  |   |  |
|---|---|---|---|
| 2019-08-19  | Japanese Market    | Food & Groceries  | -14.32
| 2019-08-19  | Helles Bier Market   | Fun  | -0.23
| 2019-08-19  | Movies   | Fun  | -8.20


