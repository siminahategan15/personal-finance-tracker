package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"atad-project/internal/finance"
	"atad-project/internal/version"
)

var started = time.Now()

type appState struct {
	mu    sync.Mutex
	store *finance.Store
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	for _, a := range os.Args[1:] {
		if a == "--version" || a == "-v" {
			return version.WriteJSON(os.Stdout)
		}
	}

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	store := finance.NewStore("data/transactions.json")
	state := &appState{store: store}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = version.WriteJSON(w)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":         "ok",
			"uptime_seconds": int(time.Since(started).Seconds()),
		})
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!doctype html><html lang="en"><meta charset="utf-8">
<title>%s</title><body>
<h1>%s</h1>
<p>Personal Finance Tracker</p>
<p>Track income, expenses, and monthly balances.</p>
</body></html>`, version.Current().App, version.Current().App)
	})

	mux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		state.mu.Lock()
		state.store.Reset()
		state.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("GET /transactions", func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		items := state.store.ListTransactions()
		if category != "" {
			filtered := make([]finance.Transaction, 0)
			for _, tx := range items {
				if strings.EqualFold(tx.Category, category) {
					filtered = append(filtered, tx)
				}
			}
			items = filtered
		}
		writeJSON(w, http.StatusOK, items)
	})

	mux.HandleFunc("POST /transactions", func(w http.ResponseWriter, r *http.Request) {
		var tx finance.Transaction
		if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
			return
		}
		if err := state.store.AddTransaction(tx); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_transaction", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, tx)
	})

	mux.HandleFunc("GET /transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/transactions/")
		tx, ok := state.store.GetTransaction(id)
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "transaction not found")
			return
		}
		writeJSON(w, http.StatusOK, tx)
	})

	mux.HandleFunc("DELETE /transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/transactions/")
		if !state.store.DeleteTransaction(id) {
			writeError(w, http.StatusNotFound, "not_found", "transaction not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("GET /summary", func(w http.ResponseWriter, r *http.Request) {
		month := r.URL.Query().Get("month")
		category := r.URL.Query().Get("category")
		summary := state.store.Summary(month)
		if category != "" {
			summary = state.store.SummaryWithCategory(month, category)
		}
		writeJSON(w, http.StatusOK, summary)
	})

	mux.HandleFunc("GET /categories", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, state.store.Categories())
	})

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", srv.Addr, "commit", version.Current().Commit)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"error": code, "message": message})
}
