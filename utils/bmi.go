package utils

func CalculateBMI(weight float64, height float64) float64 {
	h := height / 100
	return weight / (h * h)
}

func GetBMICategory(bmi float64) string {
	if bmi < 18.5 {
		return "underweight"
	} else if bmi < 25 {
		return "normal"
	} else if bmi < 30 {
		return "overweight"
	}
	return "obese"
}