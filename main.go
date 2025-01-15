package main

import (
	"encoding/json"
	"net/http"
	"strings"

	_ "email_spammer/docs" // docs is generated by Swag CLI, you have to import it.

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	"go.uber.org/zap"
)

// EmailRequest представляет собой структуру запроса
type EmailRequest struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Body      string `json:"body"`
}

// ScanResponse представляет собой структуру ответа
type ScanResponse struct {
	IsMalicious bool   `json:"is_malicious"`
	Reason      string `json:"reason"`
}

// @Summary Scan email for viruses
// @Description Scan the email message for malicious content
// @Accept json
// @Produce json
// @Param email body EmailRequest true "Email to scan"
// @Success 200 {object} ScanResponse
// @Router /scan [post]
func ScanEmail(w http.ResponseWriter, r *http.Request, logger *zap.Logger) {
	var email EmailRequest

	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Для простоты примера, мы просто проверим на наличие "malware" в тексте.
	isMalicious := false
	reason := ""

	if containsMalware(email.Body) {
		isMalicious = true
		reason = "Email body contains known malware keywords."
	}

	response := ScanResponse{
		IsMalicious: isMalicious,
		Reason:      reason,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func containsMalware(body string) bool {
	malwareKeywords := []string{"malware", "virus", "trojan"}
	for _, keyword := range malwareKeywords {
		if contains(body, keyword) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	router := mux.NewRouter()

	router.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		ScanEmail(w, r, logger) // Передаем логгер в обработчик
	}).Methods("POST")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
