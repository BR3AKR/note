<!--META--
author: .Author
created: {{.Date.Format "2006-01-02T15:04:05Z07:00"}}
edited: {{.Date.Format "2006-01-02T15:04:05Z07:00"}}
title: {{.Title}}
subtitle: {{.Subtitle}}
tags:{{range $i, $t := .Tags}}
  - {{.}}{{end}}
--END-->