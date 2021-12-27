package main

import (
	"diplom/internal/generator"
	"diplom/internal/handlers"
	"diplom/internal/models"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	os.Setenv("FYNE_THEME", "light")
	myApp := app.New()

	myWindow := myApp.NewWindow("Лабораторная работа №2")

	myWindow.Resize(fyne.Size{Width: models.WindowWidth, Height: models.WindowHeight})
	myWindow.SetFixedSize(true)

	objects := []fyne.CanvasObject{}

	radioLabel := models.CreateLabel("Вид системы")
	objects = append(objects, radioLabel)

	radio := widget.NewRadioGroup([]string{"Однопродуктовая", "Многопродуктовая"}, nil)
	radio.SetSelected("Однопродуктовая")
	radio.Resize(radio.MinSize())
	objects = append(objects, radio)

	label0 := models.CreateLabel("Вариант")
	objects = append(objects, label0)

	input0 := models.CreateInput("Вариант", "7182021")
	objects = append(objects, input0)

	label1 := models.CreateLabel("Кол-во экспериментов")
	objects = append(objects, label1)

	input1 := models.CreateInput("Кол-во экспериментов", "10")
	objects = append(objects, input1)

	output := widget.NewMultiLineEntry()
	output.Resize(fyne.NewSize(models.WindowWidth-models.DefaultWidth-40, 150))
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

	var retOne []generator.OneProductModel
	var retMulti []generator.MultiProductModel

	tableOne := widget.NewTable(
		func() (int, int) {
			switch radio.Selected {
			case "Однопродуктовая":
				{
					return len(retOne) + 1, 7
				}
			case "Многопродуктовая":
				{
					return len(retMulti) + 1, 4
				}
			}
			return 0, 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("   ")
		},
		func(i widget.TableCellID, obj fyne.CanvasObject) {
			switch radio.Selected {
			case "Однопродуктовая":
				{
					if i.Row == 0 {
						switch i.Col {
						case 0:
							obj.(*widget.Label).SetText("№")
						case 1:
							obj.(*widget.Label).SetText("ROP")
						case 2:
							obj.(*widget.Label).SetText("EOQ")
						case 3:
							obj.(*widget.Label).SetText("TC1")
						case 4:
							obj.(*widget.Label).SetText("TC2")
						case 5:
							obj.(*widget.Label).SetText("TC3")
						case 6:
							obj.(*widget.Label).SetText("TC")
						}
					} else {
						switch i.Col {
						case 0:
							obj.(*widget.Label).SetText(fmt.Sprintf("%d", retOne[i.Row-1].ExpNumber))
						case 1:
							obj.(*widget.Label).SetText(fmt.Sprintf("%d", retOne[i.Row-1].ROP))
						case 2:
							obj.(*widget.Label).SetText(fmt.Sprintf("%d", retOne[i.Row-1].EOQ))
						case 3:
							obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", retOne[i.Row-1].TC1))
						case 4:
							obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", retOne[i.Row-1].TC2))
						case 5:
							obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", retOne[i.Row-1].TC3))
						case 6:
							obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", retOne[i.Row-1].TC))
						}

					}
				}
			case "Многопродуктовая":
				{
					if i.Row == 0 {
						switch i.Col {
						case 0:
							obj.(*widget.Label).SetText("№")
						case 1:
							obj.(*widget.Label).SetText("TOC")
						case 2:
							obj.(*widget.Label).SetText("TCC")
						case 3:
							obj.(*widget.Label).SetText("TCOST")
						}
					} else {
						switch i.Col {
						case 0:
							obj.(*widget.Label).SetText(fmt.Sprintf("%d", retMulti[i.Row-1].ExpNumber))
						case 1:
							obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", retMulti[i.Row-1].TOC))
						case 2:
							obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", retMulti[i.Row-1].TCC))
						case 3:
							obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", retMulti[i.Row-1].TCOST))
						}

					}
				}
			}

		})
	tableOne.SetColumnWidth(0, 30)
	tableOne.SetColumnWidth(1, 50)
	tableOne.SetColumnWidth(2, 50)
	tableOne.SetColumnWidth(3, 100)
	tableOne.SetColumnWidth(4, 100)
	tableOne.SetColumnWidth(5, 100)
	tableOne.Resize(fyne.NewSize(models.WindowWidth-models.DefaultWidth-40, 600))
	tableOne.Move(fyne.NewPos(models.DefaultHorizontalPadding, 3*models.DefaultVerticalPadding+150))

	tableOne.OnSelected = func(id widget.TableCellID) {
		switch radio.Selected {
		case "Однопродуктовая":
			{
				retOne = generator.SortModelsOne(retOne, id.Col)
				tableOne.Refresh()
			}
		case "Многопродуктовая":
			{
				//retMulti = generator.SortModelsMulti(retMulti, id.Col)
				tableOne.Refresh()
			}
		}
	}

	radio.OnChanged = func(value string) {
		switch value {
		case "Однопродуктовая":
			{
				tableOne.SetColumnWidth(0, 30)
				tableOne.SetColumnWidth(1, 50)
				tableOne.SetColumnWidth(2, 50)
				tableOne.SetColumnWidth(3, 100)
				tableOne.SetColumnWidth(4, 100)
				tableOne.SetColumnWidth(5, 100)
			}
		case "Многопродуктовая":
			{
				tableOne.SetColumnWidth(0, 30)
				tableOne.SetColumnWidth(1, 100)
				tableOne.SetColumnWidth(2, 100)
				tableOne.SetColumnWidth(3, 100)
			}
		}
		tableOne.Refresh()
	}

	objects = append(objects, tableOne)

	mainContainer := container.NewWithoutLayout(objects...)

	button.OnTapped = func() {
		switch radio.Selected {
		case "Однопродуктовая":
			{
				var err error
				retOne, err = handlers.OneProductHandler(input0.Text, input1.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}
				tableOne.Refresh()
			}
		case "Многопродуктовая":
			{
				var err error
				retMulti, err = handlers.MultiProductHandler(input0.Text, input1.Text, output)
				if err != nil {
					errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
					errorPopUp.Move(fyne.NewPos(200, 200))
					errorPopUp.Show()
				}
				tableOne.Refresh()
			}
		}
	}

	myWindow.SetContent(mainContainer)
	myWindow.ShowAndRun()
}
