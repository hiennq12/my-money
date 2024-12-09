package noti

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hiennq12/my-money/struct_modal"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strings"
	"time"
)

// TelegramConfig holds the configuration for Telegram bot
type Config struct {
	TelegramConfig TelegramConfig `yaml:"telegram_config"`
}
type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

// TelegramMessage represents the message structure for Telegram API
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func readFileConfig(teleConfig *TelegramConfig) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	//var telegramConfig TelegramConfig
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return
	}

	teleConfig.BotToken = config.TelegramConfig.BotToken
	teleConfig.ChatID = config.TelegramConfig.ChatID
}

func PrepareData(moneyInDay *struct_modal.RowResponse) (*TelegramConfig, string) {
	// send message to tele
	// Configure your Telegram bot
	now := time.Now()
	nowString := now.Format("2006-01-02 15:04:05")
	configTele := &TelegramConfig{}
	readFileConfig(configTele)
	//configTele := TelegramConfig{
	//	BotToken: "7817584153:AAEQgSGiOE1TyouM6veW1VF1ExBg8CD1Vcw", // Replace with your bot token
	//	//ChatID:   "1022100822",                                     // Replace with your chat ID
	//	ChatID: "-4748590452", // Replace with your chat ID
	//}

	countSpendMoney := 0
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("So tien da tieu trong ngay [%v] la [%vk] \n", nowString, moneyInDay.TotalMoney))
	for k, v := range moneyInDay.Reason {
		countSpendMoney++
		builder.WriteString(fmt.Sprintf("[%v] %vk : %v \n", countSpendMoney, k, v))
	}

	message := builder.String()
	return configTele, message
}

// SendTelegramMessage sends a message to a specified Telegram chat
func SendTelegramMessage(config *TelegramConfig, message string) error {
	// Construct the Telegram Bot API URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.BotToken)

	// Create the message payload
	payload := TelegramMessage{
		ChatID: config.ChatID,
		Text:   message,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
