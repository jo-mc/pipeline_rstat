package cmdr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	//"strings"
	"math"
	"strconv"

	"github.com/pkg/errors"

	"redo_pipeline_rstat/cmdr/rs"
)

// Flags for use
var Flags struct {
	Filepath string
	Verbose  bool
}

// RunCommand runs the command
func RunCommand() error {
	if isInputFromPipe() {
		print("data is from pipe")
		return toUppercase(os.Stdin, os.Stdout)
	}
	 {
		file, e := getFile()
		if e != nil {
			return e
		}
		defer file.Close()
		return toUppercase(file, os.Stdout)
	}
}

func isInputFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}

func getFile() (*os.File, error) {
	if Flags.Filepath == "" {
		return nil, errors.New("please input a file")
	}
	if !fileExists(Flags.Filepath) {
		return nil, errors.New("the file provided does not exist")
	}
	file, e := os.Open(Flags.Filepath)
	if e != nil {
		return nil, errors.Wrapf(e,
			"unable to read the file %s", Flags.Filepath)
	}
	return file, nil
}

func toUppercase(r io.Reader, w io.Writer) error {
	rS := new(rs.RStats)
	r1 := new(rs.RStats)
	scanner := bufio.NewScanner(bufio.NewReader(r))
	var anum float64
	var countN uint64
	var statblock int = 5
	anum = anum + 0.0  // fudge says anum not used?
	for scanner.Scan() {
		//_, e := fmt.Fprintln(
			//w, strings.ToUpper(scanner.Text()))
			
			anum, err := strconv.ParseFloat(scanner.Text(), 64)
			if err == nil {
				rs.RollingStat(anum, rS)
				rs.RollingStat(anum, r1)
				} else {
				fmt.Println("ERROR, not numeric:", err)
			}
			countN = countN + 1
			if math.Mod(float64(countN),float64(statblock)) == 0 {
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>block :",countN - uint64(statblock), " to  ", countN)
				fmt.Println("The sample Mean :", r1.M1)
				fmt.Println("The estimated variance :", (r1.M2/(float64(r1.N) - 1)))
				fmt.Println("Largest Value :",r1.Max)
				fmt.Println("Smallest Value :",r1.Min)
				fmt.Println("Median : NA")
				fmt.Println("Standard Deviation : ", (math.Sqrt(r1.M2/((float64(r1.N))-1.0)))/math.Sqrt(float64(r1.N)))
				if r1.N > 0 {
				fmt.Println("Skew :",((math.Pow(float64(r1.N)-1.0, 1.5)/float64(r1.N))*r1.M3) / (math.Pow(r1.M2, 1.5)))
				} else {
					fmt.Println("Skew : 0")
				}
				if r1.N > 0 {
				fmt.Println("Kurtosis :",(((float64(r1.N)-1.0)/float64(r1.N))*(float64(r1.N) - 1.0))*r1.M4 / ((r1.M2 * r1.M2)) - 3.0) 
				} else {
					fmt.Println("Kurtosis : 0")
				}
				fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")				
				fmt.Println()
				//  and reset
				rs.ResetRoll(r1)
			}

		// if e != nil {
		// 	return e
		// }
	}

	fmt.Println("")
	fmt.Println("Overall results:")
	fmt.Println("The sample Mean :", rS.M1)
	fmt.Println("Std Dev: ", math.Sqrt(rS.M2/((float64(rS.N))-1.0)))
	fmt.Println("The estimated variance :", (rS.M2/(float64(rS.N) - 1)))
	fmt.Println("Largest Value :",rS.Max)
	fmt.Println("Smallest Value :",rS.Min)
	fmt.Println("Median : NA")
	fmt.Println("Standard Deviation of the mean : ",(math.Sqrt(rS.M2/((float64(rS.N))-1.0)))/math.Sqrt(float64(rS.N)))
	if rS.N > 0 {
	fmt.Println("Skew :",((math.Pow(float64(rS.N)-1.0, 1.5)/float64(rS.N))*rS.M3) / (math.Pow(rS.M2, 1.5)))
	} else {
		fmt.Println("Skew : 0")
	}
	if rS.N > 0 {
	fmt.Println("Kurtosis :",(((float64(rS.N)-1.0)/float64(rS.N))*(float64(rS.N) - 1.0))*rS.M4 / ((rS.M2 * rS.M2)) - 3.0) 
	} else {
		fmt.Println("Kurtosis : 0")
	}
	fmt.Println("Number of items :",rS.N)
	
	return nil
}

func fileExists(filepath string) bool {
	info, e := os.Stat(filepath)
	if os.IsNotExist(e) {
		return false
	}
	return !info.IsDir()
}
