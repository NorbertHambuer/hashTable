package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
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

var listaCNP []string

func randomTimestamp() (string, string, string) {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000

	randomNow := time.Unix(randomTime, 0)

	return randomNow.Format("2006"), randomNow.Format("01"), randomNow.Format("02")
}

func getPrimaCifraM(anN string) int {
	anNastere, _ := strconv.Atoi(anN)

	if anNastere >= 1900 && anNastere <= 1999 {
		return 1
	}

	if anNastere >= 1800 && anNastere <= 1899 {
		return 3
	}

	return 5
}

func getPrimaCifraF(anN string) int {
	anNastere, _ := strconv.Atoi(anN)

	if anNastere >= 1900 && anNastere <= 1999 {
		return 2
	}

	if anNastere >= 1800 && anNastere <= 1899 {
		return 4
	}

	return 6
}

func validareIntervalCNP(dataNasterii string, cod string) bool {
	var codCNP, dataCNP string
	for _, cnp := range listaCNP {
		dataCNP = cnp[1:7]

		if dataCNP == dataNasterii {
			codCNP = cnp[9:12]

			if codCNP == cod {
				return false
			}
		}
	}

	return true
}

func getCifraControl(cnp string) int {
	rn := []rune(cnp)
	control := []int{2, 7, 9, 1, 4, 6, 3, 5, 8, 2, 7, 9}
	sum := 0

	for index, cifra := range rn {
		sum += control[index] * int(cifra-'0')
	}

	result := sum % 11

	if result == 10 {
		return 1
	}

	return result
}

func main() {
	file, err := ioutil.ReadFile("dispersie_populatie.csv")

	if err != nil {
		fmt.Println("Eroare citire csv.")
	}

	dateFile, errFile := os.Create("listaCNP.txt")
	if errFile != nil {
		fmt.Println(errFile)
		return
	}

	lines := strings.Split(string(file), "\r")

	var valM, valF, codInterval, primaCifra int
	var totalPopulatie, totalM, totalF int
	var anN, lunaN, ziN, cnpCurent, dataNasterii string
	var seed rand.Source
	var random *rand.Rand
	total := 0
	for _, line := range lines {
		vals := strings.Split(strings.TrimPrefix(line, "\n"), ",")

		if codJudet, err := judMap[vals[0]]; err && len(listaCNP) <= 1000000 {
			valM, _ = strconv.Atoi(vals[1])
			valF, _ = strconv.Atoi(vals[2])
			totalM = 1000000 * valM / totalPopulatie
			totalF = 1000000 * valF / totalPopulatie

			total += totalM + totalF

			for i := 1; i <= totalM+totalF; i++ {
				anN, lunaN, ziN = randomTimestamp()
				dataNasterii = fmt.Sprintf("%s%s%s", anN[2:4], lunaN, ziN)

				if i <= totalM {
					primaCifra = getPrimaCifraM(anN)
				} else {
					primaCifra = getPrimaCifraF(anN)
				}

				seed = rand.NewSource(time.Now().UnixNano())
				random = rand.New(seed)
				codInterval = random.Intn(999-1) + 1

				for !validareIntervalCNP(dataNasterii, fmt.Sprintf("%03d", codInterval)) {
					seed = rand.NewSource(time.Now().UnixNano())
					random = rand.New(seed)
					codInterval = random.Intn(999-1) + 1
				}

				cnpCurent = fmt.Sprintf("%d%s%s%s", primaCifra, dataNasterii, codJudet, fmt.Sprintf("%03d", codInterval))

				cnpCurent = fmt.Sprintf("%s%d", cnpCurent, getCifraControl(cnpCurent))

				listaCNP = append(listaCNP, cnpCurent)
				_, errFile = dateFile.WriteString(cnpCurent + "\n")

				if len(listaCNP) == 1000000 {
					break
				}
			}
		} else {
			totalPopulatie, _ = strconv.Atoi(vals[3])
		}

		fmt.Println(vals[0])
	}
}
