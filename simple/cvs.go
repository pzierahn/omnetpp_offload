package simple

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func sortDictionaryInt(dict map[string]int) []string {

	terms := make([]string, len(dict))

	inx := 0
	for key := range dict {
		terms[inx] = key
		inx++
	}

	sort.Strings(terms)

	sort.SliceStable(terms, func(i, j int) bool {
		return dict[terms[i]] > dict[terms[j]]
	})

	return terms
}

func sortDictionaryFloat(dict map[string]float64) []string {

	terms := make([]string, len(dict))

	inx := 0
	for key := range dict {
		terms[inx] = key
		inx++
	}

	sort.Slice(terms, func(i, j int) bool {
		return dict[terms[i]] > dict[terms[j]]
	})

	return terms
}

func SortToCVSFloat(dict map[string]float64) string {

	sortTerms := sortDictionaryFloat(dict)

	cvs := ""

	for _, term := range sortTerms {
		cvs += fmt.Sprintf("%s, %.19f\n", term, dict[term])
	}

	return cvs
}

func SortToCVSInt(dict map[string]int) string {

	sortTerms := sortDictionaryInt(dict)

	cvs := ""

	for _, term := range sortTerms {
		cvs += fmt.Sprintf("%s, %d\n", term, dict[term])
	}

	return cvs
}

func WriteCVSInt(filename string, dict map[string]int) {

	cvs := SortToCVSInt(dict)
	err := ioutil.WriteFile(filename, []byte(cvs), 0755)
	if err != nil {
		log.Panic(err)
	}
}

func WriteCVSFloat(filename string, dict map[string]float64) {

	cvs := SortToCVSFloat(dict)
	err := ioutil.WriteFile(filename, []byte(cvs), 0755)
	if err != nil {
		log.Panic(err)
	}
}

func WriteMultiCSVFloat(filename string, heads []string, content [][]float64) {
	if len(heads) != len(content) {
		log.Panicln("len(heads) != len(content)")
	}

	text := strings.Join(heads, ",") + "\n"

	max := 0
	for _, line := range content {
		if max < len(line) {
			max = len(line)
		}
	}

	for inx := 0; inx < max; inx++ {

		parts := make([]string, len(content))

		for iny, line := range content {
			parts[iny] = fmt.Sprint(line[inx])
		}

		text += strings.Join(parts, ",") + "\n"
	}

	text = strings.TrimSpace(text)

	err := ioutil.WriteFile(filename, []byte(text), 0644)
	if err != nil {
		log.Panic(err)
	}
}
