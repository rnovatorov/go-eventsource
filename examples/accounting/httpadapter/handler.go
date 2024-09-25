package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rnovatorov/go-eventsource/examples/accounting/model"
)

type accountingService interface {
	CreateBook(
		ctx context.Context, bookID string, bookDescription string,
	) (string, error)
	CloseBook(
		ctx context.Context, bookID string,
	) error
	AddBookAccount(
		ctx context.Context, bookID string, accountName string,
		accountType model.AccountType,
	) error
	GetBookAccountBalance(
		ctx context.Context, bookID string, accountName string,
	) (uint64, error)
	EnterBookTransaction(
		ctx context.Context, bookID string, timestamp time.Time,
		accountDebited string, accountCredited string, amount uint64,
	) error
}

type Handler struct {
	mux               *http.ServeMux
	accountingService accountingService
}

func NewHandler(s accountingService) *Handler {
	h := &Handler{
		mux:               http.NewServeMux(),
		accountingService: s,
	}

	h.mux.HandleFunc("/book/create", h.handleBookCreate)
	h.mux.HandleFunc("/book/close", h.handleBookClose)
	h.mux.HandleFunc("/book/account/add", h.handleBookAccountAdd)
	h.mux.HandleFunc("/book/account/balance", h.handleBookAccountBalance)
	h.mux.HandleFunc("/book/transaction/enter", h.handleBookTransactionEnter)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) handleBookCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var payload struct {
		BookID          string `json:"book_id"`
		BookDescription string `json:"book_description"`
	}
	if err := h.unmarshalJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := h.accountingService.CreateBook(
		r.Context(), payload.BookID, payload.BookDescription,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleBookClose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var payload struct {
		BookID string `json:"book_id"`
	}
	if err := h.unmarshalJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.accountingService.CloseBook(
		r.Context(), payload.BookID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleBookAccountAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var payload struct {
		BookID      string `json:"book_id"`
		AccountName string `json:"account_name"`
		AccountType string `json:"account_type"`
	}
	if err := h.unmarshalJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.accountingService.AddBookAccount(
		r.Context(), payload.BookID, payload.AccountName,
		model.NewAccountType(payload.AccountType),
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleBookAccountBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, err := h.accountingService.GetBookAccountBalance(
		r.Context(), q.Get("book_id"), q.Get("account_name"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Balance uint64 `json:"balance"`
	}
	data, err := json.Marshal(response{
		Balance: balance,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) handleBookTransactionEnter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var payload struct {
		BookID          string `json:"book_id"`
		Timestamp       string `json:"timestamp"`
		AccountDebited  string `json:"account_debited"`
		AccountCredited string `json:"account_credited"`
		Amount          uint64 `json:"amount"`
	}
	if err := h.unmarshalJSON(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timestamp, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.accountingService.EnterBookTransaction(
		r.Context(), payload.BookID, timestamp,
		payload.AccountDebited, payload.AccountCredited, payload.Amount,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) unmarshalJSON(r *http.Request, dest any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	return json.Unmarshal(body, dest)
}
