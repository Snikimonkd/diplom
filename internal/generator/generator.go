package generator

import (
	"diplom/internal/models"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

var seed int

const a = 6364136223846793005
const c = 1442695040888963407

type Model struct {
	// полные издержки системы
	TC float64
	// полные издержки, связанные с организацией поставок
	TC1 float64
	// полные издержки, связанные с организацией поставок
	TC2 float64
	// полные потери от дефицита продукта на складе
	TC3 float64
	// текущее время
	CLOCK int
	// время очередной поставки
	T int
	// количество запаса на складе
	V1 int
	// спрос в i-ый день
	D int
	// время необходимое для выполнения j-ого заказа
	PLT int
	// объем одной поставки
	EOQ int
	// точка возобновления запаса
	ROP int
	// затраты на хранение единицы продукта в течение одного дня
	C1 float64
	// затраты на организацию одной поставки
	C2 float64
	// потери, связанные с нехваткой единицы продукта
	C3 float64
	// начальный уровень запаса
	B1 int
	// продолжительность рассматриваемого периода
	TT int
	// номер эксперимента
	ExpNumber int
}

func lcgInt() int {
	seed = (a*seed + c)
	return seed
}

func intToFloat(randInt int) float64 {
	return float64(randInt) / float64(math.MaxInt)
}

func abs(val int) int {
	if val < 0 {
		return val * -1
	}

	return val
}

func randomIntWithBorders(lower, upper int) int {
	return abs(lcgInt())%(upper-lower) + lower
}

func randomFloatWithBorders(lower, upper int) float64 {
	return math.Abs(intToFloat(lcgInt()))*float64(upper-lower) + float64(lower)
}

func LinearGenerate(lower, upper, amount, base int) []float64 {
	var ret []float64
	seed = base
	for i := 0; i < amount; i++ {
		randomInt := lcgInt()
		randomFloat := intToFloat(randomInt)

		ret = append(ret, math.Abs(randomFloat)*float64(upper-lower)+float64(lower))
	}

	return ret
}

func normalPair() (float64, float64, float64) {
	var s, x, y float64
	s = 0
	for s == 0 || s > 1 {
		x = intToFloat(lcgInt())
		y = intToFloat(lcgInt())
		s = x*x + y*y
	}

	return s, x, y
}

func NormalGenerate(mathExpectation, dispersion, amount, base int) []float64 {
	var ret []float64
	seed = base
	for i := 0; i < amount/2; i++ {
		s, x, y := normalPair()

		z0 := x * math.Sqrt(-2*math.Log(s)/s)
		z1 := y * math.Sqrt(-2*math.Log(s)/s)

		ret = append(ret, (float64(mathExpectation) + float64(dispersion)*z0), (float64(mathExpectation) + float64(dispersion)*z1))
	}

	return ret
}

func ExpGenerate(lyambda, amount, base int) []float64 {
	var ret []float64
	seed = base
	for i := 0; i < amount; i++ {
		randomInt := lcgInt()
		randomFloat := math.Abs(intToFloat(randomInt))

		res := float64(-1) / float64(lyambda) * math.Log(1-randomFloat)

		ret = append(ret, res)
	}

	return ret
}

func partition(arr []float64, low, high int) ([]float64, int) {
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

func quickSort(arr []float64, low, high int) []float64 {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func distribution(arr []float64, cols int) []int {
	sortArr := quickSort(arr, 0, len(arr)-1)

	min := sortArr[0]
	max := sortArr[len(sortArr)-1]

	groupSize := (max - min) / float64(cols)

	var distr []int
	distr = append(distr, 0)
	i := 0

	for _, v := range sortArr {
		if v <= (min + float64(i+1)*groupSize) {
			distr[i]++
		} else {
			i++
			distr = append(distr, 0)
			distr[i]++
		}
	}

	return distr
}

func findMax(arr []int) int {
	max := 0
	for _, v := range arr {
		if v > max {
			max = v
		}
	}

	return max
}

func Draw(arr []float64, cols int) []fyne.CanvasObject {
	distr := distribution(arr, cols)

	max := findMax(distr)

	width := int(models.WindowWidth-30) / cols
	height := 700 / max

	ret := []fyne.CanvasObject{}

	for i, v := range distr {
		pos := fyne.NewPos(float32(i*width)+models.DefaultHorizontalPadding, float32(750-height*v))
		size := fyne.NewSize(float32(width), float32(height*v))
		ret = append(ret, models.CreateRectangel(size, pos, color.Black))
	}

	return ret
}

func bernoulli(p float64) float64 {
	q := 1.0 - p
	U := float64(math.Abs(float64(intToFloat(lcgInt()))))
	if U > q {
		return 1
	}

	return 0
}

func binomialBernoulli(p float64, n int) float64 {
	var sum float64 = 0
	for i := 0; i < n; i++ {
		sum += bernoulli(p)
	}

	return sum
}

func Binomial(p float64, n, amount, base int) []float64 {
	seed = base
	ret := []float64{}
	for i := 0; i < amount; i++ {
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

func generateModel() Model {
	var model Model
	// начальное положение системы
	model.CLOCK = 0
	model.T = 0
	model.TC1 = 0
	model.TC2 = 0
	model.TC3 = 0
	model.TC = 0

	// данные, которые мы перебираем с целью поиска наилучшего варианта
	model.EOQ = randomIntWithBorders(0, 10)
	model.ROP = randomIntWithBorders(0, 10)

	// данные, зависящие от варианта
	model.C1 = randomFloatWithBorders(0, 10)
	model.C2 = randomFloatWithBorders(0, 100)
	model.C3 = randomFloatWithBorders(0, 8)

	model.TT = randomIntWithBorders(15, 60)

	model.B1 = randomIntWithBorders(100, 200)
	model.V1 = model.B1

	return model
}

func Modeling(variant, amount int) []Model {
	seed = variant
	var ret []Model

	for i := 0; i < amount; i++ {
		model := generateModel()
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
				model.TC3 -= float64(model.V1) * model.C3
				model.V1 = 0
			}

			model.TC1 += float64(model.V1) * model.C1

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

func SortModels(arr []Model, id int) []Model {
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