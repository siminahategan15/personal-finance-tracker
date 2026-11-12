# Personal Finance Tracker Architecture

## The shape

The app is a Go web service. A request comes in through one of its URLs, such as `/transactions`. The server in `main.go` receives it and calls the matching handler. The handler uses the finance code to update or read transactions, then sends a JSON response.

## The data

Each transaction has an ID, type, amount, category, description, and date. For this stage, the app keeps transactions in memory while it is running and saves them in `data/transactions.json`, so they can be loaded again after a restart. I plan to move the data to a database in a later stage. For now, the JSON file keeps this version simpler to build and test.

## Three decisions

1. I used Go's built-in web tools instead of adding a web framework. This keeps the project smaller and avoids an extra dependency.
2. I use a JSON file for this stage instead of a database. I plan to add a database later; the file is a simpler way to build and test the current version.
3. I separated the finance logic from the web server, but kept the project structure small. This makes the code easier to follow without adding many layers.

## What is not built yet

The database is not added yet. Other features that are not in this version include user accounts, budgets, recurring transactions, CSV export, and advanced reports. These can be added in later project stages.
