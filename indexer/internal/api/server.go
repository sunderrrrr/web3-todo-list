package api

import (
	"encoding/json"
	"net/http"

	"w3todo-indexer/internal/db"
)

type Server struct {
	db *db.Postgres
}

func New(database *db.Postgres) *Server {
	return &Server{db: database}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/transfers", s.handleTransfers)
	mux.HandleFunc("/api/todos", s.handleTodos)
	mux.HandleFunc("/api/balances", s.handleBalances)
	mux.HandleFunc("/api/rewards", s.handleRewards)
	return withCORS(mux)
}

func (s *Server) handleTransfers(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Conn.Query(`SELECT from_addr, to_addr, value, tx_hash FROM transfers ORDER BY id DESC LIMIT 20`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var result []map[string]any
	for rows.Next() {
		var from, to, value, txHash string
		rows.Scan(&from, &to, &value, &txHash)
		result = append(result, map[string]any{
			"from": from, "to": to, "value": value, "tx_hash": txHash,
		})
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Conn.Query(`SELECT id, text, owner FROM todos ORDER BY id DESC LIMIT 20`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var result []map[string]any
	for rows.Next() {
		var id uint64
		var text, owner string
		rows.Scan(&id, &text, &owner)
		result = append(result, map[string]any{
			"id": id, "text": text, "owner": owner,
		})
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleBalances(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Conn.Query(`SELECT addr, balance FROM balances WHERE balance > 0 ORDER BY balance DESC`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var result []map[string]any
	for rows.Next() {
		var addr, balance string
		rows.Scan(&addr, &balance)
		result = append(result, map[string]any{
			"addr": addr, "balance": balance,
		})
	}
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleRewards(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Conn.Query(`SELECT user_addr, amount, tx_hash FROM rewards ORDER BY id DESC LIMIT 20`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var result []map[string]any
	for rows.Next() {
		var addr, amount, txHash string
		rows.Scan(&addr, &amount, &txHash)
		result = append(result, map[string]any{
			"user_addr": addr, "amount": amount, "tx_hash": txHash,
		})
	}
	json.NewEncoder(w).Encode(result)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}
