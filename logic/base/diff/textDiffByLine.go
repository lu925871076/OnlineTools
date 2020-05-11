package diff

import (
	"com.lu/OnlineTools/common"
	"strings"
)

//文本按行diff，逻辑同精确diff，区别在于单位由一个字符到一行
type TextDiffByLine struct {
	//初始化时字符串按'\n'切分
	text1 []string
	text2 []string
	dp    [][]int8
	path  [][]string
	diffs []TextExactDiffsData
}

//初始化，按换行切分，dp，path切片初始化
func (t *TextDiffByLine) init(text1, text2 string) {
	t.text1 = strings.Split(text1, "\n")
	t.text2 = strings.Split(text2, "\n")
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
func (t *TextDiffByLine) diff() {
	t.getPath()
	t.GetDiffFromPath()
	//按行diff，不要合并
	//t.mergeDiffs()
}

//获取html显示代码
func (t *TextDiffByLine) getHtml() string {
	diffs := t.diffs
	html := ""
	for _, diff := range diffs {
		op := diff.op
		text := diff.text
		switch op {
		case DIFF_INSERT:
			html += "<ins style=\"background:#e6ffe6;\">" + text + "</ins><br/>"
		case DIFF_DELETE:
			html += "<del style=\"background:#ffe6e6;\">" + text + "</del><br/>"
		case DIFF_EQUAL:
			html += "<span>" + text + "</span><br/>"
		}
	}
	return html
}

//lcs算法生成path
func (t *TextDiffByLine) getPath() {
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
func (t *TextDiffByLine) GetDiffFromPath() {
	text1 := t.text1
	text2 := t.text2
	t.pathDiffs(len(text1), len(text2))
}

//递归生成diffs
func (t *TextDiffByLine) pathDiffs(i, j int) {
	path := t.path
	text1 := t.text1
	text2 := t.text2

	if i == 0 || j == 0 {
		for i != 0 {
			t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_DELETE, text1[i-1]))
			i--
		}
		for j != 0 {
			t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_INSERT, text2[j-1]))
			j--
		}
		return
	}
	if path[i][j] == DIFF_PATH_DIAGONA {
		t.pathDiffs(i-1, j-1)
		t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_EQUAL, text1[i-1]))
	} else if path[i][j] == DIFF_PATH_VERTICAL {
		t.pathDiffs(i-1, j)
		t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_DELETE, text1[i-1]))
	} else if path[i][j] == DIFF_PATH_HORIZONTAL {
		t.pathDiffs(i, j-1)
		t.diffs = append(t.diffs, NewTextExactDiffData(DIFF_INSERT, text2[j-1]))
	}

}

func (t *TextDiffByLine) mergeDiffs() {
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

func newTextDiffByLine(text1, text2 string) *TextDiffByLine {
	t := &TextDiffByLine{}
	t.init(text1, text2)
	return t
}

func DoTextDiffByLine(text1, text2 string) string {
	t := newTextDiffByLine(text1, text2)
	t.diff()
	return t.getHtml()
}
