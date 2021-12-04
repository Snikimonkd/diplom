package generator

import (
	"diplom/models"
	"fmt"
	"math"

	"fyne.io/fyne/v2"
)

var seed int32

const a = 1103515245
const c = 12345

func lcgInt() int32 {
	seed = (a*seed + c)
	return seed
}

func intToFloat(randInt int32) float32 {
	return float32(randInt) / float32(math.MaxInt32)
}

func Generate(lower, upper, amount, base int) []float32 {
	var ret []float32
	seed = int32(base)
	for i := 0; i < amount; i++ {
		randomInt := lcgInt()
		randomFloat := intToFloat(randomInt)

		ret = append(ret, float32(math.Abs(float64(randomFloat))*float64(upper-lower)+float64(lower)))
	}

	return ret
}

func partition(arr []float32, low, high int) ([]float32, int) {
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

func quickSort(arr []float32, low, high int) []float32 {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func distribution(arr []float32, cols int) []int {
	sortArr := quickSort(arr, 0, len(arr)-1)

	min := sortArr[0]
	max := sortArr[len(sortArr)-1]

	groupSize := (max - min) / float32(cols)

	var distr []int
	distr = append(distr, 0)
	i := 0

	for _, v := range sortArr {
		if v < (min + float32(i+1)*groupSize) {
			distr[i]++
		} else {
			i++
			distr = append(distr, 0)
			distr[i]++
		}
	}

	return distr[:len(distr)-1]
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

func Draw(arr []float32, cols int) []fyne.CanvasObject {
	distr := distribution(arr, cols)

	max := findMax(distr)

	width := int(models.WindowWidth-30) / cols
	height := 700 / max

	ret := []fyne.CanvasObject{}

	for i, v := range distr {
		pos := fyne.NewPos(float32(i*width), float32(750-height*v))
		size := fyne.NewSize(float32(width), float32(height*v))
		ret = append(ret, models.CreateRectangel(size, pos))
	}

	fmt.Println(ret)

	return ret
}
