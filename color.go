package main

func blue(str string) string {
	return "\033[1;34m" + str + "\033[0m"
}

func yellowWithBlueBG(str string) string {
	return "\033[1;33;44m" + str + "\033[0m"
}
