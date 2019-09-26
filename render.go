package main

import (
	"bytes"
	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
)

var ignoreTagsByte = [][]byte{[]byte("script"), []byte("title")}

var ignoreTagsRune = [][]rune{[]rune("script"), []rune("title")}

func MarkdownToHtml(md []byte) []byte {
	//md = strings.Replace(md, "\r", "", -1)
	md = bytes.Replace(md, []byte("\r"), nil, -1)
	baseRenderer := blackfriday.NewHTMLRenderer(
		blackfriday.HTMLRendererParameters{
			Flags:                      blackfriday.FootnoteReturnLinks | blackfriday.CommonHTMLFlags,
			FootnoteReturnLinkContents: "<sup>返回</sup>",
		})
	chromaRenderer := bfchroma.NewRenderer(bfchroma.ChromaStyle(myGitHub), bfchroma.Extend(baseRenderer))
	return blackfriday.Run(
		md,
		blackfriday.WithRenderer(baseRenderer),
		blackfriday.WithExtensions(blackfriday.Footnotes|blackfriday.AutoHeadingIDs|blackfriday.CommonExtensions),
		blackfriday.WithRenderer(chromaRenderer),
	)
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
	//return string(htmlRune)
	return "123"
}

func HighlightKeywordBytes(html []byte, key []byte) []byte {
	htmlData := make([]byte, len(html))
	copy(htmlData, html)
	keyRune := key
	htmlLen := len(htmlData)
	keyLen := len(keyRune)
	/*STATE
	  0: Outside angle bracket
	  1: Inside angle bracket
	  2: Inside special char e.g. &lt;*/
	state := 0
	matchLen := 0
	flagBegin := []byte("<span style=\"color:red;background:yellow;\">")
	flagEnd := []byte("</span>")
	beginLen := len(flagBegin)
	_ = beginLen
	//flagEndLen := len(flagEnd)
	flagLen := len(flagBegin) + len(flagEnd)
	_ = flagLen
	ignoreTagLevel := 0
	var tag = make([]byte, 0, 128)
	var tagLen int
	for i := 0; i < htmlLen; i++ {
		switch state {
		case 0:
			if htmlData[i] == '&' {
				state = 2
			}
			if htmlData[i] == '<' {
				state = 1
			} else if byteEqualIgnoreCase(htmlData[i], keyRune[matchLen]) {
				matchLen++
				if matchLen >= keyLen {
					// Find match!
					matchLen = 0
					htmlData = insertRuneSliceAtBytes(htmlData, flagBegin, flagEnd, i-keyLen+1, i+1)
					i += flagLen
					htmlLen += flagLen
				}
			} else {
				matchLen = 0
			}
		case 1:
			if htmlData[i] == '>' {
				tagName := tag[:tagLen]
				//log.Println("tag-name:", string(tagName))
				if enterIgnoreTagByte(tagName) {
					ignoreTagLevel++
				} else if exitIgnoreTagByte(tagName) {
					if ignoreTagLevel > 0 {
						ignoreTagLevel--
					}
				} else {
					state = 0
				}
				tagLen = 0
			} else {
				if len(tag) > tagLen {
					tag[tagLen] = htmlData[i]
				} else {
					tag = append(tag, htmlData[i])
				}
				tagLen++
			}
		case 2:
			if htmlData[i] == ';' {
				state = 0
			}
		}
	}
	return htmlData
}

func enterIgnoreTagByte(tagName []byte) bool {
	byteToLower(tagName)
	for k := range ignoreTagsByte {
		tag := ignoreTagsByte[k]
		if bytes.HasPrefix(tagName, tag) {
			return true
		}
	}
	return false
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

func exitIgnoreTagByte(tagName []byte) bool {
	if len(tagName) <= 0 {
		return false
	}
	byteToLower(tagName)
	if tagName[0] != '/' {
		return false
	}
	for k := range ignoreTagsByte {
		tag := ignoreTagsByte[k]
		if bytes.HasPrefix(tagName[1:], tag) {
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

func insertRuneSliceAtBytes(dst []byte, src []byte, src2 []byte, index int, index2 int) (ret []byte) {
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

func byteToLower(src []byte) {
	for k := range src {
		if 65 <= src[k] && src[k] <= 90 {
			src[k] = src[k] + 32
		}
	}
}

func runeToLower(src []rune) {
	for k := range src {
		if 65 <= src[k] && src[k] <= 90 {
			src[k] = src[k] + 32
		}
	}
}

func byteEqualIgnoreCase(a, b byte) bool {
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
