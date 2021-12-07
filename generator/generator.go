package generator

import (
	"diplom/models"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
)

var seed int

const a = 6364136223846793005
const c = 1442695040888963407

func lcgInt() int {
	seed = (a*seed + c)
	return seed
}

func intToFloat(randInt int) float64 {
	return float64(randInt) / float64(math.MaxInt)
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
	//ret = append(ret, models.CreateRectangel(fyne.NewSize(800, 800), fyne.NewPos(0, 0), color.White))

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
