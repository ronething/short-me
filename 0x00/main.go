package main

func main() {
	var (
		a *App
	)
	a = &App{
	}
	a.Initialize()
	a.Run(":8888")

}
