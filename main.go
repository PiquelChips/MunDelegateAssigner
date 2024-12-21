package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

var countries = map[string][]string{
	"P5": {
		"United States",
		"China",
		"Russia",
		"United Kingdom",
		"France",
	},
	"High": {
		"Germany",
		"Japan",
		"India",
		"Brazil",
		"Canada",
		"Australia",
		"South Korea",
		"Italy",
		"Spain",
		"Saudi Arabia",
	},
	"Medium": {
		"Mexico",
		"Indonesia",
		"Turkey",
		"Netherlands",
		"Switzerland",
		"Sweden",
		"Poland",
		"Argentina",
		"Nigeria",
		"South Africa",
		"Egypt",
		"Pakistan",
		"Vietnam",
		"UAE",
		"Israel",
	},
	"Standard": {
		"Malaysia",
		"Singapore",
		"Thailand",
		"Philippines",
		"Chile",
		"Peru",
		"Colombia",
		"Morocco",
		"Kenya",
		"Ghana",
		"Ethiopia",
		"Iraq",
		"Kuwait",
		"Qatar",
		"New Zealand",
		"Portugal",
		"Ireland",
		"Greece",
		"Denmark",
		"Norway",
		"Finland",
	},
}

func in_slice(a string, slice []string) bool {
	for _, element := range slice {
		if element == a {
			return true
		}
	}
	return false
}

func shuffle_slice(src []string) []string {
	dest := make([]string, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}

func main() {

	var input_file_name, output_file_name string
	// Get all needed input from the user
	fmt.Printf("Enter relative path of the file containing the list of delegates. The file should contain one delegate per line: ")
	fmt.Scanf("%s", &input_file_name)
	fmt.Printf("Enter relative path of the file the assignments will be written to: ")
	fmt.Scanf("%s", &output_file_name)

	// Read delegates from the file
	file, err := os.Open(input_file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var delegates []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		delegates = append(delegates, scanner.Text())
	}

	// Get all countries into a slice
	var all_countries []string
	for _, value := range countries {
		all_countries = append(all_countries, value...)
	}

	// Create variables needed for the assignments
	var assignments = make(map[string]string, len(delegates))
	shuffled_contries := shuffle_slice(all_countries)

	for i, delegate := range delegates {
        if i >= len(shuffled_contries) || i >= len(delegates) {
            assignments[delegate] = "No more countries"
            continue
        }
		assignments[delegate] = shuffled_contries[i]
	}

	// Write the assignments to the file
	assignments_file, error := os.Create(output_file_name)
	if error != nil {
		panic(error)
	}
	defer assignments_file.Close()

	writer := bufio.NewWriter(assignments_file)
	for delegate, assignment := range assignments {
        fmt.Fprintf(writer, "%s: %s\n", delegate, assignment)
	}
    writer.Flush()
}
