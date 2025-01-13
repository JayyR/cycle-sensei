# cycle-sensei

Unlock Your Cycling Potential with Cycle Sensei, the AI-powered cycling coach that adapts to your unique needs.

## Overview
cycle-sensei is a web application that helps athletes track their performance. It consists of two main modules:
1. `assistant-app`: A Next.js frontend application.
2. `backend-app`: A Go backend application.

## Features
- **Athlete Performance Tracking**: Track and analyze your cycling performance.
- **Strava Integration**: Sync your data with Strava.
- **Personalized Coaching**: Receive personalized coaching tips based on your performance data.
- **Real-time Analytics**: Get real-time analytics and insights on your cycling activities.

## Prerequisites
- Node.js (for `assistant-app`)
- Go (for `backend-app`)
- Docker (optional, for running services in containers)

## Running the Modules

### assistant-app (Next.js)

1. Navigate to the `assistant-app` directory:
    ```sh
    cd assistant-app
    ```

2. Install the dependencies:
    ```sh
    npm install
    ```

3. Run the development server:
    ```sh
    npm run dev
    ```

4. Open your browser and go to `http://localhost:3000` to see the application.

### backend-app (Go)

1. Navigate to the `backend-app` directory:
    ```sh
    cd backend-app
    ```

2. Install the dependencies:
    ```sh
    go mod tidy
    ```

3. Run the backend server:
    ```sh
    go run main.go
    ```

4. The backend server will start on `http://localhost:8080`.

## Docker (Optional)

You can also run the services using Docker.

### Running assistant-app with Docker

1. Build the Docker image:
    ```sh
    docker build -t assistant-app .
    ```

2. Run the Docker container:
    ```sh
    docker run -p 3000:3000 assistant-app
    ```

3. Open your browser and go to `http://localhost:3000` to see the application.

### Running backend-app with Docker

1. Build the Docker image:
    ```sh
    docker build -t backend-app .
    ```

2. Run the Docker container:
    ```sh
    docker run -p 8080:8080 backend-app
    ```

3. The backend server will start on `http://localhost:8080`.

## License
This project is licensed under the MIT License.
