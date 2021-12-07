package handlers

import (
	"diplom/generator"
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

	probabilityFloat, err := strconv.ParseFloat(probabilityString, 32)
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
	var resultString string
	i := 0
	for _, v := range arr {
		if i == 30 {
			i = 0
			resultString += "\n"
		}
		resultString += fmt.Sprintf("%d", int(v)) + "\t"
		i++
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

	return gistCols, nil
}