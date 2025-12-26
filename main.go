package main

import "stock-investing/cmd/stock-investing"

func main() {
	// cmd/stock-investing/main.go 의 main이 아니라,
	// 여기서는 별도 패키지로 뺀 엔트리 함수를 호출하는 식으로 구성하는 게 일반적이다.
	// 단순화를 위해 지금은 실제 로직을 cmd에 두고,
	// 나중에 리팩터링할 수 있다.
	stock_investing.Main()
}
