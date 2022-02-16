package parser

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

type (
	regexNamedGroup    string
	hintTypeDescriptor string
	hint               struct {
		variable       string
		typeDescriptor hintTypeDescriptor
		value          string
		position       int
	}
)

const (
	hintTypeName         = hintTypeDescriptor("Name")
	hintTypeDescription  = hintTypeDescriptor("Description")
	hintTypeDefaultValue = hintTypeDescriptor("Default")
	hintTypeValues       = hintTypeDescriptor("Values")
	hintTypeInvalid      = hintTypeDescriptor("invalid")

	regexNamedGroupVariable = regexNamedGroup("varname")
	regexNamedGroupType     = regexNamedGroup("key")
	regexNamedGroupValue    = regexNamedGroup("value")
)

var hintRegex = regexp.MustCompile(fmt.Sprintf(
	"^# \\$\\{(?P<%s>\\S+)\\} (?P<%s>\\S+): (?P<%s>.+)$",
	regexNamedGroupVariable,
	regexNamedGroupType,
	regexNamedGroupValue,
))

func ParseParameters(snippet string) []model.Parameter {
	hints := parseHints(snippet)
	return hintsToParameters(hints)
}

func CreateSnippet(snippet string, parameters []model.Parameter, values []string, options model.SnippetFormatOptions) string {
	if len(values) != len(parameters) {
		log.Warn().Msgf(
			"Number of parameters (%d) and number of supplied values (%d) does not match",
			len(parameters),
			len(values),
		)
		return snippet
	}

	var result string
	if options.ParamMode == model.SnippetParamModeSet {
		result = setParameters(snippet, parameters, values)
		if options.RemoveComments {
			result = pruneComments(result)
		}
	} else {
		result = replaceParameters(snippet, parameters, values)
	}

	return result
}

func setParameters(snippet string, parameters []model.Parameter, values []string) string {
	hints := parseHints(snippet)

	start := 0
	result := ""
	for i, parameter := range parameters {
		maxPosition := 0
		for _, hint := range hints {
			if hint.variable == parameter.Key {
				if hint.position > maxPosition {
					maxPosition = hint.position
				}
			}
		}

		newLine := fmt.Sprintf("%s=\"%s\"\n", parameter.Key, values[i])

		result += snippet[start:maxPosition] + newLine
		start = maxPosition
	}

	result += snippet[start:]

	return result
}

func replaceParameters(snippet string, parameters []model.Parameter, values []string) string {
	result := pruneComments(snippet)
	for i, parameter := range parameters {
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", parameter.Key), values[i])
	}
	return result
}

func hintsToParameters(hints []hint) []model.Parameter {
	var result []model.Parameter

	// put all variables into a list in order to their order of occurrence in the snippet
	var variableNames []string

	names := map[string]string{}
	descriptions := map[string]string{}
	defaults := map[string]string{}
	values := map[string][]string{}

	for _, h := range hints {
		variableNameExists := false
		for i := range variableNames {
			if varName := variableNames[i]; varName == h.variable {
				variableNameExists = true
				break
			}
		}

		if !variableNameExists {
			variableNames = append(variableNames, h.variable)
		}

		switch h.typeDescriptor {
		case hintTypeName:
			names[h.variable] = h.value
		case hintTypeDescription:
			descriptions[h.variable] = h.value
		case hintTypeDefaultValue:
			defaults[h.variable] = h.value
		case hintTypeValues:
			if parsedValues := stringutil.SplitWithEscape(h.value, ',', '\\', true); len(parsedValues) > 0 {
				if alreadyValues, ok := values[h.variable]; !ok {
					values[h.variable] = parsedValues
				} else {
					values[h.variable] = append(alreadyValues, parsedValues...)
				}
			}
		}
	}

	for _, varName := range variableNames {
		// If no name is provided, use the variable name as parameter name
		name, ok := names[varName]
		if !ok || name == "" {
			name = varName
		}

		result = append(result, model.Parameter{
			Key:          varName,
			Name:         name,
			Description:  descriptions[varName],
			DefaultValue: defaults[varName],
			Values:       values[varName],
		})
	}

	return result
}

func parseHints(snippet string) []hint {
	var result []hint

	scanner := bufio.NewScanner(strings.NewReader(snippet))
	position := 0
	for scanner.Scan() {
		line := scanner.Text()
		position += len(line) + 1

		currentHint := hint{
			typeDescriptor: hintTypeInvalid,
			position:       position,
		}

		match := hintRegex.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		for i, name := range hintRegex.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			if groupName, ok := toRegexNamedGroup(name); ok {
				switch groupName {
				case regexNamedGroupValue:
					currentHint.value = match[i]
				case regexNamedGroupVariable:
					currentHint.variable = match[i]
				case regexNamedGroupType:
					switch match[i] {
					case string(hintTypeName):
						currentHint.typeDescriptor = hintTypeName
					case string(hintTypeDescription):
						currentHint.typeDescriptor = hintTypeDescription
					case string(hintTypeDefaultValue):
						currentHint.typeDescriptor = hintTypeDefaultValue
					case string(hintTypeValues):
						currentHint.typeDescriptor = hintTypeValues
					}
				}
			}
		}

		if currentHint.isValid() {
			result = append(result, currentHint)
		}
	}

	return result
}

func toRegexNamedGroup(val string) (regexNamedGroup, bool) {
	switch val {
	case string(regexNamedGroupValue):
		return regexNamedGroupValue, true
	case string(regexNamedGroupVariable):
		return regexNamedGroupVariable, true
	case string(regexNamedGroupType):
		return regexNamedGroupType, true
	}
	return regexNamedGroupValue, false
}

func (h *hint) isValid() bool {
	return h.variable != "" && h.typeDescriptor != hintTypeInvalid && h.value != ""
}

func pruneComments(script string) string {
	scanner := bufio.NewScanner(strings.NewReader(script))
	result := ""
	for scanner.Scan() {
		line := scanner.Text()
		if hintRegex.MatchString(line) {
			continue
		}
		if result != "" {
			result += "\n"
		}
		result += line
	}
	return result
}
