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
	hintParamType      string
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
	hintTypeParamType    = hintTypeDescriptor("Type")
	hintTypeValues       = hintTypeDescriptor("Values")
	hintTypeInvalid      = hintTypeDescriptor("invalid")

	regexNamedGroupVariable = regexNamedGroup("varname")
	regexNamedGroupType     = regexNamedGroup("key")
	regexNamedGroupValue    = regexNamedGroup("value")

	paramTypePath     = hintParamType("PATH")
	paramTypePassword = hintParamType("PASSWORD")
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

	allHintValues := toHintValues(hints)

	for _, varName := range allHintValues.variableNames {
		// If no name is provided, use the variable name as parameter name
		name, ok := allHintValues.names[varName]
		if !ok || name == "" {
			name = varName
		}

		result = append(result, model.Parameter{
			Key:          varName,
			Name:         name,
			Description:  allHintValues.descriptions[varName],
			DefaultValue: allHintValues.defaults[varName],
			Values:       allHintValues.values[varName],
			Type:         mapToParameterType(allHintValues.types[varName]),
		})
	}

	return result
}

type hintValues struct {
	variableNames []string
	names         map[string]string
	descriptions  map[string]string
	defaults      map[string]string
	values        map[string][]string
	types         map[string]string
}

func toHintValues(hints []hint) hintValues {
	result := hintValues{
		names:        map[string]string{},
		descriptions: map[string]string{},
		defaults:     map[string]string{},
		values:       map[string][]string{},
		types:        map[string]string{},
	}

	for _, h := range hints {
		variableNameExists := false
		for i := range result.variableNames {
			if varName := result.variableNames[i]; varName == h.variable {
				variableNameExists = true
				break
			}
		}

		if !variableNameExists {
			result.variableNames = append(result.variableNames, h.variable)
		}

		switch h.typeDescriptor {
		case hintTypeName:
			result.names[h.variable] = h.value
		case hintTypeDescription:
			result.descriptions[h.variable] = h.value
		case hintTypeDefaultValue:
			result.defaults[h.variable] = h.value
		case hintTypeParamType:
			result.types[h.variable] = h.value

		case hintTypeValues:
			if parsedValues := stringutil.SplitWithEscape(h.value, ',', '\\', true); len(parsedValues) > 0 {
				if alreadyValues, ok := result.values[h.variable]; !ok {
					result.values[h.variable] = parsedValues
				} else {
					result.values[h.variable] = append(alreadyValues, parsedValues...)
				}
			}
		}
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
					currentHint.typeDescriptor = hintTypeDescriptor(match[i])
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

func mapToParameterType(val string) model.ParameterType {
	switch val {
	case string(paramTypePath):
		return model.ParameterTypePath
	case string(paramTypePassword):
		return model.ParameterTypePassword
	}
	return model.ParameterTypeValue
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
