# Everato API Integration

This repository contains the **Everato API**, a generalized and decoupled API integration for **Everato**. The API is designed to work independently from the **Data Handler** and **Notification Service**, ensuring modularity and scalability.

## Overview

The **Everato API** serves as the central communication hub for various components within the Everato ecosystem. It facilitates interactions between the **Main Dashboard**, **Event Dashboards**, **Metrics & Logging**, and backend services like the **Data Handler** and **Notification Service**.

### Key Features

- **Event Creation**: After creating an event, the API communicates with the **Admin Dashboard** and **Event Dashboard** for further processing.
- **Event Bus Integration**: Publishes and subscribes to events using **Kafka** and **Zookeeper**.
- **gRPC Communication**: Handles gRPC calls to acknowledge or return values to connected services.
- **Database Communication**: Ensures seamless interaction with the database for storing and retrieving data.
- **Metrics & Logging**: Provides detailed metrics and logging for monitoring and debugging.

## Architecture

The architecture is designed to ensure decoupled and efficient communication between components:

1. **Main Dashboard**: Interacts with the API for general operations.
2. **Event Dashboards**: Receives updates and communicates with the API after an event is created.
3. **Event Bus**: Utilizes **Kafka** and **Zookeeper** for event publishing and subscription.
4. **Data Handler**: Subscribes to events from the Event Bus and interacts with the database.
5. **Notification Service**: Subscribes to events from the Event Bus to send notifications.
6. **Metrics & Logging**: Monitors API operations and logs relevant data.

## Tech Stack

- **Programming Language**: Go
- **Event Bus**: Kafka & Zookeeper
- **Communication Protocol**: gRPC
- **Dashboards**: Admin Dashboard and Event Dashboard
- **Database**: Integrated with the Data Handler

## Workflow

1. **Event Creation**:

    - The API communicates with the **Admin Dashboard** and **Event Dashboard**.
    - Further communication with the database is handled via the **Data Handler**.

2. **Event Bus**:

    - The API publishes events to the Event Bus.
    - The **Data Handler** and **Notification Service** subscribe to these events.

3. **Metrics & Logging**:
    - The API sends metrics and logs to the monitoring system.

## Diagram

Below is the visual representation of the architecture:

![Everato API Integration Diagram](assets/plan_revision_00.png)

## Getting Started

To set up the API, follow these steps:

1. Clone the repository.
2. Install dependencies.
3. Configure Kafka and Zookeeper for the Event Bus.
4. Run the API service.

## Contributing

Contributions are welcome! Please follow the guidelines in the `CONTRIBUTING.md` file.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
