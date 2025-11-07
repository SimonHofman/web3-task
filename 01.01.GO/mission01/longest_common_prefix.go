package main

import "fmt"

func longestCommonPrefix(strs []string) string {
	fmt.Printf("%s: ", strs)
	char_arr := []byte(strs[0])
	for i := 0; i < len(char_arr); i++ {
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || char_arr[i] != strs[j][i] {
				return string(char_arr[:i])
			}
		}
	}
	return ""
}

func main() {
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
	fmt.Println(longestCommonPrefix([]string{"dog", "racecar", "car"}))
}
