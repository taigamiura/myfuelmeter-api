package future

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/taigamiura/my-fuel-meter-project/api/testhelper"
	"github.com/taigamiura/my-fuel-meter-project/api/utils"
)

func TestHandleWebSocket(t *testing.T) {
	mockRedis := testhelper.NewMockRedis(t)
	mockDB, mock, err := testhelper.GetNewDbMock()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(w, r, mockDB, mockRedis)
	}))
	defer server.Close()

	// WebSocketに接続
	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// テストメッセージを送信
	testMessage := utils.Coordinates{
		Latitude:  35.6895,
		Longitude: 139.6917,
		Timestamp: "2023-10-10T10:00:00Z",
		Message:   "START_TRACKING",
	}
	messageBytes, err := json.Marshal(testMessage)
	if err != nil {
		t.Fatalf("Failed to marshal test message: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// レスポンスを読み取る
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	expectedSubstring := `"status":"SUCCESS"`
	if !strings.Contains(string(message), expectedSubstring) {
		t.Errorf("Expected response to contain %s, but got %s", expectedSubstring, message)
	}

	// WebSocket接続を正常に閉じる
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		t.Fatalf("Failed to close WebSocket connection: %v", err)
	}

	// すべての期待値が満たされたことを確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}
