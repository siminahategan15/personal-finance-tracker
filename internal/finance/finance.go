package finance

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	TypeIncome  = "income"
	TypeExpense = "expense"
)

type Transaction struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type Summary struct {
	Month         string             `json:"month"`
	TotalIncome   float64            `json:"total_income"`
	TotalExpenses float64            `json:"total_expenses"`
	Net           float64            `json:"net"`
	Categories    map[string]float64 `json:"categories"`
}

type Store struct {
	mu           sync.Mutex
	transactions []Transaction
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Reset() {
	s.mu.Lock()
	s.transactions = nil
	s.mu.Unlock()
}

func (s *Store) ListTransactions() []Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Transaction, len(s.transactions))
	copy(out, s.transactions)
	return out
}

func (s *Store) GetTransaction(id string) (Transaction, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, tx := range s.transactions {
		if tx.ID == id {
			return tx, true
		}
	}
	return Transaction{}, false
}

func (s *Store) DeleteTransaction(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, tx := range s.transactions {
		if tx.ID == id {
			s.transactions = append(s.transactions[:i], s.transactions[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) AddTransaction(tx Transaction) error {
	if strings.TrimSpace(tx.Type) == "" {
		return errors.New("type is required")
	}
	if tx.Type != TypeIncome && tx.Type != TypeExpense {
		return fmt.Errorf("type must be %q or %q", TypeIncome, TypeExpense)
	}
	if tx.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if strings.TrimSpace(tx.Category) == "" {
		return errors.New("category is required")
	}
	if strings.TrimSpace(tx.Description) == "" {
		return errors.New("description is required")
	}
	if strings.TrimSpace(tx.Date) == "" {
		tx.Date = time.Now().Format("2006-01-02")
	}
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = time.Now().UTC()
	}
	if tx.ID == "" {
		tx.ID = fmt.Sprintf("tx-%d", time.Now().UnixNano())
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions = append(s.transactions, tx)
	return nil
}

func (s *Store) Summary(month string) Summary {
	s.mu.Lock()
	defer s.mu.Unlock()

	summary := Summary{Month: month, Categories: map[string]float64{}}
	for _, tx := range s.transactions {
		if month != "" && !strings.HasPrefix(tx.Date, month) {
			continue
		}
		switch tx.Type {
		case TypeIncome:
			summary.TotalIncome += tx.Amount
		case TypeExpense:
			summary.TotalExpenses += tx.Amount
		}
		summary.Categories[tx.Category] += tx.Amount
	}
	summary.Net = summary.TotalIncome - summary.TotalExpenses
	return summary
}

func (s *Store) SummaryWithCategory(month, category string) Summary {
	summary := s.Summary(month)
	if category == "" {
		return summary
	}
	filtered := Summary{Month: month, Categories: map[string]float64{}}
	for _, tx := range s.ListTransactions() {
		if month != "" && !strings.HasPrefix(tx.Date, month) {
			continue
		}
		if category != "" && !strings.EqualFold(tx.Category, category) {
			continue
		}
		switch tx.Type {
		case TypeIncome:
			filtered.TotalIncome += tx.Amount
		case TypeExpense:
			filtered.TotalExpenses += tx.Amount
		}
		filtered.Categories[tx.Category] += tx.Amount
	}
	filtered.Net = filtered.TotalIncome - filtered.TotalExpenses
	return filtered
}

func (s *Store) Categories() []string {
	seen := map[string]struct{}{}
	items := make([]string, 0)
	for _, tx := range s.ListTransactions() {
		if _, ok := seen[tx.Category]; ok {
			continue
		}
		seen[tx.Category] = struct{}{}
		items = append(items, tx.Category)
	}
	sort.Strings(items)
	return items
}
