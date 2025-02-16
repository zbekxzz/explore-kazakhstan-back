package main

import "auth/internal/utils"

func main() {
	d := utils.Info{"jack"}

	d.SendMail()
}
