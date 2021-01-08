package controller

import "strings"

// CheckNumberOfThreads ... (when user create post)
func CheckNumberOfThreads(input string) []string {
	new := strings.ToLower(input)
	done := []string{}
	var one string
	arrThreads := strings.Split(new, ",")
	for _, thread := range arrThreads {
		if thread == "" || thread == " " {
			continue
		}
		for _, symbol := range thread {
			if symbol != ' ' && ((symbol > 47 && symbol < 58) || (symbol >= 'a' && symbol <= 'z') || (symbol >= 'A' && symbol <= 'Z')) {
				one += string(symbol)
			}
		}
		if one != "" && one != " " {
			done = append(done, one)
		}
		one = ""
	}
	return done

}
