# Computer Science Culture

## 1. What is a reentrant function?

A reentrant function is a function that can be interrupted in the middle of its execution and safely called again ("re-entered") before the previous executions are complete. This property is crucial in multi-threaded or interrupt-driven environments where a function might be preempted and called again before it finishes its initial execution.

### Key Characteristics of Reentrant Functions:

1. **No Static or Global State**: Reentrant functions do not rely on static or global variables that can be modified by concurrent executions. Instead, they use local variables or pass all necessary state information as parameters.

2. **Atomic Operations**: If a reentrant function needs to perform operations that could be interrupted, it must ensure those operations are atomic (i.e., indivisible and uninterruptible).

3. **Thread Safety**: While all reentrant functions are thread-safe, not all thread-safe functions are reentrant. Thread safety generally means that a function can be safely called from multiple threads, while reentrancy specifically addresses the ability to interrupt and re-enter the function.

### Example:

A simple example of a reentrant function is one that uses only local variables:

```go
func ReentrantFunction(a int) int {
    b := a + 1
    return b
}
```

In contrast, a non-reentrant function might look like this:

```go
var b int
func NonReentrantFunction(a int) int {
    b = a + 1
    return b
}
```

In the non-reentrant example, the global variable `b` could lead to incorrect behavior if the function is interrupted and called again before the first call completes.


## 2. What is the difference between a thread, a fork, and a coroutine?

1. **Thread**:
   - A thread is the smallest unit of processing that can be scheduled by an operating system. Threads within the same process share the same memory space, which allows for efficient communication and data sharing.
   - However, this shared memory can also lead to issues like race conditions if not managed properly.
   - Threads are typically used for parallelism, allowing multiple tasks to be executed simultaneously.

2. **Fork**:
   - Forking is a method of creating a new process by duplicating an existing one. The new process (child) is an exact copy of the original (parent) but has its own memory space.
   - Forking is often used in Unix-like operating systems to create new processes. Each process runs independently, and communication between processes typically requires inter-process communication (IPC) mechanisms.
   - Forking can be resource-intensive, as it involves copying the entire memory space of the parent process.

3. **Coroutine**:
   - A coroutine is a general control structure whereby flow control is cooperatively passed between two or more routines without returning to the caller.
   - Coroutines are lighter than threads and can be used for asynchronous programming. They allow for non-blocking operations, making them ideal for I/O-bound tasks.
   - Unlike threads, coroutines run in a single thread and share the same memory space, which makes them more efficient in terms of resource usage.


## 3. What is HELM?

HELM can also mean a Kubernetes Helm chart, which is a package manager for Kubernetes applications. It is like `apt` or `homebrew` and helps manage Kubernetes applications by providing a way to define, install, and upgrade even the most complex Kubernetes applications.

## 4. What do you think about design patterns? (give examples) Have you ever used one or more design patterns?

Design patterns are proven solutions to common software design problems. They provide a template for how to solve a problem in a way that is reusable and adaptable to different situations. Some well-known design patterns include:

1. **Singleton**: Ensures that a class has only one instance and provides a global point of access to it. This is useful for managing shared resources, such as database connections or configuration settings.

2. **Observer**: Defines a one-to-many dependency between objects so that when one object changes state, all its dependents are notified and updated automatically. This is commonly used in event-driven systems.

3. **Factory Method**: Provides an interface for creating objects in a superclass but allows subclasses to alter the type of objects that will be created. This promotes loose coupling and adherence to the Open/Closed Principle.

I have used several design patterns in my projects, particularly the Singleton and Observer patterns. For example, I implemented the Observer pattern in a real-time chat application to update users' interfaces when new messages arrive.

## 5. What is NewSQL?

NewSQL is a class of modern relational database management systems that aims to provide the best of both worlds: the scalability and high performance of NoSQL databases combined with the transactional consistency (ACID guarantees) of traditional SQL databases.

Key Characteristics:
 - SQL as the Primary Interface
 - ACID Guarantees for Transactions
 - Horizontal Scalability
 - Distributed, Fault-Tolerant Architecture
 - Modern, Cloud-Native Design

Example:
 - Google Cloud Spanner
 - CockroachDB
 - tiDB