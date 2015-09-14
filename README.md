[![Go Report Card](http://goreportcard.com/badge/Depado/go-b0tsec)](http://goreportcard.com/report/Depado/go-b0tsec)

# Go-b0tsec

A pretty simple IRC Bot with plugins and middlewares.

## Plugins

A plugin is a command. End of definition. Each plugin is mapped with a command.  
To create your own plugin, take a look at the [plugins/commands/examples](https://github.com/Depado/go-b0tsec/tree/master/plugins/commands/example) directory to see an example.

To write a plugin you need to satisfy the Plugin interface which is defined like that :

```go
// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	Get(*irc.Connection, string, string, []string)
	Help(*irc.Connection, string)
}
```

I'd advise you to create a new package for each plugin and always call your plugin struct "Plugin", so that the API can remain stable and uniform.  
You can then register your plugin like that :

```go
RegisterCommand("command", new(myplugin.Plugin))
```

## Middlewares

A middleware is a function that is executed each time a message is sent, no matter what. The term "middleware" may not be the smartest choice there, but it kind of felt like a middleware so...  
I see middlewares as a way to monitor and react to things that are being said on a channel without the users specifically calling for the bot. As an example you can check the github middleware that will send the description of a github repository each time it finds a github link in a message.

To write your own middleware you need to satisfy the Middleware interface which is defined like that :

```go
// Middleware represents a single middleware.
type Middleware interface {
	Get(*irc.Connection, string, string, string)
}
```

As for the plugins, I'd advise you to create a package per middleware, and call the middleware struct "Middleware" to keep the API stable and uniform.  
You can then register your middleware like that :

```go
RegisterMiddleware(new(mymiddleware.Middleware))
```
