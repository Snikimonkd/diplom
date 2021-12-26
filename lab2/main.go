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

	var ret []generator.OneProductModel

	table := widget.NewTable(
		func() (int, int) {
			return len(ret) + 1, 7
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("   ")
		},
		func(i widget.TableCellID, obj fyne.CanvasObject) {
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
					obj.(*widget.Label).SetText(fmt.Sprintf("%d", ret[i.Row-1].ExpNumber))
				case 1:
					obj.(*widget.Label).SetText(fmt.Sprintf("%d", ret[i.Row-1].ROP))
				case 2:
					obj.(*widget.Label).SetText(fmt.Sprintf("%d", ret[i.Row-1].EOQ))
				case 3:
					obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", ret[i.Row-1].TC1))
				case 4:
					obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", ret[i.Row-1].TC2))
				case 5:
					obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", ret[i.Row-1].TC3))
				case 6:
					obj.(*widget.Label).SetText(fmt.Sprintf("%.2f", ret[i.Row-1].TC))
				}

			}
		})
	table.SetColumnWidth(0, 30)
	table.SetColumnWidth(1, 50)
	table.SetColumnWidth(2, 50)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 100)
	table.SetColumnWidth(5, 100)
	table.Resize(fyne.NewSize(models.WindowWidth-models.DefaultWidth-40, 600))
	table.Move(fyne.NewPos(models.DefaultHorizontalPadding, 3*models.DefaultVerticalPadding+150))

	table.OnSelected = func(id widget.TableCellID) {
		ret = generator.SortModels(ret, id.Col)
		table.Refresh()
	}

	objects = append(objects, table)

	oneProductContainer := container.NewWithoutLayout(objects...)

	button.OnTapped = func() {
		var err error
		ret, err = handlers.SimModelingHandler(input0.Text, input1.Text, output)
		if err != nil {
			errorPopUp := widget.NewPopUp(widget.NewLabel(err.Error()), myWindow.Canvas())
			errorPopUp.Move(fyne.NewPos(200, 200))
			errorPopUp.Show()
		}
		table.Refresh()

	}

	tabCont := container.NewTabItem("Многопродуктовая система", container.NewWithoutLayout([]fyne.CanvasObject{}...))

	tabs := container.NewAppTabs(
		container.NewTabItem("Однопродуктовая система", oneProductContainer),
		tabCont,
	)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
