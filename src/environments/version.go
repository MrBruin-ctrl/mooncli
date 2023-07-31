package environments

import (
	"fmt"
	"github.com/mitchellh/colorstring"
)

func Logo() {
	logo := "\n" +
		"[dark_gray]_|      _|    _|_|      _|_|    _|      _|                _|_|_|  _|        _|_|_|\n" +
		"[dark_gray]_|_|  _|_|  _|    _|  _|    _|  _|_|    _|              _|        _|          _|\n" +
		"[white]_|  _|  _|  _|    _|  _|    _|  _|  _|  _|  _|_|_|_|_|  _|        _|          _|\n" +
		"[white]_|      _|  _|    _|  _|    _|  _|    _|_|              _|        _|          _|\n" +
		"[white]_|      _|    _|_|      _|_|    _|      _|                _|_|_|  _|_|_|_|    _|_|_|  \n[reset]"

	logo += fmt.Sprintf("%84v\n", "按键F1获取更多帮助信息")
	_, _ = colorstring.Println(logo)
}
