package rs

import "sort"

// import (
// 	"fmt"
// 	"math"
// )

// RStats hold the moments
type RStats struct {
	rMean, rVar, rStdDev, rSkew, rKurt float64
	M1, M2, M3, M4                     float64
	N                                  uint64
	Min, Max                           float64
}

// ResetRoll zeroises
func ResetRoll(rs *RStats) {
rs.M1 = 0
rs.M2 = 0
rs.M3 = 0
rs.M4 = 0
rs.N = 0
rs.Min = 0
rs.Max = 0
}
// RQuant holds the running quantizatoin data
type RQuant struct {
	p          float64
	q, np, dnp [5]float64
	npos       [5]int64
	n          int64
}

// Reinit resets the RQunat structure for kicking off new quantization
func Reinit(rQ *RQuant, p float64) {
	// initialize positions n
	for i := 0; i < 5; i++ {
		rQ.npos[i] = int64(i) + 1
	}
	// initialize np
	rQ.np[0] = 1.0
	rQ.np[1] = 1.0 + 2.0*p
	rQ.np[2] = 1.0 + 4.0*p
	rQ.np[3] = 3.0 + 2.0*p
	rQ.np[4] = 5.0
	// initialize dn
	rQ.dnp[0] = 0.0
	rQ.dnp[1] = 0.5 * p
	rQ.dnp[2] = p
	rQ.dnp[3] = 0.5 * (1.0 + p)
	rQ.dnp[4] = 1.0
	//
	rQ.n = 0
}

//QuantRoller calculates a rolling qunatile
func QuantRoller(x float64, rQ *RQuant) {
	if rQ.n < 5 {
		rQ.q[rQ.n] = x
	} else {
		var i, k int64 = 0, -1
		if rQ.n == 5 {
			// sort first five heights
			sort.Float64s(rQ.q[:])
		}
		/* step B1: find k such that q_k <= x < q_{k+1} */
		if x < rQ.q[0] {
			rQ.q[0] = x
			k = 0
		} else if x >= rQ.q[4] {
			rQ.q[4] = x
			k = 3
		} else {
			for i = 0; i <= 3; i++ {
				if (rQ.q[i] <= x) && (x < rQ.q[i+1]) {
					k = i
					break
				}
			}
		}
		if k < 0 {
			/* we could get here if x is nan */
			//GSL_ERROR ("invalid input argument x", GSL_EINVAL);
		}

		/* step B2(a): update n_i */
		for i = k + 1; i <= 4; i++ {
			(rQ.npos[i])++
		}

		/* step B2(b): update n_i' */
		for i = 0; i < 5; i++ {
			rQ.np[i] += rQ.dnp[i]
		}

		/* step B3: update heights */
		for i = 1; i <= 3; i++ {
			ni := rQ.npos[i]
			d := rQ.np[i] - float64(ni)

			if ((d >= 1.0) && (rQ.npos[i+1]-rQ.npos[i] > 1)) || ((d <= -1.0) && (rQ.npos[i-1]-rQ.npos[i] < -1)) {
				var dsign int64
				if d > 0.0 {
					dsign = 1
				} else {
					dsign = -1
				}
				qp1 := rQ.q[i+1]
				qi := rQ.q[i]
				qm1 := rQ.q[i-1]
				np1 := rQ.npos[i+1]
				nm1 := rQ.npos[i-1]
				//qp := calc_psq(qp1, qi, qm1, dsign, np1, ni, nm1)
				outer := float64(dsign) / float64(np1-nm1)
				innerleft := (float64(ni) - float64(nm1) + float64(dsign)) * (qp1 - qi) / float64(np1-ni)
				innerright := (float64(np1-ni) - float64(dsign)) * (qi - qm1) / float64(ni-nm1)
				qp := qi + outer*(innerleft+innerright)
				//
				
				if qm1 < qp && qp < qp1 {
					rQ.q[i] = qp
				} else {
						/* use linear formula */
					rQ.q[i] = rQ.q[i] + float64(dsign)*((rQ.q[i+dsign]-qi)/(float64(rQ.npos[i+dsign])-float64(ni)))
				}

				rQ.npos[i] += dsign
				
			}
		}

	}
	(rQ.n)++

}

// RQuantResult
func RQuantResult(rQ *RQuant) float64 {
	if rQ.n >= 5 {
		return rQ.q[2]
	} else {
		/* not yet initialized */
		//gsl_sort(w->q, 1, w->n);
		sort.Float64s(rQ.q[:])
		//return gsl_stats_quantile_from_sorted_data(w->q, 1, w->n, w->p);
		//                                           rQ.q, 1, rQ.n, rQ.p
		// FUNCTION(gsl_stats,quantile_from_sorted_data)
		//(const BASE sorted_data[], rQ.q
		//const size_t stride,       1
		//const size_t n,            rQ.n
		//const double f)            rQ.p
		index := rQ.p * float64(rQ.n-1)
		lhs := int64(index)
		delta := float64(index - float64(lhs))
		var result float64

		if rQ.n == 0 {
			return 0.0
		}

		if lhs == (rQ.n - 1) {
			result = rQ.q[lhs*1]
		} else {
			result = (1-delta)*rQ.q[lhs*1] + delta*rQ.q[(lhs+1)*1]
		}
		return result
	}

}

// RollingStat calulates the  moments for a new data in
func RollingStat(x float64, rS *RStats) {

	var delta float64
	var delta_n float64
	var delta_n2 float64
	var term1 float64
	var n1 float64

	/* update Min and Max */
	if rS.N == 0 {
		rS.Min = x
		rS.Max = x
	} else {
		if x < rS.Min {
			rS.Min = x
		}
		if x > rS.Max {
			rS.Max = x
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

	// update median

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
