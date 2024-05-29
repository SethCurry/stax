package ruleparser

import (
	"regexp"
	"strings"
)

// ContentType represents the type of content in a rule, such as text,
// a reference to another rule, or a symbol.
type ContentType string

const (
	// ContentText represents a plain text content element.
	ContentText ContentType = "text"

	// ContentSymbol represents a symbol in the contents,
	// such as a mana symbol or the tap symbol.
	ContentSymbol ContentType = "symbol"

	// ContentReference represents a reference to another rule.
	ContentReference ContentType = "reference"
)

type ContentElement struct {
	Type  ContentType `json:"type"`
	Value string      `json:"value"`
}

//nolint:funlen
func parseContent(content string) []*ContentElement {
	refRegex := regexp.MustCompile(`[0-9]{3}\.[0-9]+[a-z]?`)
	elements := []*ContentElement{}
	acc := ""

	scanner := NewStringScanner(content)

	for {
		character, ok := scanner.Next()
		if !ok {
			break
		}

		switch character {
		case '{':
			if len(acc) > 0 {
				elements = append(elements, &ContentElement{
					Type:  ContentText,
					Value: acc,
				})
				acc = ""
			}

			value := scanner.ReadUntil([]rune{'}'})

			elements = append(elements, &ContentElement{
				Type:  ContentSymbol,
				Value: value,
			})

			scanner.Next()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			maybeRef := string(character) + scanner.ReadUntil([]rune{' ', ')', ',', '-'})
			maybeRef = strings.TrimSuffix(maybeRef, ".")

			lastElEndsInSeeRule := strings.HasSuffix(acc, "ee rule ")

			if refRegex.MatchString(maybeRef) || lastElEndsInSeeRule {
				if len(acc) > 0 {
					elements = append(elements, &ContentElement{
						Type:  ContentText,
						Value: acc,
					})
					acc = ""
				}

				elements = append(elements, &ContentElement{
					Type:  ContentReference,
					Value: maybeRef,
				})
			} else {
				acc += maybeRef
			}
		default:
			acc += string(character)
		}
	}

	if len(acc) > 0 {
		elements = append(elements, &ContentElement{
			Type:  ContentText,
			Value: acc,
		})
	}

	return elements
}

func convertEncoding(old string) string {
	old = strings.Replace(old, "“", "\"", -1)
	old = strings.Replace(old, "”", "\"", -1)
	old = strings.Replace(old, "’", "'", -1)
	old = strings.Replace(old, "™", "(tm)", -1)
	old = strings.Replace(old, "–", "-", -1)

	return old
}
