package main

// The bot constants
const (
	Server  = "irc.freenode.net:6667"
	Channel = "#n0sec"
	BotName = "gob0tsec"
)

/*
   Arrays can't be consts. The [...] syntax lets the compiler find out how many elements
   are stored in the array (so that its size is fixed)
   Here the compiler will translate [...]string to [20]string
*/
var EightBallAnswers = [...]string{
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
}
