package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	commands = []command{
		{
			name:   "quit",
			action: quitAction,
			usage:  "退出程序",
		},
		{
			name:   "help",
			action: helpAction,
			usage:  "打印此帮助信息",
		},
		{
			name:   "width",
			action: widthAction,
			usage:  "[int] 查看或设置一行的宽度",
		},
		{
			name:   "calc",
			action: calcAction,
			usage:  "<int>... 计算并打印一行的填充方式",
		},
	}
}

const (
	setBlock   = '\u258A' // '▊'
	unsetBlock = '\u2591' // '░'
)

var (
	width                 int
	commands              []command
	invalidCommand        = command{name: "invalid"}
	noCommandError        = errors.New("没有这个命令")
	ambiguousCommandError = errors.New("命令有歧义")
	quitWithoutError      = errors.New("退出")
	invalidFormatError    = errors.New("非法格式")
	invalidLengthError    = errors.New("长度不是正数")
	invalidNumberError    = errors.New("非法的数字")
)

type command struct {
	name   string
	action func(args []string) error
	usage  string
}

func exec(command string, args []string) error {
	cmd, err := getCommand(command)
	if err != nil {
		return err
	}
	return cmd.action(args)
}

func getCommand(cmdName string) (command, error) {
	var matches []command
	for _, cmd := range commands {
		if strings.HasPrefix(cmd.name, cmdName) {
			matches = append(matches, cmd)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	} else if len(matches) > 1 {
		return invalidCommand, ambiguousCommandError
	} else {
		return invalidCommand, noCommandError
	}
}

func quitAction(_ []string) error {
	return quitWithoutError
}

func helpAction(_ []string) error {
	for _, cmd := range commands {
		fmt.Println(cmd.name, cmd.usage)
	}
	return nil
}

func widthAction(args []string) error {
	if len(args) == 0 {
		fmt.Println(width)
		return nil
	}
	if len(args) > 1 {
		return invalidFormatError
	}

	newWidth, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	if newWidth < 0 {
		return invalidLengthError
	}

	width = newWidth
	return nil
}

func calcAction(args []string) error {
	if len(args) < 1 {
		return invalidFormatError
	}

	nums, err := parseNumbers(args)
	if err != nil {
		return err
	}
	sumNum, err := sum(nums)
	if err != nil {
		return err
	}

	diff := width - sumNum
	row := makeRow()
	setRow(row, nums, diff)
	printRow(row)
	return nil
}

func parseNumbers(args []string) ([]int, error) {
	result := make([]int, len(args))
	for i, arg := range args {
		num, err := strconv.Atoi(arg)
		if err != nil {
			return nil, err
		}
		result[i] = num
	}
	return result, nil
}

func sum(nums []int) (int, error) {
	result := 0
	for i, num := range nums {
		if i > 0 {
			result++
		}
		result += num
	}
	if result > width {
		return 0, invalidNumberError
	} else {
		return result, nil
	}
}

func makeRow() []rune {
	row := make([]rune, width)
	for i := range row {
		row[i] = unsetBlock
	}
	return row
}

func setRow(row []rune, nums []int, diff int) {
	i := 0
	for j, num := range nums {
		if j > 0 {
			i++
		}
		i += num
		setCount := num - diff
		for k := i - setCount; k <= i-1; k++ {
			row[k] = setBlock
		}
	}
}

func printRow(row []rune) {
	for i, block := range row {
		if i > 0 && i%5 == 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%c", block)
	}
	fmt.Println()
}
