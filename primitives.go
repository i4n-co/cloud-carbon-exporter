package cloudcarbonexporter

import (
	"gonum.org/v1/gonum/stat/distuv"
)

func CPUWatts(percent float64) float64 {
	wattPerCoreAdjuster := 16.5
	lambda := 3.0
	coef := 1.3

	exp := new(distuv.Exponential)
	exp.Rate = lambda

	return wattPerCoreAdjuster + wattPerCoreAdjuster*exp.CDF(float64(percent/100))*coef
}
