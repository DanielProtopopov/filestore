package redis

import (
	"context"
	"fmt"
	"log"
	"strings"
	// Package dependencies
	"github.com/redis/go-redis/v9"
)

// Constants
var constConfigKeyspaceEvents = "notify-keyspace-events"
var constKeyspaceFormat = "__keyevent@%d__"
var constKeyspacePatternFormat = "%s:%s"
var constDefaultKeyspacePattern = "*"

/*
RedisEventMsg is a message struct emitted onto the sink channel for a consumer
to process.
*/
type RedisEventMsg struct {
	// The redis key event, i.e. 'SET', 'EXPIRE' etc
	KeyEvent string
	// The key to which the event applies
	Key string
}

/*
RedisEventSink monitors events emitted from Redis data store.

NOTE: Events will not be emitted unless the Redis instance has the required
config entry set in the "redis.conf" file.

To do this, you can either:
 1. Ensure the following entry exists in your `redis.conf`:
    notify-keyspace-events "AKE"
 2. Set the entry via `CONFIG SET`
 3. Use this event sink code to automatically check/set the config

For more detailed information about Redis keyspace notifications please
see http://redis.io/topics/notifications
*/
type RedisEventSink struct {
	redisClient     *redis.Client
	keyspace        string
	pattern         string
	keyspacePattern string
	subscriber      *redis.PubSub
	EventChannel    chan *RedisEventMsg
}

/*
NewRedisEventSink is a factory function that creates a new instance of the event sink
*/
func NewRedisEventSink(client *redis.Client, databaseIndex int64, pattern string) *RedisEventSink {

	if client == nil {
		log.Panic("[RedisEventSink] Invalid redis client passed")
	}

	if databaseIndex < 0 || databaseIndex > 16 {
		log.Panic("[RedisEventSink] Invalid database index specified")
	}

	if pattern == "" {
		pattern = constDefaultKeyspacePattern
	}

	var keyspace = fmt.Sprintf(constKeyspaceFormat, databaseIndex)

	self := &RedisEventSink{
		redisClient:     client,
		keyspace:        keyspace,
		pattern:         pattern,
		keyspacePattern: fmt.Sprintf(constKeyspacePatternFormat, keyspace, pattern),
		EventChannel:    make(chan *RedisEventMsg),
	}

	self.init()
	return self
}

func (sink *RedisEventSink) init() {
	var setConfigKeyspaceEvents = false

	// Check to ensure keyspace events are enabled...
	configResultCmd := sink.redisClient.ConfigGet(context.Background(), constConfigKeyspaceEvents)
	if configResultCmd.Err() != nil {
		log.Panic(configResultCmd.Err())
	}

	if len(configResultCmd.Val()) > 0 {
		for _, configValue := range configResultCmd.Val() {
			// Returned config string MUST contain at least a "K" or "E" in it
			if !strings.ContainsAny(configValue, "KE") {
				log.Printf("[RedisEventSink] Keyspace event notifications (\"%s\") are not enabled correctly in Redis database config.", constConfigKeyspaceEvents)
				setConfigKeyspaceEvents = true
			}
		}
	} else {
		log.Printf("[RedisEventSink] Keyspace event notifications (\"%s\") are not enabled in Redis database config.", constConfigKeyspaceEvents)
		setConfigKeyspaceEvents = true
	}

	if setConfigKeyspaceEvents {
		log.Printf("[RedisEventSink] Updating config for keyspace event notifications.")
		configStatusCmd := sink.redisClient.ConfigSet(context.Background(), constConfigKeyspaceEvents, "AKE")
		if configStatusCmd.Err() != nil {
			log.Panic(configStatusCmd.Err())
		}
	}

	// Subscribe to events matching the supplied keyspace pattern
	sink.subscriber = sink.redisClient.PSubscribe(context.Background(), sink.keyspacePattern)

	log.Printf("[RedisEventSink] Subscribed to events - (Keyspace Pattern: \"%s\")", sink.keyspacePattern)

	// Start observing for events...
	go sink.listenForEvents()
}

/*
Pattern returns the pattern that this sink is observing
*/
func (sink *RedisEventSink) Pattern() string {
	return sink.pattern
}

/*
Close performs component clean up
*/
func (sink *RedisEventSink) Close() {
	if sink.subscriber != nil {
		err := sink.subscriber.Close()
		if err != nil {
			log.Printf("[RedisEventSink] Error closing events subscription: %#v", err)
		}
	}
}

/*
listenForEvents receives events and pushes them back out on a channel
NOTE: Its easy to receive vast quantities of events, so consider commenting
out logging or making it conditional.
*/
func (sink *RedisEventSink) listenForEvents() {
	for {
		if sink.subscriber != nil {
			msgi, err := sink.subscriber.Receive(context.Background())
			if err != nil {
				log.Printf("[RedisEventSink] Error receiving events: %v", err)
			} else {
				// log.Infof("Received Redis Event: %#v", msgi)

				switch msg := msgi.(type) {
				case *redis.Subscription:
					log.Printf("[RedisEventSink] Event subscribed: \"%s\", Data:\"%s\"", msg.Kind, msg.Channel)
				case *redis.Message:
					eventType := strings.Replace(msg.Channel, sink.keyspace+":", "", -1)

					keyEventMsg := &RedisEventMsg{
						KeyEvent: eventType,
						Key:      msg.Payload,
					}
					sink.EventChannel <- keyEventMsg

				default:
				}
			}
		}
	}
}
