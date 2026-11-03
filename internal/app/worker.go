package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

const demoTopic = "demo"

// subscribeDemo 订阅 demo 主题。
func subscribeDemo(ctx context.Context, pubSub *gochannel.GoChannel) (<-chan *message.Message, error) {
	messages, err := pubSub.Subscribe(ctx, demoTopic)
	if err != nil {
		return nil, fmt.Errorf("subscribe %s: %w", demoTopic, err)
	}
	slog.InfoContext(ctx, "worker subscribed", "topic", demoTopic)
	return messages, nil
}

// consumeDemo 消费订阅通道中的消息，直至 ctx 取消或通道关闭。
func consumeDemo(ctx context.Context, messages <-chan *message.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-messages:
			if !ok {
				return
			}
			handleDemoMessage(ctx, msg)
		}
	}
}

func handleDemoMessage(ctx context.Context, msg *message.Message) {
	slog.InfoContext(ctx, "worker received", "uuid", msg.UUID, "payload", string(msg.Payload))
	msg.Ack()
}
