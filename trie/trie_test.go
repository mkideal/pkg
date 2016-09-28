package trie

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const wordSize = 1000

var (
	swords = []string{
		"敏感词一",
		"敏感词二",
		"敏感词三",
		"袁世凯",
		"孙中山",
		"毛泽东",
		"蒋介石",
		"袁大头",
		"孙小头",
		"毛老头",
		"老蒋",
	}
	dict = NewTrie(swords)

	charset = []string{
		"一", "二", "三",
		"袁", "世", "凯",
		"孙", "中", "山",
		"毛", "泽", "东",
		"蒋", "介", "石",
		"袁", "大", "头",
		"小", "老", "敏",
		"聖", "公", "會",
		"主", "風", "小",
		//"學", "的", "集",
		//"誦", "訓", "練",
		//"已", "有", "多",
		//"年", "的", "歷",
		//"史", "每", "年",
		//"都", "會", "挑",
		//"選", "三", "至",
		//"四", "年", "級",
		//"有", "潛", "質",
		//"的", "學", "生",
		"參", "加", "我",
	}
	specialCharset = []string{
		"a", "b", "c",
		"d", "e", "f",
		"!", "@", "#",
	}

	testwords = func() []string {
		rand.Seed(time.Now().UnixNano())
		ws := make([]string, 0, wordSize)
		for i := 0; i < wordSize; i++ {
			buf := new(bytes.Buffer)
			for n := 0; n < 30; n++ {
				buf.WriteString(charset[rand.Intn(len(charset))])
			}
			ws = append(ws, buf.String())
		}
		return ws
	}()
)

func TestMatch(t *testing.T) {
	word := "一孙中山四"
	node, deep := dict.match(word)
	t.Log(word, deep, dict.isWordTail(node))
}

func Benchmark_TrieMatch(b *testing.B) {
	count := 0
	for i := 0; i < b.N; i++ {
		for _, tc := range testwords {
			if dict.Match(tc) {
				count++
			}
		}
		if i == 0 {
			b.Logf("matched count: %d", count)
		}
	}
}

func Benchmark_ForMatch(b *testing.B) {
	count := 0
	for i := 0; i < b.N; i++ {
		for _, tc := range testwords {
			for _, w := range swords {
				if strings.Contains(tc, w) {
					//b.Log(tc, w)
					count++
					break
				}
			}
		}
		if i == 0 {
			b.Logf("matched count: %d", count)
		}
	}
}
