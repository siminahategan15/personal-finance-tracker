package finance

import "testing"

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Fatal("expected store to be initialized")
	}
}

func TestAddTransaction(t *testing.T) {
	store := NewStore()
	transaction := Transaction{ID: "t1", Type: "expense", Amount: 25.50, Category: "food", Description: "Lunch", Date: "2026-11-07"}

	if err := store.AddTransaction(transaction); err != nil {
		t.Fatalf("AddTransaction returned error: %v", err)
	}

	list := store.ListTransactions()
	if len(list) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(list))
	}
}

func TestSummary(t *testing.T) {
	store := NewStore()
	transactions := []Transaction{
		{ID: "t1", Type: "income", Amount: 2000, Category: "salary", Description: "Paycheck", Date: "2026-11-01"},
		{ID: "t2", Type: "expense", Amount: 50, Category: "food", Description: "Groceries", Date: "2026-11-03"},
	}
	for _, tx := range transactions {
		if err := store.AddTransaction(tx); err != nil {
			t.Fatalf("AddTransaction: %v", err)
		}
	}

	summary := store.Summary("2026-11")
	if summary.TotalIncome != 2000 {
		t.Fatalf("expected income 2000, got %v", summary.TotalIncome)
	}
	if summary.TotalExpenses != 50 {
		t.Fatalf("expected expense 50, got %v", summary.TotalExpenses)
	}
	if summary.Net != 1950 {
		t.Fatalf("expected net 1950, got %v", summary.Net)
	}
}
