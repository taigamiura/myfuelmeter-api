package future

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"github.com/taigamiura/my-fuel-meter-project/api/services"
	"github.com/taigamiura/my-fuel-meter-project/api/utils"
)

var validate *validator.Validate
var db *gorm.DB // GORMのDB接続
// var Rdb *redis.Client // Redisクライアント
// WebSocketのアップグレーダー
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 必要に応じてセキュリティ設定を行う
	},
}

func init() {
	validate = validator.New()
}

// WebSocket接続を処理するハンドラー
func HandleWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB, rdb *redis.Client) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Println("Client connected")
	handleClientConnection(db, rdb, conn)
}

// クライアント接続を処理する関数
func handleClientConnection(db *gorm.DB, rdb *redis.Client, conn *websocket.Conn) {
	var startCoords utils.Coordinates
	var endCoords utils.Coordinates
	isFirstMessage := true
	isSendFinishResponse := false
	// 連続エラー回数：5以上になった場合は処理を終了する
	errorRepeatCount := 0
	// main処理
	// アプリケーションから一定時間毎にtimestamp, 緯度、経度情報がリクエストされる
	// リクエストは基本的にredisに保存され、処理の最後にDBに保存する。
	// DB保存後は移動距離・ガソリン代・移動時間を計算してレスポンスする
	// アプリケーションが接続を終了するまで処理を続行する
	for {
		if errorRepeatCount >= 5 {
			break
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("Client disconnected gracefully:", err)
				break
			}
			handleConnectionError(err)
			errorRepeatCount++
			log.Println("errorRepeatCount:", errorRepeatCount)
			continue
		}

		log.Println("Received message from client:", string(msg))
		coords, err := decodeCoordinates(msg)
		if err != nil {
			sendValidationError(conn)
			errorRepeatCount++
			log.Println("errorRepeatCount:", errorRepeatCount)
			continue
		}

		// 現在の座標の更新
		endCoords = coords
		processingStartTime, traceID := handleStartMessage(&startCoords, &isFirstMessage, coords)

		if err := services.SaveCoordinatesToRedis(rdb, &coords, traceID); err != nil {
			log.Println("Error saving coordinates to Redis:", err)
			errorRepeatCount++
			log.Println("errorRepeatCount:", errorRepeatCount)
			continue
		}

		if coords.Message == "FINISH_TRACKING" && !isSendFinishResponse {
			coordinates, err := services.GetCoordinatesFromRedis(rdb, coords.TraceID)
			if err != nil {
				log.Println("Error get coordinates to Resid: ", err)
			}
			// 隣接する値の合計を計算して出力
			var distance float64 = 0
			for i := 0; i < len(coordinates)-1; i++ {
				distance += utils.HaversineDistance(coordinates[i], coordinates[i+1])
			}
			fuelCost := utils.CalculateFuelCost(distance) // 燃料費の計算ロジックを追加
			responseMessage := prepareFinishResponse(distance, fuelCost, coordinates[0].Timestamp, coordinates[len(coordinates)-1].Timestamp)
			if err := sendMessageFinish(conn, responseMessage); err != nil {
				log.Println("Error sending finish message:", err)
			}
			isSendFinishResponse = true
			break
		} else {
			responseMessage := utils.ResponseMessage{
				Status:              "SUCCESS",
				ProcessingStartTime: processingStartTime,
				TraceID:             traceID,
			}
			if err := sendMessage(conn, responseMessage); err != nil {
				log.Println("Error sending message:", err)
			}
		}
		errorRepeatCount = 0
	}

	// トリップデータを保存（この部分もエラーハンドリングを追加することが望ましい）
	if err := services.CreateTrip(db, startCoords, endCoords); err != nil {
		log.Println("Error saving trip to database:", err)
	} else {
		log.Println("Trip saved to database")
	}
}

// Data エラー処理
func handleConnectionError(err error) {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		log.Println("Client disconnected gracefully:", err)
	} else {
		log.Println("Error reading message:", err)
	}
}

// 座標をデコードする関数
func decodeCoordinates(msg []byte) (utils.Coordinates, error) {
	var coords utils.Coordinates
	if err := json.Unmarshal(msg, &coords); err != nil {
		log.Println("Error decoding JSON:", err)
		return utils.Coordinates{}, err
	}

	if err := validate.Struct(coords); err != nil {
		log.Println("Validation error:", err)
		return utils.Coordinates{}, err
	}
	return coords, nil
}

// バリデーションエラーを送信する関数
func sendValidationError(conn *websocket.Conn) {
	responseMessage := utils.ResponseMessage{
		Status:              "VALIDATE_ERROR",
		ProcessingStartTime: "",
		TraceID:             "",
	}
	sendMessage(conn, responseMessage)
}

// 開始メッセージを処理する関数
func handleStartMessage(startCoords *utils.Coordinates, isFirstMessage *bool, coords utils.Coordinates) (string, string) {
	var processingStartTime string
	var traceID string

	if *isFirstMessage {
		*startCoords = coords
		processingStartTime = time.Now().Format(time.RFC3339)
		traceID = utils.GenerateUniqueID()
		*isFirstMessage = false
	} else {
		processingStartTime = coords.ProcessingStartTime
		traceID = coords.TraceID
	}
	return processingStartTime, traceID
}

// 終了時のメッセージ送信
func sendMessageFinish(conn *websocket.Conn, message utils.TrackFinishResponse) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	log.Println("Sent message to client: sendMessageFinish")
	return nil
}

// メッセージ送信のヘルパー関数
func sendMessage(conn *websocket.Conn, message utils.ResponseMessage) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	log.Println("Sent message to client:", message.TraceID)
	return nil
}

// トリップ終了時のレスポンスメッセージを準備する関数
func prepareFinishResponse(distance float64, fuelCost float64, startTimeStr string, endTimeStr string) utils.TrackFinishResponse {

	formattedStartTime, _ := utils.FormatTime(startTimeStr)
	formattedEndTime, _ := utils.FormatTime(endTimeStr)

	return utils.TrackFinishResponse{
		Distance:  distance,
		FuelCost:  fuelCost,
		StartTime: formattedStartTime,
		EndTime:   formattedEndTime,
	}
}
