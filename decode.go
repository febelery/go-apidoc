package goapidoc

import (
	"fmt"
	"regexp"
	"strings"
)

func unmarshal(data []string) (doc docDef) {
	for _, datum := range data {
		parts := strings.SplitN(datum, " ", 2)
		tag := strings.TrimPrefix(parts[0], "@")

		if len(parts) < 2 {
			continue
		}

		switch tag {
		case "api":
			doc.Api = parseApi(parts[1])
		case "apiVersion":
			doc.ApiVersion = strings.TrimSpace(parts[1])
		case "apiName":
			doc.ApiName = strings.TrimSpace(parts[1])
		case "apiGroup":
			doc.ApiGroup = strings.TrimSpace(parts[1])
		case "apiPrivate":
			doc.ApiPrivate = strings.TrimSpace(parts[1]) == "true"
		case "apiPermission":
			doc.ApiPermission = strings.TrimSpace(parts[1])
		case "apiDescription":
			doc.ApiDescription = strings.TrimSpace(parts[1])
		case "apiUse":
			doc.ApiUse = append(doc.ApiUse, strings.TrimSpace(parts[1]))
		case "apiDefine":
			doc.ApiDefine = append(doc.ApiDefine, strings.TrimSpace(parts[1]))
		case "apiParam":
			doc.ApiParam = append(doc.ApiParam, parserParam(parts[1]))
		case "apiHeader":
			doc.ApiHeader = append(doc.ApiHeader, parserParam(parts[1]))
		case "apiQuery":
			doc.ApiQuery = append(doc.ApiQuery, parserParam(parts[1]))
		case "apiBody":
			doc.ApiParam = append(doc.ApiParam, parserParam(parts[1]))
		case "apiSuccess":
			doc.ApiSuccess = append(doc.ApiSuccess, parserParam(parts[1]))
		case "apiSuccessExample":
			doc.ApiSuccessExample = append(doc.ApiSuccessExample, parserRespExample(strings.TrimSpace(parts[1])))
		case "apiError":
			doc.ApiError = append(doc.ApiError, parserParam(parts[1]))
		case "apiErrorExample":
			doc.ApiErrorExample = append(doc.ApiErrorExample, parserRespExample(strings.TrimSpace(parts[1])))
		}
	}

	return
}

func parseApi(str string) (api routeDef) {
	str = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(str), " ")
	base := strings.SplitN(str, " ", 3)
	if len(base) > 0 {
		base[0] = strings.Replace(base[0], "{", "", -1)
		base[0] = strings.Replace(base[0], "}", "", -1)
		api.Method = strings.TrimSpace(base[0])
	}
	if len(base) > 1 {
		api.Path = strings.TrimSpace(base[1])
	}
	if len(base) > 2 {
		api.Title = strings.TrimSpace(base[2])
	}
	return
}

func parserParam(str string) (param paramDef) {
	str = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(str), " ")
	segments := strings.SplitN(str, " ", 3)

	if len(segments) > 0 {
		// {String{2..5}="xxx",aaa}
		matches := regexp.MustCompile(`{([^{}=]+)(?:{([^{}=]+)})?(?:=([^{}=]+(?:,[^{}=]+)*)?)?}`).FindStringSubmatch(segments[0])

		if len(matches) == 4 {
			if len(matches[1]) > 0 {
				param.Type.Type = matches[1]
			}

			if len(matches[2]) > 0 {
				param.Type.Size = matches[2]
			}

			if len(matches[3]) > 0 {
				allows := strings.Split(strings.Replace(matches[3], "\"", "", -1), ",")
				if len(allows) > 0 {
					param.Type.Allows = allows
				}
			}
		}
	}

	if len(segments) > 1 {
		param.Field.Required = !(strings.Contains(segments[1], "[") && strings.Contains(segments[1], "]"))

		// [address.street=ZZZ]
		matches := regexp.MustCompile(`\[?(\w+(?:\.\w+)*)(?:=([^]]+))?]?`).FindStringSubmatch(segments[1])
		fmt.Println(segments[1], matches)
		if len(matches) > 2 {
			param.Field.Name = matches[1]
		}
		if len(matches) >= 3 {
			param.Field.Default = matches[2]
		}
	}

	if len(segments) > 2 {
		param.Description = segments[2]
	}

	return
}

func parserRespExample(str string) (resp respExampleDef) {
	segments := strings.SplitN(str, "\n", 3)

	if len(segments) > 0 {
		prefixSeg := strings.SplitN(segments[0], " ", 2)

		if len(prefixSeg) > 0 {
			prefixSeg[0] = strings.Replace(prefixSeg[0], "{", "", -1)
			prefixSeg[0] = strings.Replace(prefixSeg[0], "}", "", -1)
			resp.Type = strings.TrimSpace(prefixSeg[0])
		}

		if len(prefixSeg) > 0 {
			resp.Title = prefixSeg[1]
		}

		resp.Http = segments[1]
	}
	if len(segments) > 1 {
		resp.Http = segments[1]
	}
	if len(segments) > 2 {
		resp.Example = segments[2]
	}

	return
}
