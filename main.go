package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"gopkg.in/ini.v1"
)

var countries = map[string][]string{}

type Delegate_weight struct {
	delegate string
	weight   int
}

func main() {

	delegates_file, config_file, history_file, assignments_file := get_filenames()

	// initialize vars
	get_countries(config_file)
	delegates := get_delegates(delegates_file)
	all_countries := get_all_countries()
	important_countries := shuffle_slice(get_important_countries())
	previous_assignments, min_assignments := get_previous_assignments(history_file)
	available_countries := shuffle_slice(all_countries)

	remove_chairs(&delegates, delegates_file)
	var assignments = make(map[string]string, len(delegates))

	handle_weighted_delegate_assignments(&delegates, &previous_assignments, config_file, min_assignments, &important_countries, &available_countries, &assignments)

	// assign the remaining delegates
	for i, delegate := range delegates {
		if i >= len(available_countries) || i >= len(delegates) {
			assignments[delegate] = "No more countries"
			continue
		}
		assignments[delegate] = available_countries[i]
	}

	write_assignments(assignments_file, assignments)
	write_history(history_file, assignments, previous_assignments)
}

func in_slice(a string, slice []string) bool {
	for _, element := range slice {
		if element == a {
			return true
		}
	}
	return false
}

func slice_element_index(a string, slice []string) int {
	for i, element := range slice {
		if element == a {
			return i
		}
	}
	return -1
}

func shuffle_slice(src []string) []string {
	dest := make([]string, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}

func get_delegates(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var delegates []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		delegate := scanner.Text()
		delegates = append(delegates, strings.TrimRight(delegate, "\n"))
	}

	return delegates
}

func write_assignments(filename string, assignments map[string]string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for delegate, assignment := range assignments {
		fmt.Fprintf(writer, "%s,%s\n", delegate, assignment)
	}

	writer.Flush()
}

func write_history(filename string, assignments map[string]string, previous_assignments map[string][]string) {

	var history = make(map[string][]string)
	for delegate, previous_countries := range previous_assignments {
		var assigned_countries []string
		for _, country := range previous_countries {
			assigned_countries = append(assigned_countries, country)
		}
		history[delegate] = assigned_countries
	}

	for delegate, country := range assignments {
		history[delegate] = append(history[delegate], country)
	}

	// make lines
	var lines []string
	for delegate, assigned_countries := range history {
		line := delegate
		for _, country := range assigned_countries {
			line += "," + country
		}
		lines = append(lines, line)
	}

	// remove old history file
	if _, err := os.Stat(filename); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			err := os.Remove(filename)
			if err != nil {
				panic(err)
			}
		}
	}

	file, error := os.Create(filename)
	if error != nil {
		panic(error)
	}
	defer file.Close()

	// write the lines
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintf(writer, "%s\n", line)
	}
	writer.Flush()
}

func get_countries(filename string) {
	config_data, err := ini.Load(filename)
	if err != nil {
		panic(err)
	}

	section := config_data.Section("P5")
	countries["P5"] = section.KeyStrings()
	section = config_data.Section("High")
	countries["High"] = section.KeyStrings()
	section = config_data.Section("Medium")
	countries["Medium"] = section.KeyStrings()
	section = config_data.Section("Standard")
	countries["Standard"] = section.KeyStrings()
}

func get_all_countries() []string {
	var all_countries []string
	for _, value := range countries {
		all_countries = append(all_countries, value...)
	}
	return all_countries
}

func get_important_countries() []string {
	return append(countries["P5"], countries["High"]...)
}

func get_previous_assignments(filename string) (map[string][]string, int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	assignments := make(map[string][]string)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var min_assignments int = -1
	for _, line := range lines {
		elements := strings.Split(line, ",")
		assignments[elements[0]] = elements[1:]
		assignments_len := len(elements[1:])
		if min_assignments == -1 {
			min_assignments = assignments_len
		} else if min_assignments > assignments_len {
			min_assignments = assignments_len
		}
	}

	return assignments, min_assignments + 1
}

func get_delegate_weight(previous_assignments []string, filename string, min_assignments int) int {
	if min_assignments == 0 {
		return 0
	}

	delegate_weight := 0

	config, err := ini.Load(filename)
	if err != nil {
		panic(err)
	}

	section := config.Section("General")
	i := 0
	for _, country := range previous_assignments {

		i++
		if i > min_assignments {
			break
		}

		if in_slice(country, countries["P5"]) {
			p5_weight, error := section.Key("P5").Int()
			if error != nil {
				panic(error)
			}
			delegate_weight += p5_weight
		} else if in_slice(country, countries["Important"]) {
			important_weight, error := section.Key("Important").Int()
			if error != nil {
				panic(error)
			}
			delegate_weight += important_weight
		} else if in_slice(country, countries["Medium"]) {
			medium_weight, error := section.Key("Medium").Int()
			if error != nil {
				panic(error)
			}
			delegate_weight += medium_weight
		} else if in_slice(country, countries["Standard"]) {
			standard_weigt, error := section.Key("Standard").Int()
			if error != nil {
				panic(error)
			}
			delegate_weight += standard_weigt
		}
	}

	salt_max, error := section.Key("Salt").Int()
	if error != nil {
		panic(error)
	}
	delegate_weight += rand.Intn(salt_max)

	return delegate_weight
}

func get_filenames() (string, string, string, string) {
	delegates_file := "delegates.csv"
	config_file := "config.ini"
	history_file := "history.csv"
	assignments_file := "assignments.csv"
	return delegates_file, config_file, history_file, assignments_file
}

func handle_weighted_delegate_assignments(delegates *[]string, previous_assignments *map[string][]string, config_file string, min_assignments int, important_countries *[]string, available_countries *[]string, assignments *map[string]string) {
	var delegate_weights = make(map[string]int, len(*delegates))

	// setup weights for delegates
	for _, delegate := range *delegates {
		delegate_weights[delegate] = get_delegate_weight((*previous_assignments)[delegate], config_file, min_assignments)
	}

	var weighted_delegates []Delegate_weight
	for delegate, weight := range delegate_weights {
		weighted_delegates = append(weighted_delegates, Delegate_weight{delegate, weight})
	}

	sort.Slice(*delegates, func(i, j int) bool {
		return weighted_delegates[i].weight < weighted_delegates[j].weight
	})

	// assign the weighted delegates
	for i, country := range *important_countries {
		if !in_slice(country, *available_countries) {
			continue
		}
		delegate := (*delegates)[0]
		(*assignments)[delegate] = country
		// deleting the country and delegate from respective slice
		(*delegates) = (*delegates)[1:]
		(*available_countries)[i] = (*available_countries)[len(*available_countries)-1]
		(*available_countries) = (*available_countries)[:len(*available_countries)-1]
	}
}

func remove_chairs(delegates *[]string, delegates_file string) {
	fmt.Printf("Enter the name of the chairs as they are present in the %s file. Make sure they are separated by \",\" (no spaces):", delegates_file)
	var chair_string string
	fmt.Scanf("%s", &chair_string)
	chairs := strings.Split(chair_string, ",")
	for _, chair := range chairs {
		i := slice_element_index(chair, *delegates)
		(*delegates)[i] = (*delegates)[len(*delegates)-1]
		(*delegates) = (*delegates)[:len(*delegates)-1]
	}
}
