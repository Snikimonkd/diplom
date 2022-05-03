package generator

import (
	"diplom/internal/models"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

var seed int32

const a = 1664525
const c = 1013904223

type OneProductModel struct {
	// полные издержки системы
	TC float32
	// полные издержки, связанные с организацией поставок
	TC1 float32
	// полные издержки, связанные с организацией поставок
	TC2 float32
	// полные потери от дефицита продукта на складе
	TC3 float32
	// текущее время
	CLOCK int32
	// время очередной поставки
	T int32
	// количество запаса на складе
	V1 int32
	// спрос в i-ый день
	D int32
	// время необходимое для выполнения j-ого заказа
	PLT int32
	// объем одной поставки
	EOQ int32
	// точка возобновления запаса
	ROP int32
	// затраты на хранение единицы продукта в течение одного дня
	C1 float32
	// затраты на организацию одной поставки
	C2 float32
	// потери, связанные с нехваткой единицы продукта
	C3 float32
	// начальный уровень запаса
	B1 int32
	// продолжительность рассматриваемого периода
	TT int32

	// номер эксперимента
	ExpNumber int
}

type MultiProductModel struct {
	// затраты на организацию поставки
	TOC float32
	// затраты на хранение запасов
	TCC float32
	// полные затраты
	TCOST float32
	// объем заказа i-го продукта
	EOQ []int32
	// критический уровень запаса i-го продукта
	MOP []int32
	// предкритический уровень запаса i-го продукта
	COP []int32
	// спрос на i-ый продукт в t-ый день
	D int32 // должно быть по закону Пуассона
	// уровень запаса i-го продукта в конце t-го дня
	INV []int32
	// кол-во заказов на i-ый продукт в течение времени T
	NTO []int32
	// общее число заказов в течение времени
	TNJO int32
	// затраты на оформление одного набора заказа
	FOC int32
	// переменные затраты на заказ i-го продукта
	VOC []int32
	// ежедневные затраты на хранение единицы i-го продукта
	CC []int32

	// количество продуктов
	PRAM int32
	// текущее время
	T int32
	// номер эксперимента
	ExpNumber int32
}

func lcgInt() int32 {
	seed = (a*seed + c)
	return seed
}

func intToFloat(randInt int32) float32 {
	return float32(randInt) / float32(math.MaxInt32)
}

func abs(val int32) int32 {
	if val < 0 {
		return val * -1
	}

	return val
}

func randomIntWithBorders(lower, upper int32) int32 {
	return abs(lcgInt())%(upper-lower) + lower
}

func randomFloatWithBorders(lower, upper int32) float32 {
	return float32(math.Abs(float64(intToFloat(lcgInt())))*float64(upper-lower) + float64(lower))
}

func LinearGenerate(lower, upper, amount, base int32) []float32 {
	var ret []float32
	seed = base
	for i := 0; i < int(amount); i++ {
		randomInt := lcgInt()
		randomFloat := intToFloat(randomInt)

		ret = append(ret, float32(math.Abs(float64(randomFloat))*float64(upper-lower)+float64(lower)))
	}

	return ret
}

func normalPair() (float32, float32, float32) {
	var s, x, y float32
	s = 0
	for s == 0 || s > 1 {
		x = intToFloat(lcgInt())
		y = intToFloat(lcgInt())
		s = x*x + y*y
	}

	return s, x, y
}

func NormalGenerate(mathExpectation, dispersion, amount, base int32) []float32 {
	var ret []float32
	seed = base
	for i := 0; i < int(amount/2); i++ {
		s, x, y := normalPair()

		z0 := x * float32(math.Sqrt(float64(-2*float32(math.Log(float64(s)))/s)))
		z1 := y * float32(math.Sqrt(float64(-2*float32(math.Log(float64(s)))/s)))

		ret = append(ret, (float32(mathExpectation) + float32(dispersion)*z0), (float32(mathExpectation) + float32(dispersion)*z1))
	}

	return ret
}

func ExpGenerate(lyambda, amount, base int32) []float32 {
	var ret []float32
	seed = base
	for i := 0; i < int(amount); i++ {
		randomInt := lcgInt()
		randomFloat := math.Abs(float64(intToFloat(randomInt)))

		res := float64(-1) / float64(lyambda) * math.Log(randomFloat)

		ret = append(ret, float32(res))
	}

	return ret
}

func partition(arr []float32, low, high int32) ([]float32, int32) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func quickSort(arr []float32, low, high int32) []float32 {
	if low < high {
		var p int32
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func distribution(arr []float32, cols int32) []int32 {
	sortArr := quickSort(arr, 0, int32(len(arr)-1))

	min := sortArr[0]
	max := sortArr[len(sortArr)-1]

	groupSize := (max - min) / float32(cols)

	distr := make([]int32, cols)
	distr = append(distr, 0)
	i := 0

	current := min + groupSize

	for _, v := range sortArr {
		if v <= current {
			distr[i]++
		} else {
			i++
			current += groupSize
			distr[i]++
		}
	}

	return distr
}

func findMax(arr []int32) int32 {
	var max int32
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return max
}

func Draw(arr []float32, cols int32) []fyne.CanvasObject {
	distr := distribution(arr, cols)

	max := findMax(distr)

	width := int32(models.WindowWidth-30)/cols - 1
	height := 700 / max

	ret := []fyne.CanvasObject{}

	for i, v := range distr {
		pos := fyne.NewPos(float32(int32(i)*width)+models.DefaultHorizontalPadding+float32(i), float32(750-height*v))
		size := fyne.NewSize(float32(width), float32(height*v))
		ret = append(ret, models.CreateRectangel(size, pos, color.Black))
	}

	return ret
}

func bernoulli(p float32) float32 {
	q := 1.0 - p
	U := float32(math.Abs(float64(intToFloat(lcgInt()))))
	if U > q {
		return 1
	}

	return 0
}

func binomialBernoulli(p float32, n int32) float32 {
	var sum float32 = 0
	for i := 0; i < int(n); i++ {
		sum += bernoulli(p)
	}

	return sum
}

func Binomial(p float32, n, amount, base int32) []float32 {
	seed = base
	ret := []float32{}
	for i := 0; i < int(amount); i++ {
		ret = append(ret, binomialBernoulli(p, n))
	}

	return ret
}

func factorail(n int) int {
	ret := 1
	for i := 1; i <= n; i++ {
		ret *= i
	}

	return ret
}

func Combination(n, k int) int {
	return factorail(n) / (factorail(k) * factorail(n-k))
}

func Ver(p float64, n, k int) float64 {
	return float64(Combination(n, k)) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
}

func generateOneProductModel() OneProductModel {
	var model OneProductModel
	// начальное положение системы
	model.CLOCK = 0
	model.T = 0
	model.TC1 = 0
	model.TC2 = 0
	model.TC3 = 0
	model.TC = 0

	// данные, которые мы перебираем с целью поиска наилучшего варианта
	model.EOQ = 0
	model.ROP = 0

	// данные, зависящие от варианта
	model.C1 = randomFloatWithBorders(0, 10)
	model.C2 = randomFloatWithBorders(0, 100)
	model.C3 = randomFloatWithBorders(0, 8)

	model.TT = randomIntWithBorders(15, 60)
	model.B1 = randomIntWithBorders(100, 200)

	model.V1 = model.B1

	return model
}

func ModelingOneProduct(variant, amount int32) []OneProductModel {
	seed = variant
	var ret []OneProductModel

	for i := 0; i < int(amount); i++ {
		model := generateOneProductModel()
		model.EOQ = randomIntWithBorders(20, 50)
		model.ROP = randomIntWithBorders(40, 70)
		model.ExpNumber = i + 1
		for {
			model.D = randomIntWithBorders(0, 10)

			model.CLOCK++

			if model.CLOCK <= model.TT {
				if model.T == model.CLOCK {
					model.V1 += model.EOQ
				}
			} else {
				model.TC = model.TC1 + model.TC2 + model.TC3
				break
			}
			model.V1 -= model.D

			if model.V1 < 0 {
				model.TC3 -= float32(model.V1) * model.C3
				model.V1 = 0
			}

			model.TC1 += float32(model.V1) * model.C1

			if model.ROP >= model.V1 {
				if model.T <= model.CLOCK {
					model.TC2 += model.C2
					model.PLT = randomIntWithBorders(0, 10)
					model.T += model.CLOCK + model.PLT
				}
			}
		}

		ret = append(ret, model)
	}

	return ret
}

func SortModelsOne(arr []OneProductModel, id int) []OneProductModel {
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			switch id {
			case 0:
				{
					if arr[j].ExpNumber < arr[i].ExpNumber {
						arr[j], arr[i] = arr[i], arr[j]
					}
				}
			case 1:
				{
					if arr[j].ROP < arr[i].ROP {
						arr[j], arr[i] = arr[i], arr[j]
					}
				}
			case 2:
				{
					if arr[j].EOQ < arr[i].EOQ {
						arr[j], arr[i] = arr[i], arr[j]
					}
				}
			case 3:
				{
					if arr[j].TC1 < arr[i].TC1 {
						arr[j], arr[i] = arr[i], arr[j]
					}
				}
			case 4:
				{
					if arr[j].TC2 < arr[i].TC2 {
						arr[j], arr[i] = arr[i], arr[j]
					}
				}
			case 5:
				{
					if arr[j].TC3 < arr[i].TC3 {
						arr[j], arr[i] = arr[i], arr[j]
					}
				}
			case 6:
				{
					if arr[j].TC < arr[i].TC {
						arr[j], arr[i] = arr[i], arr[j]
					}
				}
			}
		}
	}

	return arr
}

func ExpValue(arr []float32) float32 {
	var res float32 = 0

	len := float32(len(arr))

	for _, v := range arr {
		res += v / len
	}

	return res
}

func DispValue(arr []float32) float32 {
	var exp float32 = 0.0
	var buf float32 = 0.0

	len := float32(len(arr))

	for _, v := range arr {
		buf = buf + v*v/len
		exp += v / len
	}

	return float32(math.Sqrt(float64(buf - exp*exp)))
}

func Count(arr []float32, level int32) []int32 {
	res := make([]int32, level+1)

	for _, v := range arr {
		res[int(v)]++
	}

	return res
}

func randArr(len, left, right int32) []int32 {
	var arr []int32

	for i := 0; i < int(len); i++ {
		arr = append(arr, randomIntWithBorders(left, right))
	}

	return arr
}

func generateMultiProductModel() MultiProductModel {
	var model MultiProductModel

	// начальное положение системы
	model.T = 1
	model.TOC = 0
	model.TCC = 0
	model.TCOST = 0
	model.TNJO = 0

	// данные зависящие от варианта
	model.PRAM = randomIntWithBorders(1, 5)
	model.CC = randArr(model.PRAM, 0, 100)
	model.VOC = randArr(model.PRAM, 0, 20)
	model.FOC = randomIntWithBorders(0, 20)

	return model
}

func ModelingMultiProduct(variant, amountExp int32) []MultiProductModel {
	seed = variant
	var ret []MultiProductModel

	model := generateMultiProductModel()

	for k := 0; k < int(amountExp); k++ {

		model.EOQ = randArr(model.PRAM, 50, 100)
		model.MOP = randArr(model.PRAM, 0, 10)
		model.COP = randArr(model.PRAM, 10, 20)

		model.ExpNumber = int32(k)

		model.NTO = make([]int32, model.PRAM)
		model.INV = make([]int32, model.PRAM)

		for model.T = 0; model.T < 90; model.T++ {
			for i := 0; i < int(model.PRAM); i++ {
				model.D = randomIntWithBorders(0, 10)
				model.INV[i] -= model.D
			}

			for i := 0; i < int(model.PRAM); i++ {
				if model.INV[i]-model.MOP[i] <= 0 {
					model.TNJO++
					model.TOC += float32(model.FOC)
					for j := 0; j < int(model.PRAM); j++ {
						if model.INV[j]-model.COP[j] <= 0 {
							model.NTO[j]++
							model.TOC += float32(model.VOC[j])
							model.INV[j] += model.EOQ[j]
							if model.INV[j] > 0 {
								model.TCC += float32(model.CC[j] * model.INV[j])
							}
						}
					}
				} else {
					for j := 0; j < int(model.PRAM); j++ {
						if model.INV[j] > 0 {
							model.TCC += float32(model.CC[j] * model.INV[j])
						}
					}
				}
			}

		}
		model.TCOST = model.TOC + model.TCC

		ret = append(ret, model)
	}

	return ret
}
