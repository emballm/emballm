package results

type CodeDetails struct {
	StartLine        int    `json:"start_line"`        // The line number in the source code where the vulnerable code begins
	EndLine          int    `json:"end_line"`          // The line number in the source code where the vulnerable code ends
	StartCol         int    `json:"start_col"`         // The column number in the source code where the vulnerable code starts
	EndCol           int    `json:"end_col"`           // The column number in the source code where the vulnerable code ends
	ShortDescription string `json:"short_description"` // Short description of the code snippet's purpose
}

type RawVulnerability struct {
	Title            string      `json:"title"`             // The name or title of the vulnerability found
	RuleID           string      `json:"rule_id"`           // Identifier that is a snake title of the issue
	Severity         string      `json:"severity"`          // A measure of the potential impact of the vulnerability, often classified as low, medium, high, or critical
	RemediationSteps string      `json:"remediation_steps"` // Instructions or actions needed to fix or mitigate the vulnerability
	IssueDescription string      `json:"issue_description"` // A detailed explanation of the vulnerability, including how it can be exploited and the potential consequences
	ShortDescription string      `json:"short_description"` // A brief summary of the vulnerability
	PriorityScore    float64     `json:"priority_score"`    // EPSS or a numerical value that represents the importance or urgency of addressing the vulnerability
	Code             CodeDetails `json:"code"`              // The section of the software source code that contains the vulnerability
	CodeSnippet      string      `json:"code_snippet"`      // A small segment or excerpt of the source code where the vulnerability is located
}

type Issue struct {
	Title            string  `json:"title"`
	Severity         float32 `json:"severity"`
	FileName         string  `json:"file_name"`
	RemediationSteps string  `json:"remediation_steps"`
	IssueDescription string  `json:"issue_description"`
	ShortDescription string  `json:"short_description"`
	LineNumber       int     `json:"line_number"`
	CodeSnippet      string  `json:"code_snippet"`
}

type Meta struct {
	Key        []string `json:"key"`
	SubProduct string   `json:"subproduct"`
}

type Data struct {
	Meta   Meta    `json:"meta"`
	Issues []Issue `json:"issues"`
}
