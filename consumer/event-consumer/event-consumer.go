package event_consumer

import (
	"links_tg-bot/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

const numGoroutines = 100

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() (err error) {
	eventsChan := make(chan []events.Event)

	for i := 0; i < numGoroutines; i++ {
		go c.handleEvents(eventsChan)
	}

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)

		if err != nil {
			log.Printf("[ERR] in consumer: %s\n", err.Error())

			continue
		}

		if len(gotEvents) == 0 {

			time.Sleep(1 * time.Second)

			continue
		}

		eventsChan <- gotEvents

		// Code without goroutines
		//if err := c.handleEvents(gotEvents); err != nil {
		//	log.Print(err)
		//
		//	continue
		//}
	}

}

func (c Consumer) handleEvents(ch <-chan []events.Event) {

	for events := range ch {

		for _, event := range events {
			log.Printf("got new event: %s\n", event.Text)

			if err := c.processor.Process(event); err != nil {
				log.Printf("can't handle event: %s\n", err.Error())

				continue
			}
		}
	}
}

// Code without goroutines
//func (c Consumer) handleEvents(events []events.Event) error {
//
//	for _, event := range events {
//		log.Printf("got new event: %s\n", event.Text)
//
//		if err := c.processor.Process(event); err != nil {
//			log.Printf("can't handle event: %s\n", err.Error())
//
//			continue
//		}
//	}
//
//	return nil
//}
