package conventions

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

// Delivery структура для представления информации о доставке
type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// Payment структура для представления информации о платеже
type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

// Item структура для представления информации о товаре
type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

// generateRandomString генерирует случайную строку заданной длины
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomOrder создает объект Order со случайными данными
func GenerateRandomOrder() Order {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var items []Item

	for range seededRand.Intn(3) + 1 {
		item := Item{
			ChrtID:      seededRand.Intn(1000000),
			TrackNumber: fmt.Sprintf("TRACK%d", seededRand.Intn(1000000)),
			Price:       seededRand.Intn(1000),
			Rid:         fmt.Sprintf("RID%d", seededRand.Intn(1000000)),
			Name:        generateRandomString(10),
			Sale:        seededRand.Intn(100),
			Size:        generateRandomString(2),
			TotalPrice:  seededRand.Intn(1000),
			NmID:        seededRand.Intn(1000000),
			Brand:       generateRandomString(10),
			Status:      seededRand.Intn(10),
		}

		items = append(items, item)

	}

	order := Order{
		OrderUID:    uuid.New().String(),
		TrackNumber: fmt.Sprintf("TRACK%d", seededRand.Intn(1000000)),
		Entry:       generateRandomString(4),
		Delivery: Delivery{
			Name:    generateRandomString(10),
			Phone:   fmt.Sprintf("+%d", seededRand.Intn(10000000000)),
			Zip:     fmt.Sprintf("%d", seededRand.Intn(1000000)),
			City:    generateRandomString(10),
			Address: generateRandomString(20),
			Region:  generateRandomString(10),
			Email:   fmt.Sprintf("%s@%s.com", generateRandomString(5), generateRandomString(5)),
		},
		Payment: Payment{
			Transaction:  fmt.Sprintf("TRANS%d", seededRand.Intn(1000000)),
			RequestID:    generateRandomString(10),
			Currency:     generateRandomString(3),
			Provider:     generateRandomString(10),
			Amount:       seededRand.Intn(10000),
			PaymentDT:    time.Now().Unix(),
			Bank:         generateRandomString(10),
			DeliveryCost: seededRand.Intn(5000),
			GoodsTotal:   seededRand.Intn(1000),
			CustomFee:    seededRand.Intn(100),
		},
		Items:             items,
		Locale:            generateRandomString(2),
		InternalSignature: generateRandomString(10),
		CustomerID:        generateRandomString(10),
		DeliveryService:   generateRandomString(10),
		Shardkey:          generateRandomString(1),
		SmID:              seededRand.Intn(100),
		DateCreated:       time.Now(),
		OofShard:          generateRandomString(1),
	}

	return order
}
