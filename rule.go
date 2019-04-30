package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidRule              = func(rule string) error { return fmt.Errorf("invalid rule: %s", rule) }
	ErrInvalidRulePartForRegexp = func(rulePart string) error { return fmt.Errorf("invalid rule part %s", rulePart) }
	ErrExcMarkShouldBeFirst     = errors.New("exclamation mark should be first at the rule")
)

type Ruler struct {
	rule *Rule
}

type Rule struct {
	color     string
	atSignRe  *regexp.Regexp
	excMarkRe *regexp.Regexp
}

func Parse(ruleAsAString string) (*Rule, error) {
	excMarkIndex := strings.Index(ruleAsAString, "!")
	atSignMarkIndex := strings.Index(ruleAsAString, "@")

	if excMarkIndex == -1 && atSignMarkIndex == -1 {
		return nil, ErrInvalidRule(ruleAsAString)
	}

	// excMark should be greater then atSign and atSign is found
	if excMarkIndex > atSignMarkIndex && atSignMarkIndex != -1 {
		return nil, ErrExcMarkShouldBeFirst
	}

	rule := &Rule{
		color: RandColor(),
	}
	var err error

	if atSignMarkIndex == -1 {
		rule.excMarkRe, err = parsePart(ruleAsAString, excMarkIndex+1, len(ruleAsAString))
		if err != nil {
			return nil, err
		}
	} else {
		rule.excMarkRe, err = parsePart(ruleAsAString, excMarkIndex+1, atSignMarkIndex)
		if err != nil {
			return nil, err
		}

		rule.atSignRe, err = parsePart(ruleAsAString, atSignMarkIndex+1, len(ruleAsAString))
		if err != nil {
			return nil, err
		}
	}

	return rule, nil
}

func parsePart(ruleAsAString string, start, end int) (*regexp.Regexp, error) {
	rulePart := ruleAsAString[start:end]

	re, err := regexp.Compile(rulePart)
	if err != nil {
		return nil, ErrInvalidRulePartForRegexp(rulePart)
	}

	return re, nil
}

func ExecRule(rule *Rule, logStr string) {
	if !rule.excMarkRe.MatchString(logStr) {
		log.Log().Msg("not matched")
		return
	}

	if rule.atSignRe != nil {
		s := rule.atSignRe.FindString(logStr)
		s = fmt.Sprintf("<span style=\"%s\">%s</span>", rule.color, s)
		logStr = rule.atSignRe.ReplaceAllString(logStr, s)
	}

	fmt.Println(logStr)
}