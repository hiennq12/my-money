package caculator_data

import (
	"errors"
	"fmt"
	"github.com/hiennq12/my-money/struct_modal"
	"github.com/hiennq12/my-money/util"
	"log"
	"strconv"
	"strings"
	"time"
)

func MoneySpendInDay(rowsValue *struct_modal.DataRows) (*struct_modal.RowResponse, error) {
	if rowsValue == nil || rowsValue.ValueRange == nil {
		return nil, errors.New("rows data is empty, check sheet")
	}

	resp := rowsValue.ValueRange
	now := time.Now()
	totalMoney := float64(0)
	mapMoneyReason := make(map[float64]string)
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for rowIndex, row := range resp.Values {
			if len(row) < 1 {
				continue
			}

			dayInRow, err := strconv.Atoi(row[0].(string))
			if err != nil {
				log.Fatalf("Parse day in row error: %v", err.Error())
				return nil, err
			}

			if dayInRow != now.Day() {
				continue
			}
			totalMoney, mapMoneyReason = handleRowData(row, rowIndex)
		}
	}

	return &struct_modal.RowResponse{
		TotalMoney: totalMoney,
		Reason:     mapMoneyReason,
	}, nil
}

func handleRowData(row []interface{}, rowIndex int) (float64, map[float64]string) {
	firstCell := true
	// get current month and year
	totalMoney := float64(0)
	mapMoneyReason := make(map[float64]string)
	for columnIndex, cell := range row {
		dataCell := cell.(string)
		addrCell := fmt.Sprintf("%v%v", util.MapColumn[columnIndex], rowIndex+1)
		if len(dataCell) > 0 {
			if firstCell && len(dataCell) > 0 {
				util.GetDateFromFirstCell(dataCell)
				firstCell = false
				continue
			}
			spendMoney, err := handleCellData(dataCell, addrCell)
			if err != nil {
				log.Panic(fmt.Sprintf("Error when parse data in cell[%v]. Detail err: [%v]", addrCell, err.Error()))
			}

			totalMoney += spendMoney.Money
			mapMoneyReason[spendMoney.Money] = spendMoney.Note
		}
	}

	return totalMoney, mapMoneyReason
}

func handleCellData(cell string, addrCell string) (*struct_modal.SpendMoney, error) {
	allData := strings.Split(cell, " ")
	if len(allData) < 1 {
		return nil, errors.New(fmt.Sprintf("The data in the cell has format errors (in %v)", addrCell))
	}

	money, err := strconv.ParseFloat(allData[0], 64)
	if err != nil {
		log.Panicf("Error when parse money in [%v]. Detail: %v", addrCell, err.Error())
		return nil, err
	}
	return &struct_modal.SpendMoney{
		Money: money / 1000,
		Note:  strings.Join(allData[1:], " "),
	}, nil
}
