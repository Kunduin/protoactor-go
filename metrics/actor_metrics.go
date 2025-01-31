// Copyright (C) 2017 - 2022 Asynkron.se <http://www.asynkron.se>

package metrics

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/log"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
)

const LibName string = "protoactor"

type ActorMetrics struct {
	// Mutual Exclusion Primitive to use with ActorMailboxLength
	mu *sync.Mutex

	// MetricsID
	ID string

	// Actors
	ActorFailureCount            asyncint64.Counter
	ActorMailboxLength           asyncint64.Gauge
	ActorMessageReceiveHistogram syncfloat64.Histogram
	ActorRestartedCount          asyncint64.Counter
	ActorSpawnCount              asyncint64.Counter
	ActorStoppedCount            asyncint64.Counter

	// Deadletters
	DeadLetterCount       asyncint64.Counter
	FuturesCompletedCount asyncint64.Counter

	// Futures
	FuturesStartedCount  asyncint64.Counter
	FuturesTimedOutCount asyncint64.Counter

	// Threadpool
	ThreadPoolLatency syncint64.Histogram
}

// NewActorMetrics creates a new ActorMetrics value and returns a pointer to it
func NewActorMetrics() *ActorMetrics {
	instruments := newInstruments()
	return instruments
}

// newInstruments will create instruments using a meter from
// the given provider p
func newInstruments() *ActorMetrics {
	meter := global.Meter(LibName)
	instruments := ActorMetrics{mu: &sync.Mutex{}}

	var err error
	if instruments.ActorFailureCount, err = meter.AsyncInt64().Counter(
		"protoactor_actor_failure_count",
		instrument.WithDescription("Number of actor failures"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create ActorFailureCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorMessageReceiveHistogram, err = meter.SyncFloat64().Histogram(
		"protoactor_actor_message_receive_duration_seconds",
		instrument.WithDescription("Actor's messages received duration in seconds"),
	); err != nil {
		err = fmt.Errorf("failed to create ActorMessageReceiveHistogram instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorRestartedCount, err = meter.AsyncInt64().Counter(
		"protoactor_actor_restarted_count",
		instrument.WithDescription("Number of actors restarts"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create ActorRestartedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorStoppedCount, err = meter.AsyncInt64().Counter(
		"protoactor_actor_stopped_count",
		instrument.WithDescription("Number of actors stopped"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create ActorStoppedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorSpawnCount, err = meter.AsyncInt64().Counter(
		"protoactor_actor_spawn_count",
		instrument.WithDescription("Number of actors spawn"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create ActorSpawnCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.DeadLetterCount, err = meter.AsyncInt64().Counter(
		"protoactor_deadletter_count",
		instrument.WithDescription("Number of deadletters"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create DeadLetterCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.FuturesCompletedCount, err = meter.AsyncInt64().Counter(
		"protoactor_futures_completed_count",
		instrument.WithDescription("Number of futures completed"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create FuturesCompletedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.FuturesStartedCount, err = meter.AsyncInt64().Counter(
		"protoactor_futures_started_count",
		instrument.WithDescription("Number of futures started"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create FuturesStartedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.FuturesTimedOutCount, err = meter.AsyncInt64().Counter(
		"protoactor_futures_timed_out_count",
		instrument.WithDescription("Number of futures timed out"),
		instrument.WithUnit(unit.Dimensionless),
	); err != nil {
		err = fmt.Errorf("failed to create FuturesTimedOutCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ThreadPoolLatency, err = meter.SyncInt64().Histogram(
		"protoactor_thread_pool_latency_duration_seconds",
		instrument.WithDescription("History of latency in second"),
		instrument.WithUnit(unit.Milliseconds),
	); err != nil {
		err = fmt.Errorf("failed to create ThreadPoolLatency instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	return &instruments
}

// SetActorMailboxLengthGauge makes sure access to ActorMailboxLength is sequenced
func (am *ActorMetrics) SetActorMailboxLengthGauge(gauge asyncint64.Gauge) {
	// lock our mutex
	am.mu.Lock()
	defer am.mu.Unlock()

	am.ActorMailboxLength = gauge
}
