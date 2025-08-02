package scraper

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniversalScraper_ScrapePage(t *testing.T) {
	// Create a test HTML page
	testHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Test Documentation Page</title>
    <meta name="description" content="Test page for scraping">
    <meta name="keywords" content="test, documentation, scraping">
    <meta name="author" content="Test Author">
</head>
<body>
    <main>
        <h1>JavaScript Tutorial</h1>
        <p>This is a comprehensive guide to JavaScript programming.</p>
        
        <h2>Getting Started</h2>
        <p>JavaScript is a programming language that enables interactive web pages.</p>
        
        <ul>
            <li>Easy to learn</li>
            <li>Versatile language</li>
            <li>Wide browser support</li>
        </ul>
        
        <h3>Code Example</h3>
        <pre><code class="language-javascript">
function greet(name) {
    return "Hello, " + name + "!";
}
console.log(greet("World"));
        </code></pre>
        
        <blockquote>
            JavaScript is the language of the web.
        </blockquote>
    </main>
</body>
</html>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	// Test scraping
	scraper := NewUniversalScraper()
	content, err := scraper.ScrapePage(server.URL)

	require.NoError(t, err)
	assert.NotNil(t, content)
	assert.Equal(t, "JavaScript Tutorial", content.Title)
	assert.Equal(t, "JavaScript", content.Category)
	assert.Contains(t, content.Content, "comprehensive guide")
	assert.Contains(t, content.Tags, "JavaScript")
	assert.Contains(t, content.Tags, "documentation")
	assert.True(t, len(content.Examples) >= 1)
	assert.Contains(t, content.Examples[0], "function greet")
	assert.Equal(t, server.URL, content.URL)
}

func TestUniversalScraper_ExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "H1 title",
			html:     `<html><body><h1>Main Title</h1></body></html>`,
			expected: "Main Title",
		},
		{
			name:     "Article title",
			html:     `<html><body><article><h1>Article Title</h1></article></body></html>`,
			expected: "Article Title",
		},
		{
			name:     "OG title",
			html:     `<html><head><meta property="og:title" content="OG Title"></head><body></body></html>`,
			expected: "OG Title",
		},
		{
			name:     "Document title",
			html:     `<html><head><title>Document Title</title></head><body></body></html>`,
			expected: "Document Title",
		},
		{
			name:     "No title",
			html:     `<html><body></body></html>`,
			expected: "Untitled Document",
		},
	}

	scraper := NewUniversalScraper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			title := scraper.extractTitle(doc)
			assert.Equal(t, tt.expected, title)
		})
	}
}

func TestUniversalScraper_ExtractCategoryFromURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "W3Schools HTML",
			url:      "https://www.w3schools.com/html/html_intro.asp",
			expected: "HTML",
		},
		{
			name:     "MDN JavaScript",
			url:      "https://developer.mozilla.org/en-US/docs/Web/JavaScript",
			expected: "JavaScript",
		},
		{
			name:     "GitHub repository",
			url:      "https://github.com/user/repo",
			expected: "Repository",
		},
		{
			name:     "Python docs",
			url:      "https://docs.python.org/3/tutorial/",
			expected: "Python",
		},
		{
			name:     "React docs",
			url:      "https://react.dev/learn",
			expected: "React",
		},
		{
			name:     "Generic documentation",
			url:      "https://example.com/docs",
			expected: "Documentation",
		},
	}

	scraper := NewUniversalScraper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, err := url.Parse(tt.url)
			require.NoError(t, err)

			doc, err := goquery.NewDocumentFromReader(strings.NewReader("<html><body></body></html>"))
			require.NoError(t, err)

			category := scraper.extractCategoryFromURL(parsedURL, doc)
			assert.Equal(t, tt.expected, category)
		})
	}
}

func TestUniversalScraper_DetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		class    string
		expected string
	}{
		{
			name:     "JavaScript",
			class:    "language-javascript",
			expected: "JavaScript",
		},
		{
			name:     "Python",
			class:    "language-python",
			expected: "Python",
		},
		{
			name:     "HTML",
			class:    "language-html",
			expected: "HTML",
		},
		{
			name:     "No language",
			class:    "some-other-class",
			expected: "",
		},
	}

	scraper := NewUniversalScraper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := `<code class="` + tt.class + `">test code</code>`
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
			require.NoError(t, err)

			selection := doc.Find("code")
			language := scraper.detectLanguage(selection)
			assert.Equal(t, tt.expected, language)
		})
	}
}

func TestUniversalScraper_ConvertToDocument(t *testing.T) {
	content := &ScrapedContent{
		URL:      "https://example.com/test",
		Title:    "Test Document",
		Content:  "This is test content.",
		Category: "Testing",
		Tags:     []string{"test", "documentation"},
		Examples: []string{"console.log('test');"},
		Metadata: map[string]string{
			"author": "Test Author",
		},
	}

	scraper := NewUniversalScraper()
	doc := scraper.ConvertToDocument(content)

	assert.NotNil(t, doc)
	assert.Equal(t, "Test Document", doc.Title)
	assert.Equal(t, "Testing", doc.Category)
	assert.Contains(t, doc.Content, "This is test content.")
	assert.Contains(t, doc.Content, "Code Examples")
	assert.Contains(t, doc.Content, "console.log('test');")
	assert.Equal(t, []string{"test", "documentation"}, doc.Tags)
	assert.Equal(t, "Test Author", doc.Author)
	assert.True(t, strings.HasPrefix(doc.ID, "universal_"))
}

func TestUniversalScraper_ErrorHandling(t *testing.T) {
	scraper := NewUniversalScraper()

	// Test invalid URL
	_, err := scraper.ScrapePage("invalid-url")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch page")

	// Test non-existent server
	_, err = scraper.ScrapePage("http://non-existent-server.example.com")
	assert.Error(t, err)
}

func TestUniversalScraper_HTTPErrorHandling(t *testing.T) {
	// Create test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	scraper := NewUniversalScraper()
	_, err := scraper.ScrapePage(server.URL)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP error: 404")
}

func TestUniversalScraper_DeduplicateTags(t *testing.T) {
	scraper := NewUniversalScraper()
	
	tags := []string{"javascript", "tutorial", "javascript", "documentation", "tutorial", "web"}
	result := scraper.deduplicateTags(tags)
	
	expected := []string{"javascript", "tutorial", "documentation", "web"}
	assert.Equal(t, expected, result)
}

func BenchmarkUniversalScraper_ScrapePage(b *testing.B) {
	testHTML := `
<!DOCTYPE html>
<html>
<head><title>Benchmark Test</title></head>
<body>
    <h1>Test Page</h1>
    <p>This is a test page for benchmarking.</p>
    <pre><code>console.log('test');</code></pre>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(testHTML))
	}))
	defer server.Close()

	scraper := NewUniversalScraper()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scraper.ScrapePage(server.URL)
		if err != nil {
			b.Fatal(err)
		}
	}
}