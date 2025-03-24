package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"sort"
	"strings"

	"gopkg.in/ini.v1"
)

var countries = map[string][]string{}

type Delegate_weight struct {
	delegate string
	weight   int
}

const (
	delegates_file   = "delegates.csv"
	config_file      = "config.ini"
	history_file     = "history.csv"
	assignments_file = "assignments.csv"
)

func main() {

	// initialize vars
	get_countries()
	delegates := get_delegates()
	all_countries := get_all_countries()
	important_countries := shuffle_slice(append(countries["P5"], countries["High"]...))
	previous_assignments, min_assignments := get_previous_assignments()
	available_countries := shuffle_slice(all_countries)

	remove_chairs(&delegates)
	var assignments = make(map[string]string, len(delegates))

	if previous_assignments != nil {
		handle_weighted_delegate_assignments(&delegates, &previous_assignments, min_assignments, &important_countries, &available_countries, &assignments)
	}

	if len(important_countries) > 0 {
		for _, country := range important_countries {
			i := slice_element_index(country, available_countries)
			if i == -1 {
				continue
			}
			delegate := delegates[0]
			assignments[delegate] = country
			important_countries = important_countries[1:]
			// deleting the country and delegate from respective slice
			delegates = delegates[1:]
			available_countries[i] = available_countries[len(available_countries)-1]
			available_countries = available_countries[:len(available_countries)-1]
		}
	}

	// assign the remaining delegates
	for i, delegate := range delegates {
		if i >= len(available_countries) || i >= len(delegates) {
			assignments[delegate] = "No more countries"
			continue
		}
		assignments[delegate] = available_countries[i]
	}

	write_assignments(assignments)
	write_history(assignments, previous_assignments)
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

func get_delegates() []string {
	file, err := os.Open(delegates_file)
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

func write_assignments(assignments map[string]string) {
	file, err := os.Create(assignments_file)
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

func write_history(assignments map[string]string, previous_assignments map[string][]string) {

	var history = make(map[string][]string)
	if previous_assignments != nil {
		for delegate, previous_countries := range previous_assignments {
			var assigned_countries []string
			for _, country := range previous_countries {
				assigned_countries = append(assigned_countries, country)
			}
			history[delegate] = assigned_countries
		}
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

	file, error := os.Create(history_file)
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

func get_countries() {
	config_data, err := ini.Load(config_file)
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

func get_previous_assignments() (map[string][]string, int) {
	file, err := os.Open(assignments_file)
	if errors.Is(err, os.ErrNotExist) {
		return nil, 0
	} else if err != nil {
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

	return assignments, min_assignments + 2
}

func get_delegate_weight(previous_assignments []string, min_assignments int) int {
	if min_assignments == 0 {
		return 0
	}

	delegate_weight := 0

	config, err := ini.Load(config_file)
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

		if slices.Contains(countries["P5"], country) {
			p5_weight, error := section.Key("P5").Int()
			if error != nil {
				panic(error)
			}
			delegate_weight += p5_weight
		} else if slices.Contains(countries["Important"], country) {
			important_weight, error := section.Key("Important").Int()
			if error != nil {
				panic(error)
			}
			delegate_weight += important_weight
		} else if slices.Contains(countries["Medium"], country) {
			medium_weight, error := section.Key("Medium").Int()
			if error != nil {
				panic(error)
			}
			delegate_weight += medium_weight
		} else if slices.Contains(countries["Standard"], country) {
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

func handle_weighted_delegate_assignments(delegates *[]string, previous_assignments *map[string][]string, min_assignments int, important_countries *[]string, available_countries *[]string, assignments *map[string]string) {
	var delegate_weights = make(map[string]int, len(*delegates))

	// setup weights for delegates
	for _, delegate := range *delegates {
		delegate_weights[delegate] = get_delegate_weight((*previous_assignments)[delegate], min_assignments)
	}

	var weighted_delegates []Delegate_weight
	for delegate, weight := range delegate_weights {
		weighted_delegates = append(weighted_delegates, Delegate_weight{delegate, weight})
	}

	sort.Slice(*delegates, func(i, j int) bool {
		return weighted_delegates[i].weight < weighted_delegates[j].weight
	})

	// assign the weighted delegates
	for i := range len(*important_countries) {
		j := slice_element_index((*important_countries)[0], *available_countries)
		if i == -1 {
			continue
		}
		delegate := (*delegates)[0]
		(*assignments)[delegate] = (*important_countries)[0]
		*important_countries = (*important_countries)[1:]
		// deleting the country and delegate from respective slice
		(*delegates) = (*delegates)[1:]
		(*available_countries)[j] = (*available_countries)[len(*available_countries)-1]
		(*available_countries) = (*available_countries)[:len(*available_countries)-1]
	}
}

func remove_chairs(delegates *[]string) {
	fmt.Printf("Enter the name of the chairs as they are present in the %s file. Make sure they are separated by \",\" (no spaces): ", delegates_file)
	var chair_string string
	fmt.Scanf("%s", &chair_string)
	if chair_string == "" {
		return
	}
	chairs := strings.Split(chair_string, ",")
	for _, chair := range chairs {
		i := slice_element_index(chair, *delegates)
		(*delegates)[i] = (*delegates)[len(*delegates)-1]
		(*delegates) = (*delegates)[:len(*delegates)-1]
	}
}
