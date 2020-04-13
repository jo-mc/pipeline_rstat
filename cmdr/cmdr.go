package cmdr

import (
	"bufio"
	"fmt"
	"io"
	"os"

	//"strings"
	"math"
	"strconv"
	"strings"

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
	var region string
	var lastregion string
	
	var statblock int64 = 1000000
	anum = anum + 0.0 // fudge says anum not used?
	for scanner.Scan() {
		//_, e := fmt.Fprintln(
		//w, strings.ToUpper(scanner.Text()))

		line := scanner.Text()
		s := strings.Fields(line)
		region = s[0]
		if region != lastregion {
			countN = 0
		}
		lastregion = region
		anum, err := strconv.ParseFloat(s[2], 64)
		if err == nil {
			rs.RollingStat(anum, rS)
			rs.RollingStat(anum, r1)
		} else {
			fmt.Println("ERROR, not numeric:", err)
		}
		countN = countN + 1
		if math.Mod(float64(countN), float64(statblock)) == 0 {
			fmt.Print(">:,", region, ":", countN-uint64(statblock), ",to,", countN, ",")
			fmt.Printf("Mean:,%.2f", r1.M1)
			fmt.Printf(",Var:,%.2f", (r1.M2 / (float64(r1.N) - 1)))
			fmt.Printf(",Max:,%.2f", r1.Max)
			fmt.Printf(",Min:,%.2f", r1.Min)
			fmt.Printf(",Med:,NA")
			fmt.Printf(",StdDev:,%.2f", (math.Sqrt(r1.M2/((float64(r1.N))-1.0)))/math.Sqrt(float64(r1.N)))
			if r1.N > 0 {
				fmt.Printf(",Skew:,%.2f", ((math.Pow(float64(r1.N)-1.0, 1.5)/float64(r1.N))*r1.M3)/(math.Pow(r1.M2, 1.5)))
			} else {
				fmt.Println(",Skew:,0,")
			}
			if r1.N > 0 {
				fmt.Printf(",Kurtosis:,%.2f", (((float64(r1.N)-1.0)/float64(r1.N))*(float64(r1.N)-1.0))*r1.M4/(r1.M2*r1.M2)-3.0)
			} else {
				fmt.Println(",Kurtosis:,0")
			}
			fmt.Println(",<")
			//  and reset
			rs.ResetRoll(r1)
		}

		// if e != nil {
		// 	return e
		// }
	}

	fmt.Println("")
	fmt.Println("Overall results:")
	fmt.Printf("The sample Mean :%.2f\n", rS.M1)
	fmt.Printf("Std Dev: %.2f\n", math.Sqrt(rS.M2/((float64(rS.N))-1.0)))
	fmt.Printf("The estimated variance :%.2f\n", (rS.M2 / (float64(rS.N) - 1)))
	fmt.Printf("Largest Value :%.2f\n", rS.Max)
	fmt.Printf("Smallest Value :%.2f\n", rS.Min)
	fmt.Print("Median : NA\n")
	fmt.Printf("Standard Deviation of the mean : %.2f\n", (math.Sqrt(rS.M2/((float64(rS.N))-1.0)))/math.Sqrt(float64(rS.N)))
	if rS.N > 0 {
		fmt.Printf("Skew :%.2f\n", ((math.Pow(float64(rS.N)-1.0, 1.5)/float64(rS.N))*rS.M3)/(math.Pow(rS.M2, 1.5)))
	} else {
		fmt.Println("Skew : 0")
	}
	if rS.N > 0 {
		fmt.Printf("Kurtosis :%.2f\n", (((float64(rS.N)-1.0)/float64(rS.N))*(float64(rS.N)-1.0))*rS.M4/(rS.M2*rS.M2)-3.0)
	} else {
		fmt.Println("Kurtosis : 0")
	}
	fmt.Println("Number of items :", rS.N)

	return nil
}

func fileExists(filepath string) bool {
	info, e := os.Stat(filepath)
	if os.IsNotExist(e) {
		return false
	}
	return !info.IsDir()
}
