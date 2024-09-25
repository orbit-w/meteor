# Mailbox

The `mailbox` package provides an implementation of a mailbox system for actors, allowing for efficient message passing and processing. It supports both bounded and unbounded mailboxes, ensuring flexibility in handling different workloads.

## Features

- **Bounded and Unbounded Mailboxes**: Choose between bounded mailboxes, which block when full, and unbounded mailboxes, which can grow indefinitely.
- **Priority Queue**: System messages are prioritized over user messages.
- **Concurrency**: Safe to use in concurrent environments with multiple producers and a single consumer.
- **Graceful Suspension and Resumption**: Mailboxes can be suspended and resumed, allowing for controlled processing.

## Usage

### Creating a Mailbox

To create a new bounded mailbox:

```go
import (
    "github.com/orbit-w/meteor/modules/mailbox"
)

func main() {
    mb := mailbox.Bounded(100, 10)
    // Use the mailbox
}
```

### Pushing Messages

To push a message to the mailbox:

```go
mb.Push("Hello, World!")
```

To push a system message to the mailbox:

```go
mb.PushSystemMsg("System Message")
```

### Suspending and Resuming

To suspend the mailbox:

```go
mb.Suspend()
```

To resume the mailbox:

```go
mb.Resume()
```

## API

### `IMailbox` Interface

- `Push(msg any)`
    - Pushes a user message to the mailbox.
- `PushSystemMsg(msg any)`
    - Pushes a system message to the mailbox.
- `RegInvoker(_invoker Invoker)`
    - Registers an invoker to handle messages.
- `Suspend()`
    - Suspends the mailbox.
- `Resume()`
    - Resumes the mailbox.

### `MailBox` Struct

- `Bounded(size, processLimit int) IMailbox`
    - Creates a new bounded mailbox with the specified size and process limit.

## License

This project is licensed under the MIT License.
```