package diff

import (
	"com.lu/OnlineTools/common"
	"fmt"
)

//diff接口
type BaseDiff interface {
	init(string, string)
	diff()
	getHtml() string
}

//文本精确diff
type TextExactDiff struct {
	//text1 []string
	//text2 []string
	text1 string
	text2 string
	dp    [][]int8
	path  [][]string
	diffs []TextExactDiffsData
}

type TextExactDiffsData struct {
	op   int
	text string
}

//初始化，dp，path切片初始化
func (t *TextExactDiff) init(text1, text2 string) {
	//t.text1 = strings.Split(text1, "\n")
	//t.text2 = strings.Split(text2, "\n")
	t.text1 = text1
	t.text2 = text2
	//for i := 0; i < len(t.text1); i++ {
	//	t.text1[i] += "\n"
	//}
	//for i := 0; i < len(t.text2); i++ {
	//	t.text2[i] += "\n"
	//}
	t.diffs = make([]TextExactDiffsData, 0)
	len1 := len(t.text1)
	len2 := len(t.text2)
	dp := make([][]int8, 2)
	t.dp = dp
	for i := 0; i < 2; i++ {
		dp[i] = make([]int8, len2+1)
		for j := 0; j <= len2; j++ {
			dp[i][j] = 0
		}
	}

	path := make([][]string, len1+1)
	t.path = path
	for i := 0; i <= len1; i++ {
		path[i] = make([]string, len2+1)
		for j := 0; j <= len2; j++ {
			path[i] = append(path[i], "")
		}
	}
}

//diff主流程
func (t *TextExactDiff) diff() {
	t.getPath()
	t.GetDiffFromPath()
	t.mergeDiffs()
}

//获取html显示代码
func (t *TextExactDiff) getHtml() string {
	diffs := t.diffs
	html := ""
	for _, diff := range diffs {
		op := diff.op
		text := diff.text
		switch op {
		case DIFF_INSERT:
			html += "<ins style=\"background:#e6ffe6;\">" + text + "</ins>"
		case DIFF_DELETE:
			html += "<del style=\"background:#ffe6e6;\">" + text + "</del>"
		case DIFF_EQUAL:
			html += "<span>" + text + "</span>"
		}
	}
	return html
}

//lcs算法生成path
func (t *TextExactDiff) getPath() {
	k := 1
	text1 := t.text1
	text2 := t.text2
	len1 := len(text1)
	len2 := len(text2)
	dp := t.dp
	path := t.path
	for i := 1; i <= len1; i++ {
		oppose := common.GetOppose(k)
		for j := 1; j <= len2; j++ {
			if text1[i-1] == text2[j-1] {
				dp[k][j] = dp[oppose][j-1] + 1
				path[i][j] = DIFF_PATH_DIAGONA
			} else if dp[oppose][j] < dp[k][j-1] {
				dp[k][j] = dp[k][j-1]
				path[i][j] = DIFF_PATH_HORIZONTAL
			} else {
				dp[k][j] = dp[oppose][j]
				path[i][j] = DIFF_PATH_VERTICAL
			}
		}
		k = oppose
	}
}

//根据path生成diffs
func (t *TextExactDiff) GetDiffFromPath() {
	text1 := t.text1
	text2 := t.text2
	t.pathDiffs(len(text1), len(text2))
	//todo
	fmt.Println(t.diffs)
}

//递归生成diffs
func (t *TextExactDiff) pathDiffs(i, j int) {
	path := t.path
	text1 := t.text1
	text2 := t.text2

	if i == 0 || j == 0 {
		for i != 0 {
			t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_DELETE, string(text1[i-1])))
			i--
		}
		for j != 0 {
			t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_INSERT, string(text2[j-1])))
			j--
		}
		return
	}
	if path[i][j] == DIFF_PATH_DIAGONA {
		t.pathDiffs(i-1, j-1)
		t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_EQUAL, string(text1[i-1])))
	} else if path[i][j] == DIFF_PATH_VERTICAL {
		t.pathDiffs(i-1, j)
		t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_DELETE, string(text1[i-1])))
	} else if path[i][j] == DIFF_PATH_HORIZONTAL {
		t.pathDiffs(i, j-1)
		t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_INSERT, string(text2[j-1])))
	}

}

func (t *TextExactDiff) mergeDiffs() {
	diffs := t.diffs
	if len(diffs) == 0 {
		return
	}
	ret := make([]TextExactDiffsData, 0, 0)
	for i := 0; i < len(diffs); {
		firstOp := diffs[i].op
		firstText := diffs[i].text
		j := i + 1
		for j < len(diffs) && diffs[j].op == firstOp {
			firstText += diffs[j].text
			j++
		}
		i = j
		ret = append(ret, NewTextExactDiffData(firstOp, firstText))
		if i >= len(diffs) {
			break
		}
	}
	t.diffs = ret
}

func NewTextExactDiffData(op int, text string) TextExactDiffsData {
	return TextExactDiffsData{
		op:   op,
		text: text,
	}
}

//LCS
func (t *TextExactDiff) getLCS() string {
	return t.pathLcs(len(t.text1), len(t.text2))
}

func (t *TextExactDiff) pathLcs(i, j int) string {
	if i == 0 || j == 0 {
		return ""
	}
	lcs := ""
	path := t.path
	text := t.text1
	if path[i][j] == DIFF_PATH_DIAGONA {
		t.pathLcs(i-1, j-1)
		lcs += string(text[i-1])
	} else if path[i][j] == DIFF_PATH_VERTICAL {
		t.pathLcs(i-1, j)
	} else if path[i][j] == DIFF_PATH_HORIZONTAL {
		t.pathLcs(i, j-1)
	}
	return lcs
}

func newTextExactDiff(text1, text2 string) *TextExactDiff {
	t := &TextExactDiff{}
	t.init(text1, text2)
	return t
}

func DoTextExactDiff(text1, text2 string) string {
	t := newTextExactDiff(text1, text2)
	t.diff()
	return t.getHtml()
}
