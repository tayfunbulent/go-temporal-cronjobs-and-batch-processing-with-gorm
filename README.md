# Temporal Examples with Go

## Objective

“This document provides a step-by-step guide to using Temporal with Go, exploring its key features and demonstrating how it can be leveraged in a project.”

## Why Temporal?

“Temporal provides a robust solution for managing distributed, long-running workflows, with built-in support for retries, state persistence, and fault tolerance.”

## Prerequisites

“Before getting started, ensure you have Docker, Go, and the Temporal CLI installed.”

## Running Temporal

“Use the Docker Compose configuration in the project file `engine` to run the Temporal server and related services.”

## Temporal Key Concepts

### Workflow:
- **Explanation:** A Workflow in Temporal represents a series of coordinated steps or activities that define a business process. It is durable, can run for a long time, and can handle retries, timeouts, and signals.
- **Role:** The Workflow manages the overall process, coordinating the execution of various activities and ensuring the correct sequence and handling of errors or interruptions.

### Activity:
- **Explanation:** Activities are the tasks or operations that a Workflow executes. They represent the individual units of work, such as sending an email, processing a payment, or generating a report.
- **Interaction with Workflow:** Activities are invoked by Workflows, and the results or outputs of these activities are used by the Workflow to determine the next steps.

### Task Queue:
- **Explanation:** Task Queues are used in Temporal to manage the distribution of tasks (like activities) to workers. A Task Queue holds tasks until a worker is available to execute them.
- **Configuration:** Task Queues can be configured to handle specific types of tasks, balancing the load across multiple workers or machines.

### Signals & Timers:
- **Signals Explanation:** Signals are a way to send external input to a running Workflow, allowing dynamic interaction with the Workflow while it is in progress.
- **Timers Explanation:** Timers allow a Workflow to wait or delay execution for a specific period, which is useful for scheduling activities or handling timeouts.

## Advanced Usage

### Child Workflows:
- **Explanation:** Child Workflows are Workflows that are initiated and managed by a parent Workflow. They allow the breakdown of complex processes into smaller, more manageable units.
- **Management:** Parent Workflows can start, wait for, or handle the results of Child Workflows, enabling a hierarchical workflow structure.

### Cron Jobs:
- **Explanation:** Cron Jobs in Temporal allow you to schedule recurring Workflows at specified intervals, similar to traditional cron jobs in Unix/Linux systems.
- **Setup:** Temporal provides a flexible scheduling mechanism where Workflows can be automatically triggered based on a predefined schedule.

### Parallel Processing:
- **Explanation:** Parallel Processing in Temporal refers to the ability to execute multiple activities or Workflows concurrently, allowing for more efficient processing of large tasks or datasets.
- **Usage:** This is particularly useful for tasks that can be broken down into independent units of work that do not need to be processed sequentially.

## Results and Learnings

“Using Temporal, we were able to simplify the management of distributed workflows and improve the reliability of our processes.”

## Future Work

“In the future, we plan to explore more advanced features of Temporal, such as custom search attributes and further optimization.”

## Project Structure
```bash
go-temporal-cronjobs-and-batch-processing-with-gorm/
├── activities/
│   ├── sandwichShop.go
│   ├── sendEmail.go
│   ├── sendRequestToURL.go
│   └── simulateTimeout.go
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── engine/
│   ├── dynamicconfig/
│   └── docker-compose.yml
├── internal/
│   ├── cronjobs/
│   │   ├── config.go
│   │   └── sendEmail.go
│   ├── services/
│   │   ├── sandwichShop.go
│   │   ├── simulationToSendEmail.go
│   │   └── timeoutProcess.go
│   └── zapadapter/
│       └── zapadapter.go
├── signals/
│   ├── receiver.go
│   └── sendPayment.go
├── starters/
│   ├── sandwichShop.go
│   ├── simulationToSendEmail.go
│   └── timeoutProcess.go
├── workflows/
│   ├── sandwichShop.go
│   ├── simulationToSendEmail.go
│   └── timeoutProcess.go
├── go.mod
└── go.sum
```
## Starting the Project: Running Workers and Services

To start the project after building the Docker image, you can execute the following commands in the terminal:

```bash
go run cmd/worker/main.go
go run cmd/server/main.go
```

## Details of the server/main.go file

The server file in the project is responsible for setting up and starting the HTTP server using the Echo framework. Here’s a brief overview of what this file does:

### Temporal Client Initialization:
The file begins by initializing a Temporal client (`client.Dial`) to connect to the Temporal server. This client is essential for interacting with Temporal workflows and activities.

### Cron Job Configuration:
The `ConfigureCronJobs` function from the cron jobs package is called, which sets up any recurring tasks or workflows using Temporal scheduling features.

### HTTP Server Setup with Echo:
The Echo framework is used to create a new HTTP server instance (`echo.New()`).

- A basic health check endpoint is set up at `/health-check`, which returns an “OK” response to confirm the server is running.

### Workflow Routes:
The server defines several routes under the `/start-workflow` path for triggering different Temporal workflows:

- **Sandwich Shop Workflows:**
  - `/place-order`: Triggers a workflow related to placing an order.
  - `/check-payment/:workflowID`: Checks the payment status for a given workflow.

- **Send Email Workflows:**
  - `/v1, /v2, /v3, /v4`: These routes trigger different versions of email sending workflows, each showcasing various Temporal features like linear processing, goroutines, batch processing, child workflows, and futures.
  - `/all-versions`: Triggers all versions of the email sending workflows sequentially and returns the result.

- **Timeout Workflow:**
  - `/timeout`: Triggers a workflow that simulates a timeout scenario.

### Starting the Server:
Finally, the server is started on port 3310 (`e.Start(":3310")`), and it listens for incoming HTTP requests.

### Summary
The server file serves as the entry point for HTTP requests in the Temporal project. It sets up routes to trigger different workflows, configures cron jobs, and initializes the necessary Temporal client to interact with the Temporal backend. This setup allows users to start workflows, check their statuses, and simulate various scenarios like timeouts, all through a simple HTTP API.

## Details of the worker/main.go file

The worker file in the project is responsible for setting up and managing multiple Temporal workers that handle different workflows and activities. Here’s a summary of what this file does:

### Logging Configuration:
The file begins by configuring a Zap logger with a specific logging level (`ErrorLevel`). This logger is used to capture and manage logs generated by the Temporal workers.

### Temporal Client Initialization:
A Temporal client is created (`client.Dial`) to connect with the Temporal server. This client is essential for the workers to execute workflows and activities.

### Graceful Shutdown Handling:
The program sets up signal handling for graceful shutdowns. It listens for OS interrupts (e.g., Ctrl+C) to cleanly stop the workers without losing any in-progress tasks.

### Worker Definitions:
The `configureWorkers` function defines multiple workers, each assigned to a specific Task Queue. Each worker is configured with the workflows and activities it will handle.

### Task Queues and Assigned Workflows/Activities:
- **send-email-from-cron-job:** Handles workflows for sending emails triggered by cron jobs.
- **sandwich-shop:** Manages workflows related to a sandwich shop, including making, borrowing, and serving sandwiches.
- **send-email-with-linear-process:** Executes workflows that send emails using a linear process.
- **send-email-with-goroutine-and-batch-process:** Handles workflows that send emails using goroutines and batch processing.
- **send-email-with-child-and-batch-process:** Manages workflows that involve child workflows and batch processing.
- **send-email-with-child-and-future-and-batch-process:** Handles workflows with a combination of child workflows, futures, and batch processing.
- **simulate-timeout:** Executes workflows designed to simulate a timeout scenario.

### Worker Execution:
Each worker is launched in a separate goroutine, ensuring that they run concurrently. The workers listen for tasks in their assigned Task Queues and execute the registered workflows and activities. The main program waits for an interrupt signal to shut down the workers gracefully.

### Summary
The worker file is a critical component of a project, responsible for setting up and running Temporal workers that execute various workflows and activities. It organizes the workers by Task Queues, each designed to handle specific types of tasks, such as email sending, sandwich shop operations, or timeout simulations. This structure allows for clear separation of responsibilities and efficient task processing within the Temporal system.

## Details of services

### health-check
```bash
curl --location 'http://localhost:3310/health-check'
```
To check if the server is up and running, you can use the following curl command. If the server is operational, it will return a 200 / OK response

### simulate-timeout
```bash
curl --location 'http://localhost:3310/start-workflow/timeout' 
```
ervice calls the starters.SimulateTimeout function, which returns a 200 status.
The starters.SimulateTimeout function will execute the workflows.SimulateTimeout workflow. If there is a need for any preliminary checks before the workflow is executed, they can be performed within the starters function.

The workflows.SimulateTimeout function configures and runs the activities.SimulateTimeout activity with ActivityOptions. The workflow will not be marked as completed (status: completed) until it returns nil. If there are no errors and the retry policies have not terminated the workflow, it will be displayed as status: running in the web UI.
In the workflows.SimulateTimeout function, it is specified that the workflow must be terminated within 3 seconds of starting.

However, within the activities.SimulateTimeout, there is a 10-second sleep period.
As the activity does not return a response within the specified 3 seconds, the workflow is terminated and can be viewed in the web UI with a status: failed.
This flow is prepared as an example of how to use timeouts to limit the execution of workflows.

### place-order
```bash
curl --location 'http://localhost:3310/start-workflow/sandwich-shop/place-order'
```
When a request is made to the place-order service, it sequentially prepares the necessary environment for Temporal through the services, starters, and workflows files.
The SandwichShop workflow is executed. This workflow consists of three different activities and does not close until it returns nil.

The first activity is triggered directly (MakeSandwich). A message indicating that the sandwich is being prepared is logged in the terminal. A 3-second sleep period is set. It is important to use the workflow.Sleep function within the workflow in Temporal; activities cannot be delayed using time.Sleep within a workflow.

Following this, the signals.ReceiveSignal function is called with a 10-second timeout.
The ReceiveSignal function, located in the signals file, initiates the signal process by creating a signal channel, timer, and selector. If no signal is received within 10 seconds, paymentStatus is set to false, and the function concludes. If a signal is received within 10 seconds, paymentStatus is set to true, and the function concludes.
Below is the service created for signals. The Temporal workflow ID must be added as a query parameter to this service’s endpoint.

### check-payment
```bash
curl --location 'http://localhost:3310/start-workflow/sandwich-shop/check-payment/{workflowID}'
```
If the signal is received as “true” within 10 seconds, the activities.ServeSandwich activity is executed, and the user is informed via the terminal that the order has been delivered after payment.

If no signal is received and the 10 seconds elapse, the activities.BorrowSandwich activity is executed, and the user is informed via the terminal that the order has been delivered without payment.

The workflow is then terminated.
This structure demonstrates the use of signals. Workflows can be paused, frozen, or terminated at a specific step.

### send-email/v1
```bash
curl --location 'http://localhost:3310/start-workflow/send-email/v1' 
```
The v1 endpoint simulates sending emails to 100 sequentially generated dummy email addresses by calling an activity for each email within a for loop, resulting in 100 separate activities. If an error occurs within the activity used to send the emails (e.g., a panic), the retryPolicy mechanism is triggered, which will retry the operation according to the configured settings. The errors can also be logged or subjected to additional processing if needed.

The v1 endpoint operates linearly and is the longest-running email sending workflow due to its sequential nature.

### send-email/v2
```bash
curl --location 'http://localhost:3310/start-workflow/send-email/v2' 
```
The v2 endpoint takes 100 sequentially generated dummy email addresses and groups them into batches of 10. Each batch is then processed using Temporal GoRoutines via the workflow.Go function.

This results in the creation of 10 separate goroutines, each responsible for sequentially calling an activity to send emails for its respective batch.
Progress is monitored through a channel created for this purpose, which tracks the completion of each activity. Once all activities are completed, the workflow is terminated. If any issues arise within a batch during the email sending process, the retryPolicy mechanism is triggered, allowing for the desired flow to be applied.

This setup is prepared as an example of using Temporal GoRoutines. Execution times can be monitored via terminal logs.

### send-email/v3
```bash
curl --location 'http://localhost:3310/start-workflow/send-email/v3' 
```
The v3 endpoint takes 100 sequentially generated dummy email addresses and groups them into batches of 10. These batches are then sent as a whole to a childWorkflow.

This is a linear process where each childWorkflow is completed before moving on to the next batch. The childWorkflow unpacks the batch and triggers 10 separate sendEmail activities.

This structure serves as an example of using ChildWorkflow. It allows the process to be broken down into smaller, more manageable parts with better retryPolicy implementation, making it easier to manage.

When a childWorkflow is completed, the desired value can be passed back to the parentWorkflow.

### send-email/v4
```bash
curl --location 'http://localhost:3310/start-workflow/send-email/v4' 
```
The v4 endpoint takes 100 sequentially generated dummy email addresses and groups them into batches of 10. For each batch, a child workflow (SimulateSendEmailChildWorkflow) is triggered using workflow.ExecuteChildWorkflow.

The future object returned by this function represents the eventual result of the child workflow. This future does not block the execution; instead, it allows the main workflow to continue processing other tasks or triggering more child workflows.
Each future is appended to the futures slice, effectively queuing all child workflows for concurrent execution.

After all the futures are collected, the code enters a second loop where it waits for each future to complete using future.Get(ctx, &result). This method blocks until the corresponding child workflow completes and retrieves the result from it.
The total count of successfully processed emails is accumulated by adding up the results from all futures.

The duration of the entire workflow and memory usage are measured and printed to the console.

The v4 endpoint was created to demonstrate the combined use of child workflows, futures, and batch processing in Temporal. 

### send-email/all-versions
```bash
curl --location 'http://localhost:3310/start-workflow/send-email/all-versions’ 
```
The all-versions endpoint is designed to trigger all the versions of the send-email workflows and provide a convenient way to view the results in the terminal. This endpoint allows for the simultaneous execution of various implementations of the email sending process, making it easier to compare their performance and behaviour. 

It serves as a useful tool for testing and debugging, enabling developers to assess different workflow versions in one go and examine the outcomes efficiently in the terminal logs.

### send-email with cronjob
The server folder is initialised by the main.go file, where the cronjobs.ConfigureCronJobs function is triggered. The sendEmail structure is set to run every minute. To guide users on different configurations, an example constant CronSchedule has been defined.

The sendEmail job is executed once when the machine starts, and it logs the last job status and execution time to the terminal. Additionally, the SimulateSendEmailWithCronJob workflow is configured as a cron job with the necessary settings.

The SimulateSendEmailWithCronJob is responsible for executing the sendRequestToURL activity, which triggers the send-email/all-versions endpoint. This allows all versions’ outputs to be viewed in the terminal.

This setup was prepared as an example of how to use Temporal Cronjobs.

### Some Outputs in Terminal (all-versions)
```bash
SimulateSendEmailWithTemporalChildWorkflowAndFutureAndBatchProcess | 100 emails sent and batches processed in 1.648173708s
SimulateSendEmailWithTemporalGoroutineAndBatchProcess | 100 emails sent and batches processed in 3.594585167s
SimulateSendEmailWithLinearProcess | 100 emails sent and batches processed in 9.055019417s
SimulateSendEmailWithTemporalChildWorkflowAndBatchProcess | 100 emails sent and batches processed in 8.823133667s
```
