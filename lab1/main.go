package main

import (
	"diplom/internal/handlers"
	"diplom/internal/models"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	os.Setenv("FYNE_THEME", "light")
	myApp := app.New()

	myWindow := myApp.NewWindow("Лабораторная работа №1")

	myWindow.Resize(fyne.Size{Width: models.WindowWidth, Height: models.WindowHeight})
	myWindow.SetFixedSize(true)

	objects := []fyne.CanvasObject{}

	label0 := models.CreateLabel("Закон распределения")
	objects = append(objects, label0)

	radio := widget.NewRadioGroup([]string{"Равномерный", "Нормальный", "Экспоненциальный", "Биномиальный"}, nil)
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
				input2.SetText("0")
				label2.Refresh()
				label2.Show()
				input2.SetPlaceHolder("Нижняя граница")
				input2.Refresh()
				label3.Text = "Верхняя граница"
				label3.Refresh()
				label3.Show()
				input3.SetPlaceHolder("Верхняя граница")
				input3.SetText("1")
				input3.Refresh()
				input3.Show()
				label6.Show()
				input6.Show()
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
				label6.Show()
				input6.Show()
			}
		case "Экспоненциальный":
			{

				label2.Text = "Параметр λ"
				label2.Refresh()
				input2.SetPlaceHolder("Параметр λ")
				input2.SetText("10")
				input2.Refresh()
				label3.Hide()
				input3.Hide()
				label6.Show()
				input6.Show()

			}
		case "Биномиальный":
			{
				label2.Text = "Вероятность"
				label2.Refresh()
				input2.SetPlaceHolder("Вероятность")
				input2.SetText("0.5")
				input2.Refresh()
				label3.Text = "Степень полинома"
				label3.Show()
				label3.Refresh()
				input3.SetPlaceHolder("Степень полинома")
				input3.Show()
				input3.SetText("4")
				input3.Refresh()
				label6.Hide()
				input6.Hide()
			}
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
				newGistCols, err := handlers.LinearButtonHandler(input1.Text, input2.Text, input3.Text, input4.Text, input5.Text, input6.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}

				tabCont.Content = container.NewWithoutLayout(newGistCols...)
			}
		case "Нормальный":
			{
				newGistCols, err := handlers.NormalButtonHandler(input1.Text, input2.Text, input3.Text, input4.Text, input5.Text, input6.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}

				tabCont.Content = container.NewWithoutLayout(newGistCols...)
			}
		case "Экспоненциальный":
			{
				newGistCols, err := handlers.ExpButtonHandler(input1.Text, input2.Text, input4.Text, input5.Text, input6.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}

				tabCont.Content = container.NewWithoutLayout(newGistCols...)
			}
		case "Биномиальный":
			{
				newGistCols, err := handlers.BinomialButtonHandler(input1.Text, input2.Text, input3.Text, input4.Text, input5.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}

				tabCont.Content = container.NewWithoutLayout(newGistCols...)
			}
		}
	}

	myWindow.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		key := *ke

		if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
			button.Tapped(nil)
		}
	})

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
