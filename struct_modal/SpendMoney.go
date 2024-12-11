package struct_modal

import "google.golang.org/api/sheets/v4"

type SpendMoney struct {
	Money float64 `json:"money"`
	Note  string  `json:"note"`
}

type DataRows struct {
	ValueRange *sheets.ValueRange `json:"value_range"`
}

type RowResponse struct {
	//TotalMoney    float64 `json:"total_money"`
	DailySpending   float64            `json:"daily_spending"`
	MonthlySpending float64            `json:"monthly_spending"`
	Reason          map[float64]string `json:"reason"`
}
