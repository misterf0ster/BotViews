package report

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Send(controlURL string, orderID, cycles int, botID string) {
	payload := fmt.Sprintf(`{"order_id":%d,"cycles":%d,"bot_id":"%s"}`, orderID, cycles, botID)
	resp, err := http.Post(controlURL, "application/json", strings.NewReader(payload))
	if err != nil {
		log.Printf("[%s] report POST error: %v", botID, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("[%s] Reported to controller: order_id=%d cycles=%d status=%s", botID, orderID, cycles, resp.Status)
}
