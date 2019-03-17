package rules

import (
	"fmt"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"regexp"
	"strings"
)

var regexVersion3013 = regexp.MustCompile(`.+[==|@].+`)

// validateDL3013 Pin versions in pip. Instead of `pip install <package>` use `pip install <package>==<version>`
func validateDL3013(node *parser.Node, file string) (rst []string, err error) {
	for _, child := range node.Children {
		if child.Value == "run" {
			isPip, isInstall, length := false, false, len(rst)
			for _, v := range strings.Fields(child.Next.Value) {
				switch v {
				case "pip":
					isPip = true
				case "install":
					if isPip {
						isInstall = true
					}
				case "&&":
					isPip, isInstall = false, false
				default:
					if isPip && isInstall && !regexVersion3013.MatchString(v) && length == len(rst) {
						rst = append(rst, fmt.Sprintf("%s:%v DL3013 Pin versions in pip. Instead of `pip install <package>` use `pip install <package>==<version>`\n", file, child.StartLine))
					}
					isPip, isInstall = false, false
				}
			}
		}
	}
	return rst, nil
}
