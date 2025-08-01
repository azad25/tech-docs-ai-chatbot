package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"tech-docs-ai/internal/app"

	"github.com/PuerkitoBio/goquery"
)

// W3SchoolsScraper scrapes content from W3Schools website.
type W3SchoolsScraper struct {
	client *http.Client
}

// NewW3SchoolsScraper creates a new W3Schools scraper.
func NewW3SchoolsScraper() *W3SchoolsScraper {
	return &W3SchoolsScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ScrapedContent represents the content extracted from a W3Schools page.
type ScrapedContent struct {
	URL       string            `json:"url"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Category  string            `json:"category"`
	Tags      []string          `json:"tags"`
	Examples  []string          `json:"examples"`
	Metadata  map[string]string `json:"metadata"`
	Timestamp time.Time         `json:"timestamp"`
}

// ScrapePage scrapes a single W3Schools page.
func (s *W3SchoolsScraper) ScrapePage(url string) (*ScrapedContent, error) {
	log.Printf("Scraping: %s", url)

	// Make HTTP request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract content
	content := &ScrapedContent{
		URL:       url,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}

	// Extract title
	content.Title = strings.TrimSpace(doc.Find("h1").First().Text())
	if content.Title == "" {
		content.Title = strings.TrimSpace(doc.Find("title").Text())
	}

	// Extract category from URL or breadcrumb
	content.Category = s.extractCategory(url, doc)

	// Extract main content
	content.Content = s.extractMainContent(doc)

	// Extract code examples
	content.Examples = s.extractExamples(doc)

	// Extract tags
	content.Tags = s.extractTags(doc, content.Category)

	// Extract metadata
	content.Metadata = s.extractMetadata(doc, url)

	return content, nil
}

// extractCategory extracts the category from URL or page content.
func (s *W3SchoolsScraper) extractCategory(url string, doc *goquery.Document) string {
	// Try to extract from URL first
	if strings.Contains(url, "/html/") {
		return "HTML"
	} else if strings.Contains(url, "/css/") {
		return "CSS"
	} else if strings.Contains(url, "/js/") {
		return "JavaScript"
	} else if strings.Contains(url, "/python/") {
		return "Python"
	} else if strings.Contains(url, "/sql/") {
		return "SQL"
	} else if strings.Contains(url, "/php/") {
		return "PHP"
	} else if strings.Contains(url, "/java/") {
		return "Java"
	} else if strings.Contains(url, "/cpp/") {
		return "C++"
	} else if strings.Contains(url, "/csharp/") {
		return "C#"
	} else if strings.Contains(url, "/react/") {
		return "React"
	} else if strings.Contains(url, "/bootstrap/") {
		return "Bootstrap"
	} else if strings.Contains(url, "/jquery/") {
		return "jQuery"
	} else if strings.Contains(url, "/nodejs/") {
		return "Node.js"
	} else if strings.Contains(url, "/mongodb/") {
		return "MongoDB"
	} else if strings.Contains(url, "/git/") {
		return "Git"
	} else if strings.Contains(url, "/typescript/") {
		return "TypeScript"
	} else if strings.Contains(url, "/django/") {
		return "Django"
	} else if strings.Contains(url, "/postgresql/") {
		return "PostgreSQL"
	}

	// Try to extract from breadcrumb or navigation
	breadcrumb := doc.Find(".breadcrumb, .nav, .w3-bar").Text()
	if strings.Contains(breadcrumb, "HTML") {
		return "HTML"
	} else if strings.Contains(breadcrumb, "CSS") {
		return "CSS"
	} else if strings.Contains(breadcrumb, "JavaScript") {
		return "JavaScript"
	} else if strings.Contains(breadcrumb, "Python") {
		return "Python"
	} else if strings.Contains(breadcrumb, "SQL") {
		return "SQL"
	}

	return "Web Development"
}

// extractMainContent extracts the main tutorial content.
func (s *W3SchoolsScraper) extractMainContent(doc *goquery.Document) string {
	var content strings.Builder

	// Remove navigation, ads, and other non-content elements
	doc.Find("nav, .w3-bar, .w3-sidebar, .w3-hide, .ad, .advertisement, script, style").Remove()

	// Extract content from main content areas with better structure
	mainSelectors := []string{
		"#main", ".main", ".content", ".tutorial-content",
		".w3-container", ".w3-row", ".w3-col",
		"article", "section", ".chapter",
	}

	for _, selector := range mainSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			// Extract headings with better hierarchy
			s.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
				tagName := s.Get(0).Data
				text := strings.TrimSpace(s.Text())
				if text != "" {
					content.WriteString(fmt.Sprintf("\n%s: %s\n", strings.ToUpper(tagName), text))
				}
			})

			// Extract paragraphs with better formatting
			s.Find("p").Each(func(i int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				if text != "" && len(text) > 10 { // Only include substantial paragraphs
					content.WriteString(text + "\n\n")
				}
			})

			// Extract lists with better structure
			s.Find("ul, ol").Each(func(i int, s *goquery.Selection) {
				listType := "•"
				if s.Get(0).Data == "ol" {
					listType = "1."
				}

				s.Find("li").Each(func(i int, s *goquery.Selection) {
					text := strings.TrimSpace(s.Text())
					if text != "" {
						if s.Get(0).Data == "ol" {
							content.WriteString(fmt.Sprintf("%d. %s\n", i+1, text))
						} else {
							content.WriteString(fmt.Sprintf("%s %s\n", listType, text))
						}
					}
				})
				content.WriteString("\n")
			})

			// Extract tables with better formatting
			s.Find("table").Each(func(i int, s *goquery.Selection) {
				content.WriteString("Table:\n")
				s.Find("tr").Each(func(i int, s *goquery.Selection) {
					var row strings.Builder
					s.Find("td, th").Each(func(i int, s *goquery.Selection) {
						cellText := strings.TrimSpace(s.Text())
						if cellText != "" {
							row.WriteString(cellText + " | ")
						}
					})
					if row.Len() > 0 {
						content.WriteString(row.String() + "\n")
					}
				})
				content.WriteString("\n")
			})

			// Extract definition lists
			s.Find("dl").Each(func(i int, s *goquery.Selection) {
				s.Find("dt").Each(func(i int, s *goquery.Selection) {
					term := strings.TrimSpace(s.Text())
					if term != "" {
						content.WriteString(fmt.Sprintf("Term: %s\n", term))
					}
				})
				s.Find("dd").Each(func(i int, s *goquery.Selection) {
					definition := strings.TrimSpace(s.Text())
					if definition != "" {
						content.WriteString(fmt.Sprintf("Definition: %s\n\n", definition))
					}
				})
			})

			// Extract blockquotes
			s.Find("blockquote").Each(func(i int, s *goquery.Selection) {
				quote := strings.TrimSpace(s.Text())
				if quote != "" {
					content.WriteString(fmt.Sprintf("Quote: %s\n\n", quote))
				}
			})
		})
	}

	return strings.TrimSpace(content.String())
}

// extractExamples extracts code examples from the page with better structure.
func (s *W3SchoolsScraper) extractExamples(doc *goquery.Document) []string {
	var examples []string

	// Find code examples with better selection
	codeSelectors := []string{
		"pre", "code", ".w3-code", ".w3-example", ".example",
		".code-example", ".demo-code", ".tryit-btn",
	}

	for _, selector := range codeSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			code := strings.TrimSpace(s.Text())
			if code != "" && len(code) > 10 { // Only include substantial code examples
				// Clean up the code
				code = strings.ReplaceAll(code, "\n\n", "\n")
				code = strings.TrimSpace(code)

				// Add language hint if available
				if lang, exists := s.Attr("class"); exists {
					if strings.Contains(lang, "html") {
						code = "HTML:\n" + code
					} else if strings.Contains(lang, "css") {
						code = "CSS:\n" + code
					} else if strings.Contains(lang, "js") || strings.Contains(lang, "javascript") {
						code = "JavaScript:\n" + code
					} else if strings.Contains(lang, "python") {
						code = "Python:\n" + code
					} else if strings.Contains(lang, "sql") {
						code = "SQL:\n" + code
					}
				}

				examples = append(examples, code)
			}
		})
	}

	return examples
}

// extractTags extracts relevant tags from the page.
func (s *W3SchoolsScraper) extractTags(doc *goquery.Document, category string) []string {
	tags := []string{category, "tutorial", "documentation", "w3schools"}

	// Add category-specific tags
	switch category {
	case "HTML":
		tags = append(tags, "web", "markup", "semantic")
	case "CSS":
		tags = append(tags, "styling", "design", "layout")
	case "JavaScript":
		tags = append(tags, "programming", "frontend", "es6")
	case "Python":
		tags = append(tags, "programming", "backend", "data-science")
	case "SQL":
		tags = append(tags, "database", "query", "data")
	case "React":
		tags = append(tags, "frontend", "javascript", "framework")
	case "Node.js":
		tags = append(tags, "backend", "javascript", "server")
	}

	// Extract additional tags from content
	content := doc.Text()
	if strings.Contains(content, "API") {
		tags = append(tags, "api")
	}
	if strings.Contains(content, "function") {
		tags = append(tags, "functions")
	}
	if strings.Contains(content, "class") {
		tags = append(tags, "classes")
	}
	if strings.Contains(content, "object") {
		tags = append(tags, "objects")
	}

	return tags
}

// extractMetadata extracts additional metadata from the page.
func (s *W3SchoolsScraper) extractMetadata(doc *goquery.Document, url string) map[string]string {
	metadata := make(map[string]string)

	// Extract meta tags
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			if content, exists := s.Attr("content"); exists {
				metadata[name] = content
			}
		}
	})

	// Extract author information
	author := doc.Find(".author, .byline").Text()
	if author != "" {
		metadata["author"] = strings.TrimSpace(author)
	} else {
		metadata["author"] = "W3Schools"
	}

	// Extract last modified date
	modified := doc.Find(".modified, .updated, .date").Text()
	if modified != "" {
		metadata["last_modified"] = strings.TrimSpace(modified)
	}

	metadata["source"] = "W3Schools"
	metadata["url"] = url

	return metadata
}

// ConvertToDocument converts scraped content to a Document for storage.
func (s *W3SchoolsScraper) ConvertToDocument(content *ScrapedContent) *app.Document {
	// Combine content and examples
	fullContent := content.Content
	if len(content.Examples) > 0 {
		fullContent += "\n\nCode Examples:\n"
		for i, example := range content.Examples {
			fullContent += fmt.Sprintf("\nExample %d:\n%s\n", i+1, example)
		}
	}

	return &app.Document{
		ID:        fmt.Sprintf("w3s_%d", content.Timestamp.UnixNano()),
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
