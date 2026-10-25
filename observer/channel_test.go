package observer_test

import (
	"testing"
	"time"

	"github.com/matthiasharzer/service-kit/observer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestSubscriber struct {
	Notifications []string
}

func (s *TestSubscriber) OnMessage(notification string) {
	s.Notifications = append(s.Notifications, notification)
}

func TestNewSyncChannel(t *testing.T) {
	t.Run("sends notifications to a subscriber", func(t *testing.T) {
		testChannel := observer.NewSyncChannel[string]()
		testSubscriber := &TestSubscriber{}
		testChannel.Subscribe(testSubscriber.OnMessage)

		testChannel.Publish("Test Notification 1")
		testChannel.Publish("Test Notification 2")

		require.Len(t, testSubscriber.Notifications, 2)
		assert.Equal(t, "Test Notification 1", testSubscriber.Notifications[0])
		assert.Equal(t, "Test Notification 2", testSubscriber.Notifications[1])
	})

	t.Run("does not send notifications to unsubscribed subscribers", func(t *testing.T) {
		testChannel := observer.NewSyncChannel[string]()
		testSubscriber := &TestSubscriber{}
		unsubscribe := testChannel.Subscribe(testSubscriber.OnMessage)
		unsubscribe()

		testChannel.Publish("Test Notification 1")

		assert.Len(t, testSubscriber.Notifications, 0)
	})

	t.Run("handles multiple subscribers", func(t *testing.T) {
		testChannel := observer.NewSyncChannel[string]()
		subscriber1 := &TestSubscriber{}
		subscriber2 := &TestSubscriber{}
		testChannel.Subscribe(subscriber1.OnMessage)
		testChannel.Subscribe(subscriber2.OnMessage)

		testChannel.Publish("Test Notification")

		require.Len(t, subscriber1.Notifications, 1)
		require.Len(t, subscriber2.Notifications, 1)
		assert.Equal(t, "Test Notification", subscriber1.Notifications[0])
		assert.Equal(t, "Test Notification", subscriber2.Notifications[0])
	})

	t.Run("clears all subscribers", func(t *testing.T) {
		testChannel := observer.NewSyncChannel[string]()
		subscriber1 := &TestSubscriber{}
		subscriber2 := &TestSubscriber{}
		testChannel.Subscribe(subscriber1.OnMessage)
		testChannel.Subscribe(subscriber2.OnMessage)

		testChannel.Clear()

		testChannel.Publish("Test Notification")

		assert.Len(t, subscriber1.Notifications, 0)
		assert.Len(t, subscriber2.Notifications, 0)
	})

	t.Run("runs subscribers synchronously", func(t *testing.T) {
		testChannel := observer.NewSyncChannel[string]()

		receivedAt := map[string]time.Time{}

		testChannel.Subscribe(func(s string) {
			receivedAt["subscriber1"] = time.Now()

			// this blocks the second channel
			time.Sleep(100 * time.Millisecond)
		})
		testChannel.Subscribe(func(s string) {
			receivedAt["subscriber2"] = time.Now()
		})
		testChannel.Publish("test")

		require.Len(t, receivedAt, 2)
		require.Contains(t, receivedAt, "subscriber1")
		require.Contains(t, receivedAt, "subscriber2")

		firstReceived := receivedAt["subscriber1"]
		secondReceived := receivedAt["subscriber2"]

		assert.Less(t, firstReceived, secondReceived)
		delta := secondReceived.Sub(firstReceived)
		assert.GreaterOrEqual(t, delta, 100*time.Millisecond)
	})
}

func TestNewAsyncChannel(t *testing.T) {
	t.Run("runs subscribers asynchronously", func(t *testing.T) {
		testChannel := observer.NewAsyncChannel[string]()

		receivedAt := map[string]time.Time{}

		testChannel.Subscribe(func(s string) {
			receivedAt["subscriber1"] = time.Now()

			// this should not block the second channel
			time.Sleep(100 * time.Millisecond)
		})
		testChannel.Subscribe(func(s string) {
			receivedAt["subscriber2"] = time.Now()
		})
		testChannel.Publish("test")

		// we need to wait for the subscribers to finish
		time.Sleep(200 * time.Millisecond)

		require.Len(t, receivedAt, 2)
		require.Contains(t, receivedAt, "subscriber1")
		require.Contains(t, receivedAt, "subscriber2")

		firstReceived := receivedAt["subscriber1"]
		secondReceived := receivedAt["subscriber2"]

		assert.WithinDuration(t, firstReceived, secondReceived, 20*time.Millisecond)
	})
}
