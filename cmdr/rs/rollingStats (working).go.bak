package rs

// import (
// 	"fmt"
// 	"math"
// )

type RStats struct {
	rMean, rVar, rStdDev, rSkew, rKurt float64
	M1, M2, M3, M4                     float64
	N                                  uint64
	min, max						float64
}

func RollingStat(x float64, rS *RStats) {

	var delta float64
	var delta_n float64
	var delta_n2 float64
	var term1 float64
	var n1 float64

	  /* update min and max */
	  if (rS.N == 0) {
		rS.min = x;
		rS.max = x;
	  } else {
		if (x < rS.min) {
		  rS.min = x;
		}
		if (x > rS.max) {
		  rS.max = x;
		}
	  }
  

  /* update mean and variance */
	n1 = float64(rS.N)
	rS.N = rS.N + 1
	delta = x - rS.M1
	delta_n = delta / float64(rS.N)
	delta_n2 = delta_n * delta_n
	term1 = delta * delta_n * n1
	n1 = n1 + 1
	rS.M1 = rS.M1 + delta_n
	rS.M4 = rS.M4 + term1*delta_n2*(n1*n1-3*n1+3) + 6*delta_n2*rS.M2 - 4*delta_n*rS.M3
	rS.M3 = rS.M3 + term1*delta_n*(n1-2) - 3*delta_n*rS.M2
	rS.M2 = rS.M2 + term1
}

//n = n + 1
//adelta = x - mean
//delta_n = adelta / n
//delta_nsq = delta_n * delta_n
//term1 = adelta * delta_n * (n - 1)
//mean = mean + delta_n
//M4 = M4 + term1 * delta_nsq * (n * n - 3 * n + 3) + 6 * delta_nsq * M2 - 4 * delta_n * M3
//M3 = M3 + term1 * delta_n * (n - 2) - 3 * delta_n * M2
//M2 = M2 + term1

// func main() {
// 	rS := new(rStats)
// 	rollingStat(17.2, rS)
// 	fmt.Println("Mean :", rS.M1, "Std Dev: ", math.Sqrt(rS.M2/((float64(rS.n))-1.0)))
// 	rollingStat(18.1, rS)
// 	fmt.Println("Mean :", rS.M1, "Std Dev: ", math.Sqrt(rS.M2/((float64(rS.n))-1.0)))
// 	rollingStat(16.5, rS)
// 	fmt.Println("Mean :", rS.M1, "Std Dev: ", math.Sqrt(rS.M2/((float64(rS.n))-1.0)))
// 	rollingStat(18.3, rS)
// 	fmt.Println("Mean :", rS.M1, "Std Dev: ", math.Sqrt(rS.M2/((float64(rS.n))-1.0)))
// 	rollingStat(12.6, rS)
// 	fmt.Println("Mean :", rS.M1, "Std Dev: ", math.Sqrt(rS.M2/((float64(rS.n))-1.0)))

// }
