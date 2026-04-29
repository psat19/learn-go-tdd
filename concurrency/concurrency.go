package main

type WebsiteChecker func(string) bool

type Result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	resultChannel := make(chan Result)

	for _, url := range urls {
		go func(u string) {
			resultChannel <- Result{url, wc(u)}
		}(url)
	}

	for range len(urls) {
		r := <-resultChannel
		results[r.string] = r.bool
	}

	return results
}
