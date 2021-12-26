package handlers

import (
	"diplom/internal/generator"
	"diplom/internal/models"
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func BinomialButtonHandler(baseString, probabilityString, levelString, amountToGenerateString, amountToDrawString string, output *widget.Entry) ([]fyne.CanvasObject, error) {
	baseInt, err := strconv.Atoi(baseString)
	if err != nil {
		return nil, errors.New("неправильно задана база генератора")
	}

	probabilityFloat, err := strconv.ParseFloat(probabilityString, 64)
	if err != nil {
		return nil, errors.New("неправильно задана вероятность")
	}

	if probabilityFloat > 1 || probabilityFloat < 0 {
		return nil, errors.New("неправильно задана вероятность")
	}

	levelInt, err := strconv.Atoi(levelString)
	if err != nil {
		return nil, errors.New("неправилтно задана степень полинома")
	}

	if levelInt < 1 {
		return nil, errors.New("степень полинома не может быть меньше 1")
	}

	amountToGenerateInt, err := strconv.Atoi(amountToGenerateString)
	if err != nil {
		return nil, errors.New("неправильно задано количество для генерации")
	}

	amountToDrawInt, err := strconv.Atoi(amountToDrawString)
	if err != nil {
		return nil, errors.New("неправильно задан размер выборки")
	}

	if amountToDrawInt > amountToGenerateInt {
		return nil, errors.New("размер выборки не может быть больше количества для генерации")
	}

	arr := generator.Binomial(probabilityFloat, levelInt, amountToGenerateInt, baseInt)

	resultString := "Теоретические значения вероятностей для полинома " + levelString + " степени:\n"

	for i := 0; i <= levelInt; i++ {
		resultString += "P(" + strconv.Itoa(i) + ")=" + fmt.Sprintf("%f\n", generator.Ver(probabilityFloat, levelInt, i))
	}

	i := 0
	for _, v := range arr {
		if i == 10 {
			i = 0
			resultString += "\n"
		}
		resultString += fmt.Sprintf("%d", int(v)) + "\t|  "
		i++
	}
	resultString += "\n"

	resAmount := generator.Count(arr, levelInt)

	for i, v := range resAmount {
		resultString += "Кол-во сгенерированных " + fmt.Sprint(i) + " : " + fmt.Sprint(v) + "\n"
	}

	gistCols := generator.Draw(arr, levelInt+1)

	output.SetText(resultString)

	return gistCols, nil
}

func ExpButtonHandler(baseString, lyambdaString, amountToGenerateString, amountToDrawString, colsString string, output *widget.Entry) ([]fyne.CanvasObject, error) {
	baseInt, err := strconv.Atoi(baseString)
	if err != nil {
		return nil, errors.New("неправильно задана база генератора")
	}

	lyambdaInt, err := strconv.Atoi(lyambdaString)
	if err != nil {
		return nil, errors.New("неправильно задана лямбда")
	}
	if lyambdaInt <= 0 {
		return nil, errors.New("лямбда не может быть меньше или равна 0")
	}

	amountToGenerateInt, err := strconv.Atoi(amountToGenerateString)
	if err != nil {
		return nil, errors.New("неправильно задано количество для генерации")
	}

	amountToDrawInt, err := strconv.Atoi(amountToDrawString)
	if err != nil {
		return nil, errors.New("неправильно задан размер выборки")
	}

	if amountToDrawInt > amountToGenerateInt {
		return nil, errors.New("размер выборки не может быть больше количества для генерации")
	}

	colsInt, err := strconv.Atoi(colsString)
	if err != nil {
		return nil, errors.New("неправильно задано количество столбцов гистограммы")
	}

	if colsInt > amountToDrawInt {
		return nil, errors.New("количество стобцов гистограммы не может быть больше выборки")
	}

	arr := generator.ExpGenerate(lyambdaInt, amountToGenerateInt, baseInt)
	var resultString string
	for _, v := range arr {
		resultString += fmt.Sprintf("%f", v) + "\n"
	}

	gistCols := generator.Draw(arr[:amountToDrawInt], colsInt)

	output.SetText(resultString)

	expValue := generator.ExpValue(arr)
	dispValue := generator.DispValue(arr)

	expValLabel := models.CreateLabel("Мат. ожидание = " + fmt.Sprintf("%.4f", expValue))
	dispValLabel := models.CreateLabel("Дисперсия = " + fmt.Sprintf("%.4f", dispValue))

	expValLabel.Move(fyne.NewPos(0, 0))
	dispValLabel.Move(fyne.NewPos(400, 0))

	gistCols = append(gistCols, expValLabel, dispValLabel)

	return gistCols, nil
}

func LinearButtonHandler(baseString, lowerString, upperString, amountToGenerateString, amountToDrawString, colsString string, output *widget.Entry) ([]fyne.CanvasObject, error) {
	baseInt, err := strconv.Atoi(baseString)
	if err != nil {
		return nil, errors.New("неправильно задана база генератора")
	}

	lowerInt, err := strconv.Atoi(lowerString)
	if err != nil {
		return nil, errors.New("неправильно задана нижняя граница")
	}

	upperInt, err := strconv.Atoi(upperString)
	if err != nil {
		return nil, errors.New("неправилтно задана верхняя граница")
	}

	if lowerInt > upperInt {
		return nil, errors.New("нижняя граница не может быть больше верхней")
	}

	amountToGenerateInt, err := strconv.Atoi(amountToGenerateString)
	if err != nil {
		return nil, errors.New("неправильно задано количество для генерации")
	}

	colsInt, err := strconv.Atoi(colsString)
	if err != nil {
		return nil, errors.New("неправильно задано количество столбцов гистограммы")
	}

	amountToDrawInt, err := strconv.Atoi(amountToDrawString)
	if err != nil {
		return nil, errors.New("неправильно задан размер выборки")
	}

	if amountToDrawInt > amountToGenerateInt {
		return nil, errors.New("размер выборки не может быть больше количества для генерации")
	}

	if colsInt > amountToDrawInt {
		return nil, errors.New("количество стобцов гистограммы не может быть больше выборки")
	}

	arr := generator.LinearGenerate(lowerInt, upperInt, amountToGenerateInt, baseInt)
	var resultString string
	for _, v := range arr {
		resultString += fmt.Sprintf("%f", v) + "\n"
	}

	gistCols := generator.Draw(arr[:amountToDrawInt], colsInt)

	output.SetText(resultString)

	expValue := generator.ExpValue(arr)
	dispValue := generator.DispValue(arr)

	expValLabel := models.CreateLabel("Мат. ожидание = " + fmt.Sprintf("%.4f", expValue))
	dispValLabel := models.CreateLabel("Дисперсия = " + fmt.Sprintf("%.4f", dispValue))

	expValLabel.Move(fyne.NewPos(0, 0))
	dispValLabel.Move(fyne.NewPos(400, 0))

	gistCols = append(gistCols, expValLabel, dispValLabel)

	return gistCols, nil
}

func NormalButtonHandler(baseString, mathExpectationString, dispersionString, amountToGenerateString, amountToDrawString, colsString string, output *widget.Entry) ([]fyne.CanvasObject, error) {
	baseInt, err := strconv.Atoi(baseString)
	if err != nil {
		return nil, errors.New("неправильно задана база генератора")
	}

	mathExpectationInt, err := strconv.Atoi(mathExpectationString)
	if err != nil {
		return nil, errors.New("неправильно задано мат. ожидание")
	}

	dispersionInt, err := strconv.Atoi(dispersionString)
	if err != nil {
		return nil, errors.New("неправилтно задана дисперсия")
	}

	if dispersionInt < 0 {
		return nil, errors.New("дисперсия не может быть меньше нуля")
	}

	amountToGenerateInt, err := strconv.Atoi(amountToGenerateString)
	if err != nil {
		return nil, errors.New("неправильно задано количество для генерации")
	}

	colsInt, err := strconv.Atoi(colsString)
	if err != nil {
		return nil, errors.New("неправильно задано количество столбцов гистограммы")
	}

	amountToDrawInt, err := strconv.Atoi(amountToDrawString)
	if err != nil {
		return nil, errors.New("неправильно задан размер выборки")
	}

	if amountToDrawInt > amountToGenerateInt {
		return nil, errors.New("размер выборки не может быть больше количества для генерации")
	}

	if colsInt > amountToDrawInt {
		return nil, errors.New("количество стобцов гистограммы не может быть больше выборки")
	}

	arr := generator.NormalGenerate(mathExpectationInt, dispersionInt, amountToGenerateInt, baseInt)
	var resultString string
	for _, v := range arr {
		resultString += fmt.Sprintf("%f", v) + "\n"
	}

	gistCols := generator.Draw(arr[:amountToDrawInt], colsInt)

	output.SetText(resultString)

	expValue := generator.ExpValue(arr)
	dispValue := generator.DispValue(arr)

	expValLabel := models.CreateLabel("Мат. ожидание = " + fmt.Sprintf("%.4f", expValue))
	dispValLabel := models.CreateLabel("Дисперсия = " + fmt.Sprintf("%.4f", dispValue))

	expValLabel.Move(fyne.NewPos(0, 0))
	dispValLabel.Move(fyne.NewPos(400, 0))

	gistCols = append(gistCols, expValLabel, dispValLabel)

	return gistCols, nil
}

func SimModelingHandler(variantString, amountString string, output *widget.Entry) ([]generator.Model, error) {
	variantInt, err := strconv.Atoi(variantString)
	if err != nil {
		return nil, errors.New("неправильно задан вариант")
	}

	amountInt, err := strconv.Atoi(amountString)
	if err != nil {
		return nil, errors.New("неправильно задано кол-во экспериментов")
	}

	ret := generator.Modeling(variantInt, amountInt)

	startingVluesString := "Исходные данные варианта:\n"
	startingVluesString += "C1=" + fmt.Sprintf("%.2f", ret[0].C1) + "\n"
	startingVluesString += "C2=" + fmt.Sprintf("%.2f", ret[0].C2) + "\n"
	startingVluesString += "C3=" + fmt.Sprintf("%.2f", ret[0].C3) + "\n"
	startingVluesString += "B1=" + strconv.Itoa(ret[0].B1) + "\n"
	startingVluesString += "TT=" + strconv.Itoa(ret[0].TT)
	output.SetText(startingVluesString)

	return ret, nil
}
