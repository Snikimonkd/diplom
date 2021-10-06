package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

var seed uint64

const windowHeight = 600
const windowWidth = 700

const a = uint64(6364136223846793005)
const c = uint64(1442695040888963407)
const m = uint64(18446744073709551615)

func lcg() uint64 {
	seed = (a*seed + c) % m
	return seed
}

func linearRandom(lower float64, upper float64) float64 {
	return float64(lcg())/float64(math.MaxUint64)*(upper-lower) + lower
}

func normalRandom(expected float64, variance float64) (float64, float64) {
	var s, x, y float64
	s = 0
	for s == 0 || s > 1 {
		x = linearRandom(-1, 1)
		y = linearRandom(-1, 1)
		s = x*x + y*y
	}

	z0 := expected + x*math.Sqrt(-2*math.Log(s)/s)*math.Sqrt(variance)
	z1 := expected + y*math.Sqrt(-2*math.Log(s)/s)*math.Sqrt(variance)

	return z0, z1
}

func findMinAndMaxFloat(a []float64) (min float64, max float64) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func findMinAndMaxInt(a []int) (min int, max int) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func mathExpectation(arr []float64) float64 {
	expectation := 0.0
	len := float64(len(arr))
	for _, v := range arr {
		expectation += v / len
	}

	return expectation
}

func mathDeviation(arr []float64) float64 {
	expectation := 0.0
	buf := 0.0
	len := float64(len(arr))
	for _, v := range arr {
		buf = buf + v*v/len
		expectation += v / len
	}

	return buf - expectation*expectation
}

func draw(cr *cairo.Context, arr []float64) {
	cr.SetSourceRGB(255, 255, 255)
	cr.Rectangle(0, 0, windowWidth, windowHeight)
	cr.Fill()
	_, max := findMinAndMaxFloat(arr)
	length1 := float64(windowWidth-10) / float64(len(arr))
	height1 := float64(windowHeight-10) / max

	for i, v := range arr {
		cr.SetSourceRGB(0, 0, 0)
		cr.Rectangle(float64(i)*length1, windowHeight-height1*v, length1, windowHeight)
	}

	cr.Fill()
}

func distribution(arr []float64, min float64, grouplen float64, groups int) []int {
	var distr []int
	distr = append(distr, 0)
	i := 0
	for _, v := range arr {
		if v < (float64(i+1)*grouplen + min) {
			distr[i]++
		} else {
			i++
			distr = append(distr, 0)
			distr[i]++
		}
	}

	return distr
}

func draw2(cr *cairo.Context, arr []float64, groupsStr string) {
	groupsInt, err := strconv.Atoi(groupsStr)
	if err != nil {
		fmt.Println("error getting groupsInt")
		return
	}

	sort.Float64s(arr)

	min := arr[0]
	max := arr[len(arr)-1]

	grouplen := (max - min) / float64(groupsInt)

	distr := distribution(arr, min, grouplen, groupsInt)

	cr.SetSourceRGB(255, 255, 255)
	cr.Rectangle(0, 0, windowWidth, windowHeight)
	cr.Fill()

	_, maxGroup := findMinAndMaxInt(distr)

	length1 := float64(windowWidth-10) / float64(groupsInt)
	height1 := float64(windowHeight-10) / float64(maxGroup)

	for i, v := range distr {
		cr.SetSourceRGB(0, 0, 0)
		cr.Rectangle(float64(i)*length1, windowHeight-height1*float64(v), length1, windowHeight)
	}

	cr.Fill()
}

func linearArr(baseStr string, lowerStr string, upperStr string, amountStr string) ([]float64, error) {
	baseInt, err := strconv.Atoi(baseStr)
	if err != nil {
		return nil, err
	}

	lowerFloat, err := strconv.ParseFloat(lowerStr, 64)
	if err != nil {
		return nil, err
	}

	upperFloat, err := strconv.ParseFloat(upperStr, 64)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return nil, err
	}

	var ret []float64
	seed = uint64(baseInt)
	for i := 0; i < amount; i++ {
		newRand := linearRandom(lowerFloat, upperFloat)
		ret = append(ret, newRand)
	}

	return ret, nil
}

func normalArray(baseStr string, expectedStr string, varianceStr string, amountStr string) ([]float64, error) {
	baseInt, err := strconv.Atoi(baseStr)
	if err != nil {
		return nil, err
	}

	expectedFloat, err := strconv.ParseFloat(expectedStr, 64)
	if err != nil {
		return nil, err
	}

	varianceFloat, err := strconv.ParseFloat(varianceStr, 64)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return nil, err
	}

	var ret []float64
	seed = uint64(baseInt)
	for i := 0; i < amount/2; i++ {
		z0, z1 := normalRandom(expectedFloat, varianceFloat)
		ret = append(ret, z0)
		ret = append(ret, z1)
	}

	return ret, nil
}

func main() {
	// Инициализируем GTK.
	gtk.Init(nil)

	// Создаём билдер
	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	err = b.AddFromFile("test.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Получаем объект главного окна по ID
	obj, err := b.GetObject("window")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	win := obj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Получаем кнопку
	obj, err = b.GetObject("generateButton")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	generateButton := obj.(*gtk.Button)

	// Получаем поле ввода базы генератора
	obj, err = b.GetObject("baseEntry")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	baseEntry := obj.(*gtk.Entry)

	// Получаем первое поле ввода
	obj, err = b.GetObject("entry1")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	entry1 := obj.(*gtk.Entry)

	// Получаем второе поле ввода
	obj, err = b.GetObject("entry2")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	entry2 := obj.(*gtk.Entry)

	// Получаем поел ввода кол-ва столбцов диаграммы
	obj, err = b.GetObject("entry3")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	groupsEntry := obj.(*gtk.Entry)

	// Получаем поле ввода количетсва дял генерации
	obj, err = b.GetObject("amountEntry")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	amountEntry := obj.(*gtk.Entry)

	// Получаем окно с текстом
	obj, err = b.GetObject("textbox")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	textbox := obj.(*gtk.TextView)
	textbox.SetEditable(false)

	// Получаем лэйбл мат ожидание
	obj, err = b.GetObject("expectation")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	expectation := obj.(*gtk.Label)

	// Получаем лэйбл квадратичное отклонение
	obj, err = b.GetObject("deviation")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	deviation := obj.(*gtk.Label)

	// Получаем окно рисования
	obj, err = b.GetObject("area51")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	area51 := obj.(*gtk.DrawingArea)

	// Получаем переключатель на равномерное распределение
	obj, err = b.GetObject("radio1")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	radio1 := obj.(*gtk.RadioButton)

	// Получаем переключатель на нормальное распределение
	obj, err = b.GetObject("radio2")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	radio2 := obj.(*gtk.RadioButton)

	// Получаем надпись над полем ввода 1
	obj, err = b.GetObject("label1")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	label1 := obj.(*gtk.Label)

	// Получаем надпись над полем ввода 2
	obj, err = b.GetObject("label2")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	label2 := obj.(*gtk.Label)

	// Сигнал при нажатии на первый переключатель
	radio1.Connect("toggled", func() {
		label1.SetText("Нижняя граница")
		label2.SetText("Верхняя граница")
	})

	// Сигнал при нажатии на первый переключатель
	radio2.Connect("clicked", func() {
		label1.SetText("Мат ожидание")
		label2.SetText("Дисперсия")
	})

	var arr []float64

	// Сигнал по нажатию на кнопку
	generateButton.Connect("clicked", func() {
		if radio1.GetActive() {
			// получаем базу генератора
			base, err := baseEntry.GetText()
			if err != nil {
				fmt.Println("error geting base")
				return
			}
			// получаем нижнюю границу
			lower, err := entry1.GetText()
			if err != nil {
				fmt.Println("error geting lower")
				return
			}

			// получаем верхнюю границу
			upper, err := entry2.GetText()
			if err != nil {
				fmt.Println("error geting upper")
				return
			}

			// получаем количество для генерации
			amount, err := amountEntry.GetText()
			if err != nil {
				fmt.Println("error geting amount")
				return
			}

			// получаем указатель на начало текстового поля для вывода
			buffer, err := textbox.GetBuffer()
			if err != nil {
				fmt.Println("error geting textbox buffer")
				return
			}

			arr, err = linearArr(base, lower, upper, amount)
			if err != nil {
				fmt.Println("error making thinhgs")
				return
			}

			expectation.SetText(fmt.Sprintf("Математическое ожидание =%g", mathExpectation(arr)))
			deviation.SetText(fmt.Sprintf("Дисперсия =%g", mathDeviation(arr)))

			start := buffer.GetStartIter()
			end := buffer.GetEndIter()
			buffer.Delete(start, end)
			for _, v := range arr {
				buffer.Insert(start, fmt.Sprintf("%g", v)+"\n")
			}

		}

		if radio2.GetActive() {
			// получаем базу генератора
			base, err := baseEntry.GetText()
			if err != nil {
				fmt.Println("error geting base")
				return
			}

			// получаем мат ожидание
			expected, err := entry1.GetText()
			if err != nil {
				fmt.Println("error geting expected")
				return
			}

			// получаем дисперсию
			variance, err := entry2.GetText()
			if err != nil {
				fmt.Println("error geting variance")
				return
			}

			// получаем количество для генерации
			amount, err := amountEntry.GetText()
			if err != nil {
				fmt.Println("error geting amount")
				return
			}

			// получаем указатель на начало текстового поля для вывода
			buffer, err := textbox.GetBuffer()
			if err != nil {
				fmt.Println("error geting textbox buffer")
				return
			}

			arr, err = normalArray(base, expected, variance, amount)
			if err != nil {
				fmt.Println("error making thinhgs")
				return
			}

			expectation.SetText(fmt.Sprintf("Математическое ожидание =%g", mathExpectation(arr)))
			deviation.SetText(fmt.Sprintf("Дисперсия =%g", mathDeviation(arr)))

			start := buffer.GetStartIter()
			end := buffer.GetEndIter()
			buffer.Delete(start, end)
			for _, v := range arr {
				buffer.Insert(start, fmt.Sprintf("%g", v)+"\n")
			}
		}

		area51.Connect("draw", func(area51 *gtk.DrawingArea, cr *cairo.Context) {
			groups, err := groupsEntry.GetText()
			if err != nil {
				fmt.Println("error geting groups")
				return
			}
			draw2(cr, arr, groups)
		})

	})

	// Отображаем все виджеты в окне
	win.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()
}
