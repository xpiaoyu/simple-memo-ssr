package main

import (
	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/styles"
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

func MarkdownToHtml(md string) string {
	md = strings.Replace(md, "\r", "", -1)
	return string(blackfriday.Run([]byte(md), blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.ChromaStyle(styles.Dracula)))))
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
	flagBegin := []rune("<b style=\"color:red;background:yellow;\">")
	flagEnd := []rune("</b>")
	beginLen := len(flagBegin)
	_ = beginLen
	//flagEndLen := len(flagEnd)
	flagLen := len(flagBegin) + len(flagEnd)
	_ = flagLen
	ignoreTagLevel := 0
	var tag []rune
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
					//htmlRune = insertRuneSliceAt(htmlRune, flagBegin, i-keyLen+1)
					//log.Println("html:", string(htmlRune))
					//htmlRune = insertRuneSliceAt(htmlRune, flagEnd, i+beginLen+1)
					//log.Println("html:", string(htmlRune))
					htmlRune = insertRuneSliceAt(htmlRune, flagBegin, flagEnd, i-keyLen+1, i+1)
					i += flagLen
					htmlLen += flagLen
				}
			} else {
				matchLen = 0
			}
		case 1:
			if htmlRune[i] == '>' {
				tagName := string(tag)
				//log.Println("tag-name:", tagName)
				if enterIgnoreTag(tagName) {
					ignoreTagLevel++
				} else if exitIgnoreTag(tagName) {
					if ignoreTagLevel > 0 {
						ignoreTagLevel--
					}
				} else {
					state = 0
				}
				tag = nil
			} else {
				tag = append(tag, htmlRune[i])
			}
		case 2:
			if htmlRune[i] == ';' {
				state = 0
			}
		}
	}
	return string(htmlRune)
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
	//prefix := dst[:index]
	//suffix := dst[index:]
	//ret = append(ret, dst[:index]...)
	//ret = append(ret, src...)
	//ret = append(ret, dst[index:]...)
	//log.Println("cap:", cap(dst), "len:", len(dst), "index:", index, "index2", index2)
	len1 := len(src)
	len2 := len(src2)
	//if cap(dst)-len(dst) >= len1+len2 {
	//	ret = dst
	//} else {
	//	ret = make([]rune, len(dst)+len1+len2)
	//}
	//ret = make([]rune, len(dst)+len1+len2)
	//ret = dst
	lenDst := len(dst)
	//ret = dst
	//ret = append(ret, src...)
	//ret = append(ret, src2...)
	ret = make([]rune, len(dst)+len1+len2)
	copy(ret, dst[:index])
	copy(ret[index+len1:], dst[index:index2])
	copy(ret[index2+len1+len2:], dst[index2:lenDst])
	copy(ret[index:], src)
	copy(ret[index2+len1:], src2)
	return
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
