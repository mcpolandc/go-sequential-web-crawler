package main

import "fmt"

// MapPrettyPrint - prints given data structure
// key1
// 		val[0]
// 		val[1]
// 		val[2]
// key2
// 		val[0]
// 		val[1]
// ..
func PrintSitemap(sitemap *ThreadSafeMap, site string) (err error) {

	fmt.Printf("\n\n*** Printing sitemap for \"%s\" ***\n\n", site)
	// This is really bad, nested for loops! O(n^2)
	// This will be the only synchronous part of the
	// program that happens at the end to wrap up the
	// processing that has occurred
	for key, arr := range sitemap.items {

		fmt.Printf("%s\n", key)

		// example of 'type assertion', .(<type>) after variable
		for _, val := range arr.([]string) {
			fmt.Printf("\t├ %s\n", val)
		}
	}
	return err
}
