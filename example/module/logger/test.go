package logger

var log Logger

func Do() {
	a := "123"
	//log := logger.NewConsoleLogger("debug")
	log = NewFileLogger("debug",
		"./logger/log", "test.log", 10*1024)
	for {
		log.Debug("hello world a:%s", a)
		log.Trace("hello world")
		log.Info("hello world")
		log.Warn("hello world")
		log.Error("hello world")
		log.Fatal("hello world")
	}
}
