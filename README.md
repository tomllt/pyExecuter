# pyExecuter

pyExecuter is a high-performance Python script execution program based on [devchat-ai/gopool](https://github.com/devchat-ai/gopool). It leverages Go's parallel processing capabilities and Python's flexibility to provide a robust, efficient, and scalable solution for executing Python scripts.

## Features

### 1. Parallel Task Scheduling
- Utilizes devchat-ai/gopool's Goroutines pool for efficient concurrent task management
- Implements intelligent task scheduling based on CPU and memory resources
- Dynamically allocates Python script tasks to the worker pool for maximum resource utilization

### 2. Task Queue Management
- Supports priority-based task queues for urgent task execution
- Offers both FIFO and LIFO queue modes to suit different task processing needs
- Includes automatic retry mechanism for failed tasks

### 3. Multi-threading and Multi-processing Options
- Provides flexible Python script execution modes, supporting both multi-threading and multi-processing
- Automatically selects the appropriate execution mode based on task type (e.g., CPU-intensive vs I/O-intensive)

### 4. Load Balancing and Resource Management
- Implements efficient load balancing to distribute system resources evenly
- Dynamically adjusts concurrent task numbers based on current system resource usage
- Utilizes devchat-ai/gopool's built-in mechanisms to manage Goroutines within reasonable limits

### 5. Task Monitoring and Logging
- Offers detailed task monitoring for real-time tracking of Python script execution status, runtime, and resource consumption
- Records comprehensive task execution logs, including start/end times, output, and errors
- Provides a web-based dashboard for visualizing task status and resource utilization

### 6. Fault Handling and Recovery Mechanisms
- Automatically captures exceptions and records detailed error information for failed tasks
- Supports configurable task retries with customizable retry counts and intervals
- Implements task recovery mechanisms to resume execution from previous states after unexpected crashes

### 7. Multi-node Support and Distributed Execution
- Supports distributed environments for executing tasks across multiple nodes
- Provides inter-node communication and synchronization mechanisms
- Utilizes lightweight RPC or message queues (e.g., Redis, RabbitMQ) for task distribution and result collection

### 8. Task Timeout Control
- Sets execution timeout for each Python task to prevent resource hogging
- Automatically terminates and cleans up resources for timed-out tasks

### 9. Result Collection and Processing
- Automatically collects and stores task results in structured formats (e.g., databases or files)
- Offers a unified result aggregation interface for easy retrieval and processing
- Supports callback functions for custom result processing logic

### 10. High-performance Python Script Execution Interface
- Provides a flexible Python execution interface with dynamic parameter passing
- Supports dynamic execution of Python code passed as strings
- Integrates with Python virtual environments for dependency isolation

### 11. REST API and WebSocket Support
- Offers REST API for task submission, status querying, and result retrieval
- Implements WebSocket for real-time task status monitoring

### 12. Configurability and Extensibility
- Provides configurable execution parameters (e.g., concurrency level, max resource usage, log levels)
- Supports plugin-based extensions for adding new task types, monitoring tools, etc.

### 13. Security Controls
- Executes tasks in isolated environments to prevent system threats from malicious scripts
- Utilizes sandboxing technologies (e.g., Docker or virtual environments) for task isolation
- Implements authentication mechanisms to ensure only authorized users can submit and manage tasks

### 14. Data Persistence and Task State Management
- Persists task states (pending, running, completed, failed) for recovery after program restarts
- Supports periodic saving of intermediate states for long-running tasks

### 15. Flexible Task Dependency Management
- Supports task dependency definitions
- Manages task dependencies using Directed Acyclic Graphs (DAG) for efficient, ordered execution

## Installation

[Instructions on how to install pyExecuter]

## Usage

[Basic usage examples and code snippets]

## Configuration

[Details on how to configure pyExecuter]

## Contributing

We welcome contributions to pyExecuter! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) file for details on how to get started.

## License

[Your chosen license]

## Contact

[Your contact information or project maintainer details]

## Acknowledgements

- [devchat-ai/gopool](https://github.com/devchat-ai/gopool) for providing the foundational Goroutines pool implementation.
- [Any other libraries or resources you want to acknowledge]