package console

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
)

type QueryRevenueStatisticsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewQueryRevenueStatisticsLogic Query revenue statistics
func NewQueryRevenueStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRevenueStatisticsLogic {
	return &QueryRevenueStatisticsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRevenueStatisticsLogic) QueryRevenueStatistics() (resp *types.RevenueStatisticsResponse, err error) {
	if strings.ToLower(os.Getenv("PPANEL_MODE")) == "demo" {
		return l.mockRevenueStatistics(), nil
	}

	var today, monthly, all types.OrdersStatistics
	now := time.Now()
	// Get today's revenue statistics
	todayData, err := l.svcCtx.OrderModel.QueryDateOrders(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryRevenueStatisticsLogic] QueryDateOrders error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryDateOrders error: %v", err)
	} else {
		today = types.OrdersStatistics{
			AmountTotal:        todayData.AmountTotal,
			NewOrderAmount:     todayData.NewOrderAmount,
			RenewalOrderAmount: todayData.RenewalOrderAmount,
		}
	}
	// Get monthly's revenue statistics
	monthlyData, err := l.svcCtx.OrderModel.QueryMonthlyOrders(l.ctx, now)
	if err != nil {
		l.Errorw("[QueryRevenueStatisticsLogic] QueryDateOrders error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryDateOrders error: %v", err)
	} else {
		monthly = types.OrdersStatistics{
			AmountTotal:        monthlyData.AmountTotal,
			NewOrderAmount:     monthlyData.NewOrderAmount,
			RenewalOrderAmount: monthlyData.RenewalOrderAmount,
			List:               make([]types.OrdersStatistics, 0),
		}
	}

	// Get all revenue statistics
	allData, err := l.svcCtx.OrderModel.QueryTotalOrders(l.ctx)
	if err != nil {
		l.Errorw("[QueryRevenueStatisticsLogic] QueryTotalOrders error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryTotalOrders error: %v", err)
	} else {
		all = types.OrdersStatistics{
			AmountTotal:        allData.AmountTotal,
			NewOrderAmount:     allData.NewOrderAmount,
			RenewalOrderAmount: allData.RenewalOrderAmount,
			List:               make([]types.OrdersStatistics, 0),
		}
	}
	return &types.RevenueStatisticsResponse{
		Today:   today,
		Monthly: monthly,
		All:     all,
	}, nil
}

// mockRevenueStatistics is a mock function to simulate revenue statistics data.
func (l *QueryRevenueStatisticsLogic) mockRevenueStatistics() *types.RevenueStatisticsResponse {
	now := time.Now()

	// Generate daily data for the past 7 days (oldest first)
	monthlyList := make([]types.OrdersStatistics, 7)
	for i := 0; i < 7; i++ {
		dayDate := now.AddDate(0, 0, -(6 - i))
		baseAmount := int64(25000 + ((6 - i) * 3000) + ((6-i)%3)*8000)
		monthlyList[i] = types.OrdersStatistics{
			Date:               dayDate.Format("2006-01-02"),
			AmountTotal:        baseAmount,
			NewOrderAmount:     int64(float64(baseAmount) * 0.68),
			RenewalOrderAmount: int64(float64(baseAmount) * 0.32),
		}
	}

	// Generate monthly data for the past 6 months (oldest first)
	allList := make([]types.OrdersStatistics, 6)
	for i := 0; i < 6; i++ {
		monthDate := now.AddDate(0, -(5 - i), 0)
		baseAmount := int64(1800000 + ((5 - i) * 200000) + ((5-i)%2)*500000)
		allList[i] = types.OrdersStatistics{
			Date:               monthDate.Format("2006-01"),
			AmountTotal:        baseAmount,
			NewOrderAmount:     int64(float64(baseAmount) * 0.68),
			RenewalOrderAmount: int64(float64(baseAmount) * 0.32),
		}
	}

	return &types.RevenueStatisticsResponse{
		Today: types.OrdersStatistics{
			AmountTotal:        35888,
			NewOrderAmount:     22888,
			RenewalOrderAmount: 13000,
		},
		Monthly: types.OrdersStatistics{
			AmountTotal:        888888,
			NewOrderAmount:     588888,
			RenewalOrderAmount: 300000,
			List:               monthlyList,
		},
		All: types.OrdersStatistics{
			AmountTotal:        12888888,
			NewOrderAmount:     8588888,
			RenewalOrderAmount: 4300000,
			List:               allList,
		},
	}
}
