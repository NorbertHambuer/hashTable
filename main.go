package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var judMap = map[string]string{
	"Alba":                 "01",
	"Arad":                 "02",
	"Arges":                "03",
	"Bacau":                "04",
	"Bihor":                "05",
	"Bistrita-Nasaud":      "06",
	"Botosani":             "07",
	"Brasov":               "08",
	"Braila":               "09",
	"Buzau":                "10",
	"Caras-Severin":        "11",
	"Cluj":                 "12",
	"Constanta":            "13",
	"Covasna":              "14",
	"Dambovita":            "15",
	"Dolj":                 "16",
	"Galati":               "17",
	"Gorj":                 "18",
	"Harghita":             "19",
	"Hunedoara":            "20",
	"Ialomita":             "21",
	"Iasi":                 "22",
	"Ilfov":                "23",
	"Maramures":            "24",
	"Mehedinti":            "25",
	"Mures":                "26",
	"Neamt":                "27",
	"Olt":                  "28",
	"Prahova":              "29",
	"Satu Mare":            "30",
	"Salaj":                "31",
	"Sibiu":                "32",
	"Suceava":              "33",
	"Teleorman":            "34",
	"Timis":                "35",
	"Tulcea":               "36",
	"Vaslui":               "37",
	"Valcea":               "38",
	"Vrancea":              "39",
	"Bucuresti":            "40",
	"Bucuresti Sectorul 1": "41",
	"Bucuresti Sectorul 2": "42",
	"Bucuresti Sectorul 3": "43",
	"Bucuresti Sectorul 4": "44",
	"Bucuresti Sectorul 5": "45",
	"Bucuresti Sectorul 6": "46",
	"Calarasi":             "51",
	"Giurgiu":              "52",
}

func randomTimestamp() string {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	randomNow := time.Unix(randomTime, 0).Format("2006-01-02")

	return randomNow
}

func main() {
	/*dateFile, err := os.Create("dates.txt")
	if err != nil {
		fmt.Println(err)
		return
	}


	for i:=0;i<1000000;i++{
		_, err = dateFile.WriteString(randomTimestamp() + "\n")
	}*/

	file, err := ioutil.ReadFile("dispersie_populatie.csv")

	if err != nil {
		fmt.Println("Eroare citire csv.")
	}

	lines := strings.Split(string(file), "\r")

	total := 0
	var valM, valF int
	var dif int
	var totalPopulatie, totalM, totalF int
	for _, line := range lines {
		vals := strings.Split(strings.TrimPrefix(line, "\n"), ",")

		if _, err := judMap[vals[0]]; err {
			valM, _ = strconv.Atoi(vals[1])
			valF, _ = strconv.Atoi(vals[2])
			fmt.Println(fmt.Sprintf("%s %d %d || %d %d", vals[0], valM/dif, valF/dif, 1000000*valM/totalPopulatie, 1000000*valF/totalPopulatie))
			total += 1000000*valM/totalPopulatie + 1000000*valF/totalPopulatie
		} else {
			totalPopulatie, _ = strconv.Atoi(vals[3])
			dif = int(math.Round(float64(totalPopulatie / 1000000)))
		}
	}
	fmt.Println(total)
}
