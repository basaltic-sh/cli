package main

import (
	"strings"
	"unicode"
)

// serviceAliases give the longer service names a short form. Typing
// `loadbalancer` before every command gets old.
var serviceAliases = map[string][]string{
	"loadbalancer": {"lb"},
	"network":      {"net"},
	"certificate":  {"cert"},
	"telemetry":    {"tel"},
	"database":     {"db"},
}

// serviceShort is the one-line description shown in `basaltic --help`.
// Derived from the specification's title, which reads as "Basaltic Compute
// API" and is not what a command listing wants.
var serviceShort = map[string]string{
	"audit":        "Audit logs",
	"billing":      "Invoices, credits, payments and prices",
	"certificate":  "TLS certificates",
	"compute":      "Instances, images, flavors, keypairs and pools",
	"database":     "Managed database clusters",
	"dns":          "DNS zones and records",
	"iam":          "Identity, access, accounts and organizations",
	"kms":          "Encryption keys",
	"loadbalancer": "Load balancers, listeners, rules and target groups",
	"network":      "VPCs, subnets, gateways, routes and security groups",
	"quota":        "Account quotas",
	"secrets":      "Secrets and their versions",
	"storage":      "Volumes, snapshots, buckets and objects",
	"telemetry":    "Logs, metrics and traces",
}

// exported renders a name as an exported Go identifier, matching the SDK's
// own rules closely enough for the private helper names this generator emits.
func exported(s string) string {
	var b strings.Builder
	for _, w := range splitWords(s) {
		if w == "" {
			continue
		}
		b.WriteString(strings.ToUpper(w[:1]))
		b.WriteString(w[1:])
	}
	return b.String()
}

func unexported(s string) string {
	e := exported(s)
	if e == "" {
		return "x"
	}
	return strings.ToLower(e[:1]) + e[1:]
}

func splitWords(s string) []string {
	var words []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 0 {
			words = append(words, strings.ToLower(cur.String()))
			cur.Reset()
		}
	}
	runes := []rune(s)
	for i, r := range runes {
		switch {
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			flush()
		case unicode.IsUpper(r):
			prevLower := i > 0 && unicode.IsLower(runes[i-1])
			prevUpper := i > 0 && unicode.IsUpper(runes[i-1])
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if prevLower || (prevUpper && nextLower) {
				flush()
			}
			cur.WriteRune(r)
		default:
			cur.WriteRune(r)
		}
	}
	flush()
	return words
}

// plural is used for a resource's alias, so `instances` reaches the same
// commands as `instance`.
func plural(s string) string {
	switch {
	case strings.HasSuffix(s, "y") && len(s) > 1 && !isVowel(s[len(s)-2]):
		return strings.TrimSuffix(s, "y") + "ies"
	case strings.HasSuffix(s, "s"), strings.HasSuffix(s, "x"), strings.HasSuffix(s, "z"),
		strings.HasSuffix(s, "ch"), strings.HasSuffix(s, "sh"):
		return s + "es"
	}
	return s + "s"
}

func isVowel(c byte) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

// title capitalises a resource for a help line: "instance-pool" reads as
// "Instance pools".
func title(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return s
	}
	words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
	return strings.Join(words, " ")
}

// firstSentence trims a description down to what fits on a help line.
func firstSentence(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	if i := strings.Index(s, ". "); i > 0 {
		s = s[:i]
	}
	return strings.TrimSuffix(strings.TrimSpace(s), ".")
}
