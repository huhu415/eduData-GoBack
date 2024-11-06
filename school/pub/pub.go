// Package pub 提供不同学校都可以用的一些公共函数
package pub

import (
	"container/list"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type StuType int

const (
	// UG 本科生
	UG StuType = 1
	// PG 研究生
	PG StuType = 2
)

type SchoolName string

const (
	// HRBUST 哈尔滨理工大学
	HRBUST SchoolName = "hrbust"
	// HLJU 黑龙江大学
	HLJU SchoolName = "hlju"
	// NEAU 东北农业大学
	NEAU SchoolName = "neau"
)

// ExtractWeekRange 提取周数范围, 返回起始周, 结束周, 单双周, 单双周默认是5, 0是双周, 1是单周, 便于运算
func ExtractWeekRange(text string) (startWeek, endWeek, evenOrOdd int, err error) {
	startWeek, endWeek, evenOrOdd = 0, 0, 5
	if []rune(text)[0] == '第' {
		// 匹配形式 : 第3周
		// 正则表达式提取数字
		WeekSinge := regexp.MustCompile("[0-9]+").FindAllString(text, 1)
		var atoi int
		atoi, err = strconv.Atoi(WeekSinge[0])
		if err != nil {
			err = fmt.Errorf("形式 第x周 无法解析起始周: %v", err)
			return
		}
		// 如果是这种形式, 那么起始周与结束周都是这个数字
		startWeek, endWeek = atoi, atoi
		// fmt.Println("起始周-结束周", startWeek, "-", endWeek)
	} else {
		// 匹配形式 : 1-15周 或 1-15单周
		matchWeekRange := regexp.MustCompile(`(\d+)-(\d+)`).FindStringSubmatch(text)
		if len(matchWeekRange) == 0 {
			err = errors.New("无法匹配周数范围, x-x")
			return
		}
		// 第一个数
		startWeek, err = strconv.Atoi(matchWeekRange[1])
		if err != nil {
			err = fmt.Errorf("形式 x-x周 无法解析起始周: %v", err)
			return
		}
		// 第二个数
		endWeek, err = strconv.Atoi(matchWeekRange[2])
		if err != nil {
			err = fmt.Errorf("形式 x-x周 无法解析结束周: %v", err)
			return
		}
		// 判断单双周
		switch {
		case strings.Contains(text, "单"):
			evenOrOdd = 1
		case strings.Contains(text, "双"):
			evenOrOdd = 0
		}
		// fmt.Println("起始周-结束周", startWeek, "-", endWeek)
	}
	return
}

// 提取第几节课
func ExtractCoruse(text string) (start, end int, err error) {
	// 匹配形式 : 1-15周 或 1-15单周
	matchWeekRange := regexp.MustCompile(`(\d+)-(\d+)`).FindStringSubmatch(text)
	if len(matchWeekRange) == 0 {
		err = errors.New("无法匹配周数范围, x-x")
		return
	}
	// 第一个数
	start, err = strconv.Atoi(matchWeekRange[1])
	if err != nil {
		err = fmt.Errorf("形式 x-x周 无法解析起始周: %v", err)
		return
	}
	// 第二个数
	end, err = strconv.Atoi(matchWeekRange[2])
	if err != nil {
		err = fmt.Errorf("形式 x-x周 无法解析结束周: %v", err)
		return
	}

	return
}

// NewColorList 初始化颜色队列, 用于给课程上色, 一共19个, 应该用不完
func NewColorList() *list.List {
	queue := list.New()
	queue.PushBack("#2e1f54")
	queue.PushBack("#52057f")
	queue.PushBack("#bf033b")
	queue.PushBack("#f00a36")
	queue.PushBack("#ff6908")
	queue.PushBack("#ffc719")
	queue.PushBack("#598c14")
	queue.PushBack("#335238")
	queue.PushBack("#4a8594")
	queue.PushBack("#051736")
	queue.PushBack("#706357")
	queue.PushBack("#b0a696")
	queue.PushBack("#004eaf")
	queue.PushBack("#444444")
	queue.PushBack("#c1d1e0")
	queue.PushBack("#c1d1e0")
	queue.PushBack("#faa918")
	queue.PushBack("#8f1010")
	queue.PushBack("#d2ea32")
	return queue
}

// ChineseToNumber 汉字数字转阿拉伯数字
func ChineseToNumber(chinese rune) (int, error) {
	switch chinese {
	case '一':
		return 1, nil
	case '二':
		return 2, nil
	case '三':
		return 3, nil
	case '四':
		return 4, nil
	case '五':
		return 5, nil
	case '六':
		return 6, nil
	case '日':
		return 7, nil
	case '天':
		return 7, nil
	default:
		return 0, fmt.Errorf("不支持的汉字数字: %s", string(chinese))
	}
}
