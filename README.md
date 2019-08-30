# halp 

#### Interactive command line to import a csv of transaction into YNAB

Given a [csv](##Howto) import into YNAB what is not present there, eg:

**bank** has transaction 1, transaction 2, and transaction 3

**ynab** has transaction 2 

**halp** will suggest the creation of the transaction 1 and 2, giving the opportunity to select a category. 



## How to


- n26 
  - [Export transactions n26](https://support.n26.com/en-eu/fixing-an-issue/payments-and-transfers/how-to-export-a-list-of-my-transactions)
 
 
 The csv should be delimited with a semilon `;` and must have the following structure, without a header.
 
 First column - The date in ISO format (e.g. 2016-12-01)
 Second column -  The Payee name
 Third column - The category (not used yet)
 Fourth column - The transaction amount
 
 |  |  |   |  |
|---|---|---|---|
| 2019-08-19  | Japanese Market    | Food & Groceries  | -14.32
| 2019-08-19  | Helles Bier Market   | Fun  | -0.23
| 2019-08-19  | Movies   | Fun  | -8.20

