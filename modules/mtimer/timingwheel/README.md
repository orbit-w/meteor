# Scheduler

Scheduler is a Go package that manages the scheduling of timers using a Hierarchical Time Wheel. It provides an efficient way to handle delayed tasks with high precision and low overhead.

## Features

- Add tasks with a specified delay.
- Start and stop the scheduler.
- Graceful shutdown ensuring all pending tasks are executed.
- High precision timing with millisecond granularity.

## Usage

### Creating a Scheduler

To create a new Scheduler instance:

```go
import (
    "github.com/orbit-w/meteor/modules/mtimer/timingwheel"
)

func main() {
    scheduler := timingwheel.NewScheduler()
    scheduler.Start()
    defer scheduler.GracefulStop(context.Background())
}
```

### Adding Tasks

To add a new task to the scheduler with a specified delay:

```go
scheduler.Add(5*time.Second, func(args ...any) {
    fmt.Println("Task executed")
})
```

### Stopping the Scheduler

To stop the scheduler immediately:

```go
scheduler.Stop()
```

To stop the scheduler gracefully, ensuring all pending tasks are executed:

```go
scheduler.GracefulStop(context.Background())
```

## API

### `IScheduler` Interface

- `Add(delay time.Duration, callback func(...any), args ...any) *TimerTask`
    - Adds a new task to the scheduler with the given delay and callback.
    - Parameters:
        - `delay`: The duration to wait before executing the callback.
        - `callback`: The function to execute after the delay.
        - `args`: Additional arguments to pass to the callback function.
    - Returns: A pointer to the created `TimerTask`.

- `Start()`
    - Initiates the Scheduler, starting the internal processes required for scheduling tasks.

- `GracefulStop(ctx context.Context) error`
    - Stops the Scheduler gracefully, ensuring all pending timers are executed before stopping.
    - Parameters:
        - `ctx`: A context used to control the timeout for stopping the Scheduler.
    - Returns: An error if the Scheduler fails to close within the context's timeout.

- `Stop()`
    - Stops the Scheduler immediately.

### `Scheduler` Struct

- `NewScheduler() *Scheduler`
    - Creates a new Scheduler instance.

- `Add(delay time.Duration, callback func(...any), args ...any) *TimerTask`
    - Adds a new task to the scheduler with the given delay and callback.

- `Start()`
    - Starts the Scheduler, initiating the ticking.

- `GracefulStop(ctx context.Context) error`
    - Stops the Scheduler gracefully, ensuring all pending timers are executed before stopping.

- `Stop()`
    - Stops the Scheduler immediately.

## License

This project is licensed under the MIT License.
```