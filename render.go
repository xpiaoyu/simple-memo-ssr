package main

import (
	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
	"strings"
)

var test = `# Java 笔记

### 遍历某个类中的所有方法是否有特定的注解

参数 | 示例值 | 参数描述
--- | --- | ---
productId | IF | 产品ID 如: IF, IC, IH, TS, TF, T
date| 2019-7-5 | 交易日期

` + "```java" + `
public static List<Method> getMethodsAnnotatedWithMethodXY(final Class<?> type) {
    final List<Method> methods = new ArrayList<Method>();
    Class<?> klass = type;
    while (klass != Object.class) { // need to iterated thought hierarchy in order to retrieve methods from above the current instance
        // iterate though the list of methods declared in the class represented by klass variable, and add those annotated with the specified annotation
        final List<Method> allMethods = new ArrayList<Method>(Arrays.asList(klass.getDeclaredMethods()));
        for (final Method method : allMethods) {
            if (method.isAnnotationPresent(MethodXY.class)) {
                MethodXY annotInstance = method.getAnnotation(MethodXY.class);
                if (annotInstance.x() == 3 && annotInstance.y() == 2) {         
                    methods.add(method);
                }
            }
        }
        // move to the upper class in the hierarchy in search for more methods
        klass = klass.getSuperclass();
    }
    return methods;
}
` + "```"

var ignoreTags = []string{"script", "title"}

var ignoreTagsRune = [][]rune{[]rune("script"), []rune("title")}

func MarkdownToHtml(md string) string {
	md = strings.Replace(md, "\r", "", -1)
	return string(blackfriday.Run([]byte(md), blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.ChromaStyle(myGitHub)))))
}

func HighlightKeyword(html string, key string) string {
	htmlRune := []rune(html)
	keyRune := []rune(key)
	htmlLen := len(htmlRune)
	keyLen := len(keyRune)
	/*STATE
	  0: Inside angle bracket
	  1: Outside angle bracket
	  2: Inside special char e.g. &lt;*/
	state := 0
	matchLen := 0
	flagBegin := []rune("<span style=\"color:red;background:yellow;\">")
	flagEnd := []rune("</span>")
	beginLen := len(flagBegin)
	_ = beginLen
	//flagEndLen := len(flagEnd)
	flagLen := len(flagBegin) + len(flagEnd)
	_ = flagLen
	ignoreTagLevel := 0
	var tag = make([]rune, 0, 128)
	var tagLen int
	for i := 0; i < htmlLen; i++ {
		switch state {
		case 0:
			if htmlRune[i] == '&' {
				state = 2
			}
			if htmlRune[i] == '<' {
				state = 1
			} else if runeEqualIgnoreCase(htmlRune[i], keyRune[matchLen]) {
				matchLen++
				if matchLen >= keyLen {
					// Find match!
					matchLen = 0
					htmlRune = insertRuneSliceAt(htmlRune, flagBegin, flagEnd, i-keyLen+1, i+1)
					i += flagLen
					htmlLen += flagLen
				}
			} else {
				matchLen = 0
			}
		case 1:
			if htmlRune[i] == '>' {
				//tagName := string(tag[:tagLen])
				tagName := tag[:tagLen]
				//log.Println("tag-name:", tagName)
				if enterIgnoreTagRune(tagName) {
					ignoreTagLevel++
				} else if exitIgnoreTagRune(tagName) {
					if ignoreTagLevel > 0 {
						ignoreTagLevel--
					}
				} else {
					state = 0
				}
				tagLen = 0
			} else {
				if len(tag) > tagLen {
					tag[tagLen] = htmlRune[i]
				} else {
					tag = append(tag, htmlRune[i])
				}
				tagLen++
			}
		case 2:
			if htmlRune[i] == ';' {
				state = 0
			}
		}
	}
	return string(htmlRune)
}

func enterIgnoreTagRune(tagName []rune) bool {
	runeToLower(tagName)
	for k := range ignoreTagsRune {
		tag := ignoreTagsRune[k]
		if runeHasPrefix(tagName, tag) {
			return true
		}
	}
	return false
}

func exitIgnoreTagRune(tagName []rune) bool {
	if len(tagName) <= 0 {
		return false
	}
	runeToLower(tagName)
	if tagName[0] != '/' {
		return false
	}
	for k := range ignoreTagsRune {
		tag := ignoreTagsRune[k]
		if runeHasPrefix(tagName[1:], tag) {
			return true
		}
	}
	return false
}

func runeHasPrefix(s, prefix []rune) bool {
	if len(s) < len(prefix) {
		return false
	}
	prefixLen := len(prefix)
	for i := 0; i < prefixLen; i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

func enterIgnoreTag(tagName string) bool {
	tagName = strings.ToLower(tagName)
	for _, v := range ignoreTags {
		if strings.HasPrefix(tagName, v) {
			return true
		}
	}
	return false
}

func exitIgnoreTag(tagName string) bool {
	if len(tagName) <= 0 {
		return false
	}
	tagName = strings.ToLower(tagName)
	if tagName[0] != '/' {
		return false
	}
	for _, v := range ignoreTags {
		if strings.HasPrefix(tagName, "/"+v) {
			return true
		}
	}
	return false
}

func insertRuneSliceAt(dst []rune, src []rune, src2 []rune, index int, index2 int) (ret []rune) {
	len1 := len(src)
	len2 := len(src2)
	lenDst := len(dst)
	ret = dst
	ret = append(ret, src...)
	ret = append(ret, src2...)
	copy(ret[index2+len1+len2:], dst[index2:lenDst])
	copy(ret[index+len1:], dst[index:index2])
	copy(ret, dst[:index])
	copy(ret[index:], src)
	copy(ret[index2+len1:], src2)
	return
}

func runeToLower(src []rune) {
	for k := range src {
		if 65 <= src[k] && src[k] <= 90 {
			src[k] = src[k] + 32
		}
	}
}

func runeEqualIgnoreCase(a, b rune) bool {
	if a == b {
		return true
	}
	if 97 <= a && a <= 122 {
		a -= 32
	}
	if 97 <= b && b <= 122 {
		b -= 32
	}
	return a == b
}
