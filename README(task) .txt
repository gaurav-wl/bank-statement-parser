Welcome to the Test Pack.

Your test pack will include:
* This readme file
* The bank statement file (statement.csv)


Objective:
We want to build a microservice that parses the bank statement to deduce the total value of payments made
on 6 March 2011 for each currency. All lines of the bank statement that represent a payment will contain a payment
reference in one of it's narratives. A payment reference is guaranteed to be the letters "PAY" then 6 digits and then
two more letters (e.g. PAY000001YB). The currencies that we deal with for the purpose of this exercise are GBP, EUR,
USD and CAD. If there is another currency in a bank statement then it should be in the output as well.


The microservice should be built using Golang.

This microservice need to produce as an output the account balances per each currency.
It should be possible for us to run it on our local machines with a few commands.

Sample Output:
Totals
GBP 3,432.21
USD 14,432.34
EUR 43,543.23
CAD 2,321.34

Once the code is ready, just give us access to the repo (this is our Git: cursol)

