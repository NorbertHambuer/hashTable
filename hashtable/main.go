package main

import (
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
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

var sqrtHashTable = map[int][]Persoana{}
var primeHashTable = map[int][]Persoana{}
var seqHashTable = map[int][]Persoana{}

var listaCNP []string
var primeNumber = 700001
var listaPersoane []Persoana

type Persoana struct{
	cnp, nume string
}

type hashFunction func(string) int


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

func getNume(gender string) string {
	var result map[string]interface{}

	for {
		res, errReq := http.Get("https://api.namefake.com/romanian-romania/" + gender)

		if errReq != nil {
			continue
		}

		response := res.Body

		body, errRead := ioutil.ReadAll(response)

		if errRead != nil {
			errReq = errRead
			continue
		}

		errUnmarshal := json.Unmarshal(body, &result)

		if errUnmarshal != nil {
			errReq = errUnmarshal
			continue
		}

		numePers := result["name"].(string)
		words := strings.Split(numePers, " ")

		if len(words) > 2 {
			numePers = fmt.Sprintf("%s %s", words[1], words[2])
		}

		return numePers
	}
}

func generareCNP() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	file, err := ioutil.ReadFile("dispersie_populatie.csv")

	if err != nil {
		fmt.Println("Eroare citire csv.")
	}

	dateFile, errFile := os.Create("listaCNP.txt")
	if errFile != nil {
		fmt.Println(errFile)
		return
	}

	persFile, errFile := os.Create("listaPersoane.txt")
	if errFile != nil {
		fmt.Println(errFile)
		return
	}

	lines := strings.Split(string(file), "\r")

	var valM, valF int64
	var codInterval, primaCifra int
	var val int
	var totalPopulatie, totalM, totalF, total int64
	var anN, lunaN, ziN, cnpCurent, dataNasterii, numePers string
	var seed rand.Source
	var random *rand.Rand

	for _, line := range lines {
		vals := strings.Split(strings.TrimPrefix(line, "\n"), ",")

		if codJudet, err := judMap[vals[0]]; err && len(listaCNP) <= 1000000 {
			val, _ = strconv.Atoi(vals[1])
			valM = int64(val)
			val, _ = strconv.Atoi(vals[2])
			valF = int64(val)
			totalM = 1000000 * valM / totalPopulatie
			totalF = 1000000 * valF / totalPopulatie

			total += totalM + totalF
			fmt.Println(fmt.Sprintf("Total %s %d", vals[0], totalM+totalF))
			for i := 1; int64(i) <= totalM+totalF; i++ {
				anN, lunaN, ziN = randomTimestamp()
				dataNasterii = fmt.Sprintf("%s%s%s", anN[2:4], lunaN, ziN)

				if int64(i) <= totalM {
					primaCifra = getPrimaCifraM(anN)
					numePers = getNume("male")
				} else {
					primaCifra = getPrimaCifraF(anN)
					numePers = getNume("female")
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
				_, errFile = persFile.WriteString(fmt.Sprintf("%s, %s \n", numePers, cnpCurent))

				if len(listaCNP) == 1000000 {
					break
				}
			}
		} else {
			val, _ = strconv.Atoi(vals[3])
			totalPopulatie = int64(val)
		}
	}
}

func sqrtHash(cnp string) int {
	cnpInt, _ := strconv.Atoi(cnp)

	return int(math.Sqrt(float64(cnpInt)))
}

func primeHash(cnp string) int {
	cnpInt, _ := strconv.Atoi(cnp)

	return cnpInt % primeNumber
}

func seqHash(cnp string) int {
	val, _ := strconv.Atoi(cnp[1:4])

	return val
}

func adaugareCNP(table map[int][]Persoana, hashFunc hashFunction, pers Persoana){
	index := hashFunc(pers.cnp)

	if lista, err := table[index]; err{
		lista = append(lista, pers)
		table[index] = lista
	}else{
		table[index] = make([]Persoana,1,1000)
		table[index][0] = pers
	}
}

func cautarePersoanaCNP(table map[int][]Persoana, hashFunc hashFunction, cnp string) Persoana{
	index := hashFunc(cnp)

	if lista, err := table[index]; err{
		if len(lista) == 1{
			return lista[0]
		}

		for _, pers := range lista{
			if pers.cnp == cnp{
				return pers
			}
		}

		return Persoana{}
	}

	return Persoana{}
}

func cautareCNP(table map[int][]Persoana, hashFunc hashFunction, cnp string) int{
	index := hashFunc(cnp)
	iterarii := 4

	if lista, err := table[index]; err{
		iterarii++
		if len(lista) == 1{
			iterarii++
			return iterarii
		}

		iterarii+=2
		for _, pers := range lista{
			iterarii++
			if pers.cnp == cnp{
				iterarii++
				return iterarii
			}
		}

		iterarii++
		return iterarii
	}

	iterarii++
	return iterarii
}

func incarcareListaCNP()  {
	f, _ := os.Open("listaPersoane.txt")

	r := csv.NewReader(f)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		pers := Persoana{
			cnp:  strings.TrimSpace(record[1]),
			nume: strings.TrimSpace(record[0]),
		}
		listaPersoane = append(listaPersoane, pers)

		adaugareCNP(sqrtHashTable, sqrtHash, pers)
		adaugareCNP(primeHashTable, primeHash, pers)
		adaugareCNP(seqHashTable, seqHash, pers)
	}
}

func cautereCNPRandom(){
	var seed rand.Source
	var random *rand.Rand
	var in int

	sqrtFile, err := os.Create("sqrt.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	primeFile, err := os.Create("prime.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	seqFile, err := os.Create("seq.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	var it int
	for n := 0; n < 1000; n++ {
		seed = rand.NewSource(time.Now().UnixNano())
		random = rand.New(seed)
		in = random.Intn(1000000-1) + 1

		pers := listaPersoane[in]


		it = cautareCNP(sqrtHashTable,sqrtHash, pers.cnp)
		_, err = sqrtFile.WriteString(fmt.Sprintf("%d \n", it))

		it = cautareCNP(primeHashTable,primeHash, pers.cnp)
		_, err = primeFile.WriteString(fmt.Sprintf("%d \n",  it))

		it = cautareCNP(seqHashTable,seqHash, pers.cnp)
		_, err = seqFile.WriteString(fmt.Sprintf("%d \n", it))

		time.Sleep(10)
	}
}

func main() {
	//generareCNP()
	incarcareListaCNP()
	cautereCNPRandom()
}
