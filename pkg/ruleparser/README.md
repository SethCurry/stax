# ruleparser

`ruleparser` is a utility format for parsing the Magic: The Gathering
comprehensive rules into an intermediate format that can be used to
generate HTML, JSON, or other formats.

## Usage

You will typically only need to interact with `ruleparser.ParseRulesFile` if
you have a .txt of the comprehensive rules that you would like to parse. If
you have an io.Reader, you can use `ruleparser.ParseRules` instead.

You can get a copy of the rules from the [Wizards of the Coast website](https://magic.wizards.com/en/rules). Download the TXT file.

```go
package main

import (
  "fmt"
  "os"

  "github.com/SethCurry/stax/pkg/ruleparser"
)

func main() {
  rules, err := ruleparser.ParseRulesFile("rules.txt")
  if err != nil {
    panic(err)
  }

  // do something with the rules
}
```
