package main

import (
	"diplom/generator"
	"diplom/models"
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func expButtonHandler(baseString, lyambdaString, amountToGenerateString, amountToDrawString, colsString string, output *widget.Entry) ([]fyne.CanvasObject, error) {
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

func linearButtonHandler(baseString, lowerString, upperString, amountToGenerateString, amountToDrawString, colsString string, output *widget.Entry) ([]fyne.CanvasObject, error) {
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

func normalButtonHandler(baseString, mathExpectationString, dispersionString, amountToGenerateString, amountToDrawString, colsString string, output *widget.Entry) ([]fyne.CanvasObject, error) {
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

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Laba")
	myWindow.Resize(fyne.Size{Width: models.WindowWidth, Height: models.WindowHeight})
	myWindow.SetFixedSize(true)

	objects := []fyne.CanvasObject{}

	label0 := models.CreateLabel("Закон распределения")
	objects = append(objects, label0)

	radio := widget.NewRadioGroup([]string{"Равномерный", "Нормальный", "Экспоненциальный", "Дискретный"}, nil)
	radio.SetSelected("Равномерный")
	radio.Resize(radio.MinSize())
	objects = append(objects, radio)

	label1 := models.CreateLabel("База генератора")
	objects = append(objects, label1)

	input1 := models.CreateInput("База генератора", "12345678")
	objects = append(objects, input1)

	label2 := models.CreateLabel("Нижняя граница")
	objects = append(objects, label2)

	input2 := models.CreateInput("Нижняя граница", "0")
	objects = append(objects, input2)

	label3 := models.CreateLabel("Верхняя граница")
	objects = append(objects, label3)

	input3 := models.CreateInput("Верхняя граница", "1")
	objects = append(objects, input3)

	label4 := models.CreateLabel("Количество для генерации")
	objects = append(objects, label4)

	input4 := models.CreateInput("Количество для генерации", "100")
	objects = append(objects, input4)

	label5 := models.CreateLabel("Выборка")
	objects = append(objects, label5)

	input5 := models.CreateInput("Выборка", "100")
	objects = append(objects, input5)

	label6 := models.CreateLabel("Кол-во столбцов гистограммы")
	objects = append(objects, label6)

	input6 := models.CreateInput("Кол-во столбцов гистограммы", "16")
	objects = append(objects, input6)

	output := widget.NewMultiLineEntry()
	output.Resize(fyne.NewSize(models.WindowWidth-models.DefaultWidth-40, models.WindowHeight-50))
	output.Move(fyne.NewPos(models.DefaultHorizontalPadding, models.DefaultVerticalPadding))
	output.Disable()

	button := models.CreateButton("Сгенерировать")
	objects = append(objects, button)

	sumHeight := float32(0)
	for i, v := range objects {
		if i != 0 {
			sumHeight += objects[i-1].Size().Height
			v.Move(fyne.NewPos(models.WindowWidth-models.DefaultWidth-models.DefaultHorizontalPadding, sumHeight+float32(i)*models.DefaultVerticalPadding))
		} else {
			v.Move(fyne.NewPos(models.WindowWidth-models.DefaultWidth-models.DefaultHorizontalPadding, 0))
		}
	}

	objects = append(objects, output)

	mainContainer := container.NewWithoutLayout(objects...)

	radio.OnChanged = func(value string) {
		switch value {
		case "Равномерный":
			{
				label2.Text = "Нижняя граница"
				label2.Refresh()
				input2.SetPlaceHolder("Нижняя граница")
				input2.Refresh()
				label3.Text = "Верхняя граница"
				label3.Show()
				label3.Refresh()
				input3.SetPlaceHolder("Верхняя граница")
				input3.Show()
				input3.Refresh()
			}
		case "Нормальный":
			{
				label2.Text = "Мат. ожидание"
				label2.Refresh()
				input2.SetPlaceHolder("Мат. ожидание")
				label3.Show()
				input2.Refresh()
				label3.Text = "Дисперсия"
				label3.Refresh()
				input3.SetPlaceHolder("Дисперсия")
				input3.Show()
				input3.Refresh()
			}
		case "Экспоненциальный":
			{
				{
					label2.Text = "Параметр λ"
					label2.Refresh()
					input2.SetPlaceHolder("Параметр λ")
					input2.SetText("10")
					input2.Refresh()
					label3.Hide()
					input3.Hide()
				}
			}
		case "Дискретный":

		}
	}

	tabCont := container.NewTabItem("Гистограмма", container.NewWithoutLayout([]fyne.CanvasObject{}...))

	tabs := container.NewAppTabs(
		container.NewTabItem("Сгенерированные числа", mainContainer),
		tabCont,
	)

	button.OnTapped = func() {
		switch radio.Selected {
		case "Равномерный":
			{
				newGistCols, err := linearButtonHandler(input1.Text, input2.Text, input3.Text, input4.Text, input5.Text, input6.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}

				tabCont.Content = container.NewWithoutLayout(newGistCols...)
			}
		case "Нормальный":
			{
				newGistCols, err := normalButtonHandler(input1.Text, input2.Text, input3.Text, input4.Text, input5.Text, input6.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}

				tabCont.Content = container.NewWithoutLayout(newGistCols...)
			}
		case "Экспоненциальный":
			{
				newGistCols, err := expButtonHandler(input1.Text, input2.Text, input4.Text, input5.Text, input6.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}

				tabCont.Content = container.NewWithoutLayout(newGistCols...)
			}
		}
	}

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
