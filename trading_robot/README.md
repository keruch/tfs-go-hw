# Golang cryptocurrency trading bot

This project contains the implementation of a cryptocurrency trading bot on [Kraken](https://futures.kraken.com/) exchange, written in Golang.


##Setup config
Before running, you must fill the `config.toml` due to example `config/config.example.toml`. This file should be placed in `config` directory.
You can create private and public keys in your Kraken Futures account.
Read how to do it [here](https://support.kraken.com/hc/en-us/articles/360022839451-Generate-API-keys).
You can read how to get Telegram token [here](https://core.telegram.org/bots). 


##Working with bot
To launch the bot, run
```
go run ./cmd/traging_robot
```

Once launched, you should subscribe to a pair. You can do it with an HTTP request:
```
POST <address>/subscribe/<ticker>
```
`<address>` must be specified in the config file, `<ticker>` can be obtained from [Kraken support](https://support.kraken.com/hc/en-us/articles/360022835891-Ticker-symbols).
Example request: 
```
POST localhost:8091/subscribe/PI_ETHUSD
```
Also, you can unsubscribe pair in the same way:
```
POST localhost:8091/unsubscribe/PI_ETHUSD
```

So far, you cannot change the pair for trading in runtime, but in this way you can stop the bot for a while, and then, if necessary, start it again without terminating the process.

**Attention!** Do not even try to change pair in runtime. This can lead to undefined behavior!

You can also change the settings of orders at runtime. Available settings: position size, position price multiplier (for successful execution of ioc orders with low liquidity).
You can do it with:
```
POST <address>/quantity/<value>
POST <address>/multiplier/<value>
```
`quantity` must be a positive integer, `multiplier` is a positive floating point number that modifies your price using the formula:
```
price = price * (1 + multiplier),    buy  case
price = price * (1 - multiplier),    sell case 
```

Bot can be gracefully terminated with the SIGHUP, SIGINT, SIGTERM, and SIGQUIT signals.

##Endpoints list:
```
POST /subscribe/<ticker>
POST /unsubscribe/<ticker>
POST /quantity/<value>
POST /multiplier/<value>
```