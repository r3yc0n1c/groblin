package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"math/rand"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"
)

type Domain struct {
	Domains []string `json:"domains"`
}

var defaultUAs struct {
	UserAgents []string `json:"userAgents"`
}

var categories struct {
	Categories []string `json:"categories"`
}

var categoryPattern string


func getDomainList(filePath string) ([]string, error) {
	file, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	switch ext := filepath.Ext(filePath); ext {
	case ".csv":
		csvReader := csv.NewReader(strings.NewReader(string(file)))
		records, err := csvReader.ReadAll()
		if err != nil {
			return nil, err
		}
		return slices.Concat(records[1:]...), nil
	case ".json":
		var records Domain
		if err := json.Unmarshal(file, &records); err != nil {
			return nil, err
		}
		return records.Domains, nil
	default:
		log.Error("Unsupported file extension:", ext)
		os.Exit(1)
	}
	return nil, nil
}

func loadUserAgent() error {
	file, err := os.ReadFile("config/user_agent.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &defaultUAs); err != nil {
		return err
	}

	return nil
}

func loadCategory() error {
	categoryFile, err := os.ReadFile("config/category.json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal(categoryFile, &categories); err != nil {
		return err
	}

	categoryPattern = strings.Join(categories.Categories, "|")
	return nil
}

func CrawlDomain(domain string, wg *sync.WaitGroup, mu *sync.Mutex, results map[string][]string) {
	defer wg.Done()

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.MaxDepth(2),
		colly.CacheDir("./cache"),
	)

	productURLs := []string{}

	// smart regex to match product links
	regexPattern := fmt.Sprintf(`https?://%s/(%s)(/[\w-]+)*`, regexp.QuoteMeta(domain), categoryPattern)

	re := regexp.MustCompile(regexPattern)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if re.MatchString(link) {
			mu.Lock()
			productURLs = append(productURLs, link)
			mu.Unlock()
		}
	})

	c.OnRequest(func(r *colly.Request) {
		userAgent := defaultUAs.UserAgents[rand.Intn(len(defaultUAs.UserAgents))]
		r.Headers.Set("User-Agent", userAgent)		
	})

	err := c.Visit("https://" + domain)
	if err != nil {
		log.Errorf("Error visiting %s: %v", domain, err)
	}

	mu.Lock()
	results[domain] = productURLs
	log.Printf("%s: %d product links found", domain, len(productURLs))
	mu.Unlock()
}

func main() {
	// cmd arg parse
	fileArg := flag.String("file", "", "Enter the file path")
	n := flag.Int("n", 1, "Number of domains to explore at a time")
	flag.Parse()

	filePath := *fileArg

	if filePath == "" {
		log.Error("No file, I'm sleeping...")
		flag.Usage()
		os.Exit(1)
	}

	log.Infof("DOMAIN FILE: %s", filePath)

	domains, err := getDomainList(filePath)

	if err != nil {
		log.Fatal("Error loading domains:", err)
		os.Exit(1)
	}

	// load crawler config
	loadUserAgent()
	loadCategory()

	results := make(map[string][]string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// semaphore to limit concurrency
	sem := make(chan struct{}, *n)

	for _, domain := range domains {
		wg.Add(1)
		sem <- struct{}{}
		go func(domain string) {
			defer func() { <-sem }()
			CrawlDomain(domain, &wg, &mu, results)
		}(domain)
	}

	wg.Wait()

	resultData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling results to JSON: %v", err)
	}

	if err := os.WriteFile("out/results.json", resultData, 0644); err != nil {
		log.Fatalf("Error writing results to file: %v", err)
	}

	log.Debug("Crawling completed. Results saved to results.json")
}
