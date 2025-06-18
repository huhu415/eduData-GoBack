package pub

import "testing"

func TestFullWidthToHalfWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "全角引号",
			input:    "“中国近现代史”经典文献导读",
			expected: "\"中国近现代史\"经典文献导读",
		},
		{
			name:     "全角圆括号",
			input:    "【25春期末考试（120分钟）考试】中国古代史（Ⅱ）",
			expected: "[25春期末考试(120分钟)考试]中国古代史(Ⅱ)",
		},
		{
			name:     "全角方括号",
			input:    "【】",
			expected: "[]",
		},
		{
			name:     "全角数字",
			input:    "０１２３４５６７８９",
			expected: "0123456789",
		},
		{
			name:     "全角字母",
			input:    "ＡＢＣａｂｃ",
			expected: "ABCabc",
		},
		{
			name:     "全角标点符号",
			input:    "！＂＃％＆＇＊＋，－．／：；＜＝＞？＠",
			expected: "!\"#%&'*+,-./:;<=>?@",
		},
		{
			name:     "全角空格",
			input:    " 　 ",
			expected: "   ",
		},
		{
			name:     "混合字符串",
			input:    "这是一个测试（Ｔｅｓｔ）１２３",
			expected: "这是一个测试(Test)123",
		},
		{
			name:     "纯中文字符",
			input:    "中文字符不变",
			expected: "中文字符不变",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "已经是半角的字符",
			input:    "Hello World 123!",
			expected: "Hello World 123!",
		},
		{
			name:     "全角下划线和波浪号",
			input:    "＿～",
			expected: "_~",
		},
		{
			name:     "全角大括号",
			input:    "｛123｝",
			expected: "{123}",
		},
		{
			name:     "课程表常见字符",
			input:    "数学（１－２节）",
			expected: "数学(1-2节)",
		},
		{
			name:     "时间表示",
			input:    "８：００－９：４０",
			expected: "8:00-9:40",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FullWidthToHalfWidth(tt.input)
			if result != tt.expected {
				t.Errorf("FullWidthToHalfWidth(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
