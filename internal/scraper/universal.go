package scraper

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tech-docs-ai/internal/types"

	"github.com/PuerkitoBio/goquery"
)

// UniversalScraper scrapes content from any website with intelligent content extraction
type UniversalScraper struct {
	client *http.Client
}

// NewUniversalScraper creates a new universal scraper
func NewUniversalScraper() *UniversalScraper {
	return &UniversalScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Allow up to 10 redirects
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// ScrapePage scrapes content from any URL with intelligent content extraction
func (s *UniversalScraper) ScrapePage(targetURL string) (*ScrapedContent, error) {
	log.Printf("Universal scraping: %s", targetURL)

	// Validate URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Make HTTP request with proper headers
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; TechDocsAI/1.0; +https://github.com/techdocs-ai)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract content
	content := &ScrapedContent{
		URL:       targetURL,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}

	// Extract title using multiple strategies
	content.Title = s.extractTitle(doc)

	// Extract category based on URL and content
	content.Category = s.extractCategoryFromURL(parsedURL, doc)

	// Extract main content using intelligent selectors
	content.Content = s.extractMainContent(doc)

	// Extract code examples
	content.Examples = s.extractCodeExamples(doc)

	// Extract tags based on content analysis
	content.Tags = s.extractTags(doc, content.Category, parsedURL)

	// Extract metadata
	content.Metadata = s.extractMetadata(doc, targetURL, parsedURL)

	return content, nil
}

// extractTitle extracts the page title using multiple strategies
func (s *UniversalScraper) extractTitle(doc *goquery.Document) string {
	// Strategy 1: Look for main heading
	if title := strings.TrimSpace(doc.Find("h1").First().Text()); title != "" {
		return title
	}

	// Strategy 2: Look for article title
	if title := strings.TrimSpace(doc.Find("article h1, .article-title, .post-title").First().Text()); title != "" {
		return title
	}

	// Strategy 3: Look for page title meta tag
	if title := doc.Find("meta[property='og:title']").AttrOr("content", ""); title != "" {
		return strings.TrimSpace(title)
	}

	// Strategy 4: Use document title
	if title := strings.TrimSpace(doc.Find("title").Text()); title != "" {
		return title
	}

	return "Untitled Document"
}

// extractCategoryFromURL extracts category based on URL patterns and content
func (s *UniversalScraper) extractCategoryFromURL(parsedURL *url.URL, doc *goquery.Document) string {
	host := strings.ToLower(parsedURL.Host)
	path := strings.ToLower(parsedURL.Path)

	// Check for known documentation sites
	switch {
	case strings.Contains(host, "w3schools.com"):
		return s.extractW3SchoolsCategory(path)
	case strings.Contains(host, "developer.mozilla.org") || strings.Contains(host, "mdn.mozilla.org"):
		return s.extractMDNCategory(path)
	case strings.Contains(host, "stackoverflow.com"):
		return "Q&A"
	case strings.Contains(host, "github.com"):
		return "Repository"
	case strings.Contains(host, "docs.python.org"):
		return "Python"
	case strings.Contains(host, "nodejs.org"):
		return "Node.js"
	case strings.Contains(host, "reactjs.org") || strings.Contains(host, "react.dev"):
		return "React"
	case strings.Contains(host, "vuejs.org"):
		return "Vue.js"
	case strings.Contains(host, "angular.io"):
		return "Angular"
	case strings.Contains(host, "go.dev") || strings.Contains(host, "golang.org"):
		return "Go"
	case strings.Contains(host, "rust-lang.org"):
		return "Rust"
	case strings.Contains(host, "java.com") || strings.Contains(host, "oracle.com/java"):
		return "Java"
	}

	// Analyze path for technology indicators
	if category := s.analyzePath(path); category != "" {
		return category
	}

	// Analyze content for technology indicators
	return s.analyzeContent(doc)
}

// extractW3SchoolsCategory extracts category from W3Schools URL
func (s *UniversalScraper) extractW3SchoolsCategory(path string) string {
	switch {
	case strings.Contains(path, "/html/"):
		return "HTML"
	case strings.Contains(path, "/css/"):
		return "CSS"
	case strings.Contains(path, "/js/"):
		return "JavaScript"
	case strings.Contains(path, "/python/"):
		return "Python"
	case strings.Contains(path, "/sql/"):
		return "SQL"
	case strings.Contains(path, "/react/"):
		return "React"
	case strings.Contains(path, "/nodejs/"):
		return "Node.js"
	default:
		return "Web Development"
	}
}

// extractMDNCategory extracts category from MDN URL
func (s *UniversalScraper) extractMDNCategory(path string) string {
	switch {
	case strings.Contains(path, "/html"):
		return "HTML"
	case strings.Contains(path, "/css"):
		return "CSS"
	case strings.Contains(path, "/javascript"):
		return "JavaScript"
	case strings.Contains(path, "/api"):
		return "Web API"
	case strings.Contains(path, "/http"):
		return "HTTP"
	default:
		return "Web Development"
	}
}

// analyzePath analyzes URL path for technology indicators
func (s *UniversalScraper) analyzePath(path string) string {
	technologies := map[string]string{
		"html":       "HTML",
		"css":        "CSS",
		"javascript": "JavaScript",
		"js":         "JavaScript",
		"python":     "Python",
		"java":       "Java",
		"golang":     "Go",
		"rust":       "Rust",
		"php":        "PHP",
		"ruby":       "Ruby",
		"csharp":     "C#",
		"cpp":        "C++",
		"react":      "React",
		"vue":        "Vue.js",
		"angular":    "Angular",
		"nodejs":     "Node.js",
		"docker":     "Docker",
		"kubernetes": "Kubernetes",
		"aws":        "AWS",
		"azure":      "Azure",
		"gcp":        "Google Cloud",
	}

	for tech, category := range technologies {
		if strings.Contains(path, tech) {
			return category
		}
	}

	return ""
}

// analyzeContent analyzes page content for technology indicators
func (s *UniversalScraper) analyzeContent(doc *goquery.Document) string {
	content := strings.ToLower(doc.Text())
	
	// Count occurrences of different technologies
	techCounts := make(map[string]int)
	technologies := map[string]string{
		"html":       "HTML",
		"css":        "CSS",
		"javascript": "JavaScript",
		"python":     "Python",
		"java":       "Java",
		"golang":     "Go",
		"rust":       "Rust",
		"php":        "PHP",
		"ruby":       "Ruby",
		"react":      "React",
		"vue":        "Vue.js",
		"angular":    "Angular",
		"node":       "Node.js",
		"docker":     "Docker",
		"kubernetes": "Kubernetes",
	}

	for tech, category := range technologies {
		count := strings.Count(content, tech)
		if count > 0 {
			techCounts[category] = count
		}
	}

	// Return the most mentioned technology
	maxCount := 0
	var bestCategory string
	for category, count := range techCounts {
		if count > maxCount {
			maxCount = count
			bestCategory = category
		}
	}

	if bestCategory != "" {
		return bestCategory
	}

	return "Documentation"
}

// extractMainContent extracts the main content using intelligent selectors
func (s *UniversalScraper) extractMainContent(doc *goquery.Document) string {
	var content strings.Builder

	// Remove unwanted elements
	doc.Find("script, style, nav, header, footer, aside, .sidebar, .navigation, .menu, .ads, .advertisement, .social, .share").Remove()

	// Try different content selectors in order of preference
	contentSelectors := []string{
		"main",
		"article",
		".content",
		".main-content",
		".post-content",
		".entry-content",
		".article-content",
		".documentation",
		".docs",
		"#content",
		"#main",
		".container .row .col",
		"body",
	}

	var selectedContent *goquery.Selection
	for _, selector := range contentSelectors {
		if selection := doc.Find(selector).First(); selection.Length() > 0 {
			selectedContent = selection
			break
		}
	}

	if selectedContent == nil {
		selectedContent = doc.Find("body")
	}

	// Extract structured content
	s.extractStructuredContent(selectedContent, &content)

	return strings.TrimSpace(content.String())
}

// extractStructuredContent extracts content in a structured way
func (s *UniversalScraper) extractStructuredContent(selection *goquery.Selection, content *strings.Builder) {
	selection.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("h1, h2, h3, h4, h5, h6") {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				level := s.Get(0).Data[1:] // Extract number from h1, h2, etc.
				content.WriteString(fmt.Sprintf("\n%s %s\n", strings.Repeat("#", int(level[0]-'0')), text))
			}
		} else if s.Is("p") {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 10 {
				content.WriteString(text + "\n\n")
			}
		} else if s.Is("ul, ol") {
			s.Find("li").Each(func(i int, li *goquery.Selection) {
				text := strings.TrimSpace(li.Text())
				if text != "" {
					if s.Is("ol") {
						content.WriteString(fmt.Sprintf("%d. %s\n", i+1, text))
					} else {
						content.WriteString(fmt.Sprintf("• %s\n", text))
					}
				}
			})
			content.WriteString("\n")
		} else if s.Is("blockquote") {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				content.WriteString(fmt.Sprintf("> %s\n\n", text))
			}
		} else if s.Is("table") {
			content.WriteString("| Table |\n|-------|\n")
			s.Find("tr").Each(func(i int, tr *goquery.Selection) {
				var row []string
				tr.Find("td, th").Each(func(j int, cell *goquery.Selection) {
					row = append(row, strings.TrimSpace(cell.Text()))
				})
				if len(row) > 0 {
					content.WriteString("| " + strings.Join(row, " | ") + " |\n")
				}
			})
			content.WriteString("\n")
		}
	})
}

// extractCodeExamples extracts code examples from the page
func (s *UniversalScraper) extractCodeExamples(doc *goquery.Document) []string {
	var examples []string

	// Common code selectors
	codeSelectors := []string{
		"pre code",
		"pre",
		".highlight",
		".code",
		".codehilite",
		".language-",
		"code[class*='language-']",
		".example code",
		".code-example",
	}

	for _, selector := range codeSelectors {
		doc.Find(selector).Each(func(i int, sel *goquery.Selection) {
			code := strings.TrimSpace(sel.Text())
			if code != "" && len(code) > 10 {
				// Try to detect language from class
				language := s.detectLanguage(sel)
				if language != "" {
					code = fmt.Sprintf("%s:\n%s", language, code)
				}
				examples = append(examples, code)
			}
		})
	}

	return examples
}

// detectLanguage detects programming language from element classes
func (s *UniversalScraper) detectLanguage(selection *goquery.Selection) string {
	class := selection.AttrOr("class", "")
	
	languages := map[string]string{
		"javascript": "JavaScript",
		"js":         "JavaScript",
		"python":     "Python",
		"java":       "Java",
		"html":       "HTML",
		"css":        "CSS",
		"php":        "PHP",
		"ruby":       "Ruby",
		"go":         "Go",
		"rust":       "Rust",
		"cpp":        "C++",
		"csharp":     "C#",
		"sql":        "SQL",
		"bash":       "Bash",
		"shell":      "Shell",
		"json":       "JSON",
		"xml":        "XML",
		"yaml":       "YAML",
	}

	for lang, name := range languages {
		if strings.Contains(strings.ToLower(class), lang) {
			return name
		}
	}

	return ""
}

// extractTags extracts relevant tags from the page
func (s *UniversalScraper) extractTags(doc *goquery.Document, category string, parsedURL *url.URL) []string {
	tags := []string{category, "documentation"}

	// Add domain-specific tags
	host := strings.ToLower(parsedURL.Host)
	if strings.Contains(host, "github.com") {
		tags = append(tags, "github", "repository")
	} else if strings.Contains(host, "stackoverflow.com") {
		tags = append(tags, "stackoverflow", "qa")
	} else if strings.Contains(host, "medium.com") {
		tags = append(tags, "medium", "article")
	}

	// Extract meta keywords
	if keywords := doc.Find("meta[name='keywords']").AttrOr("content", ""); keywords != "" {
		keywordList := strings.Split(keywords, ",")
		for _, keyword := range keywordList {
			if trimmed := strings.TrimSpace(keyword); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	// Analyze content for additional tags
	content := strings.ToLower(doc.Text())
	contentTags := []string{
		"tutorial", "guide", "documentation", "api", "reference",
		"example", "demo", "sample", "code", "programming",
		"development", "web", "mobile", "frontend", "backend",
	}

	for _, tag := range contentTags {
		if strings.Contains(content, tag) {
			tags = append(tags, tag)
		}
	}

	return s.deduplicateTags(tags)
}

// deduplicateTags removes duplicate tags
func (s *UniversalScraper) deduplicateTags(tags []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, tag := range tags {
		if !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}

	return result
}

// extractMetadata extracts metadata from the page
func (s *UniversalScraper) extractMetadata(doc *goquery.Document, targetURL string, parsedURL *url.URL) map[string]string {
	metadata := make(map[string]string)

	// Basic metadata
	metadata["url"] = targetURL
	metadata["domain"] = parsedURL.Host
	metadata["source"] = "universal-scraper"

	// Extract meta tags
	metaTags := map[string]string{
		"description":    "meta[name='description']",
		"author":         "meta[name='author']",
		"keywords":       "meta[name='keywords']",
		"og:title":       "meta[property='og:title']",
		"og:description": "meta[property='og:description']",
		"og:type":        "meta[property='og:type']",
		"twitter:title":  "meta[name='twitter:title']",
	}

	for key, selector := range metaTags {
		if content := doc.Find(selector).AttrOr("content", ""); content != "" {
			metadata[key] = content
		}
	}

	// Extract publication date
	dateSelectors := []string{
		"meta[property='article:published_time']",
		"meta[name='date']",
		"time[datetime]",
		".date",
		".published",
		".post-date",
	}

	for _, selector := range dateSelectors {
		if date := doc.Find(selector).First().AttrOr("datetime", ""); date != "" {
			metadata["published_date"] = date
			break
		} else if date := strings.TrimSpace(doc.Find(selector).First().Text()); date != "" {
			metadata["published_date"] = date
			break
		}
	}

	// Extract language
	if lang := doc.Find("html").AttrOr("lang", ""); lang != "" {
		metadata["language"] = lang
	}

	return metadata
}

// ConvertToDocument converts scraped content to a Document for storage
func (s *UniversalScraper) ConvertToDocument(content *ScrapedContent) *types.Document {
	// Combine content and examples
	fullContent := content.Content
	if len(content.Examples) > 0 {
		fullContent += "\n\n## Code Examples\n\n"
		for i, example := range content.Examples {
			fullContent += fmt.Sprintf("### Example %d\n\n```\n%s\n```\n\n", i+1, example)
		}
	}

	return &types.Document{
		ID:        fmt.Sprintf("universal_%d", content.Timestamp.UnixNano()),
		Title:     content.Title,
		Content:   fullContent,
		Category:  content.Category,
		Tags:      content.Tags,
		Author:    content.Metadata["author"],
		CreatedAt: content.Timestamp,
		UpdatedAt: content.Timestamp,
		Metadata:  content.Metadata,
	}
}