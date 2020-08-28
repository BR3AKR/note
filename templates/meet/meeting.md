---
author: {{.Author}}
created: {{.Date.Format "2006-01-02T15:04:05Z07:00"}}
tags: 
  - meeting
---
# {{.Date.Format "2006-01-02"}} - {{.Title}}

**Participants**: {{range $i, $a := .Participants}}{{- if $i}}, {{end}}{{- $a}}{{end}}

## Action Items

| ? | Owner | Notes |
| - | ----- | ----- |
|✅❌| Sean  | Something is done or it isn't |

## Agenda
{{range $i, $t := .Agenda}}
### {{.}}
{{end}}