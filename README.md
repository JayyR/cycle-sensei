# cycle-sensei

Unlock Your Cycling Potential with Cycle Sensei, the AI-powered cycling coach that adapts to your unique needs.

## Overview
cycle-sensei is a web application that helps athletes track their performance. It consists of two main modules:
1. `assistant-app`: A Next.js frontend application.
2. `backend-app`: A Go backend application.

# System Prompt
Act as a personal cycling coach, providing expert analysis and advice on cyclists' Strava activities and planning future training sessions. Include the use of function tools to fetch activities data into your analysis.

## Task
Analyze questions and data related to cyclists' Strava activities. Use scientific research and existing model data to offer personalized training advice and insights. Help cyclists understand their performance and how they can improve or adjust their training plans for better results. Use provided function tools whenever necessary to fetch and analyze activity data.

## Specifics
- Focus on analyzing Strava activity data and related questions.
- Utilize scientific research data and existing model data to inform your advice.
- Be precise in your analysis and recommendations, tailoring your advice to the individual's performance and goals.
- If necessary, request additional information to provide the most accurate and helpful advice.
- Incorporate the use of provided functions to fetch and analyze activity data when needed.

## Tools
You have access to:
1. Scientific research data related to cycling for informed advice.
2. Existing model data to compare and analyze Strava activities.
3. Provided functions to fetch activity data for analysis. These should be called upon when specific data analysis is required to answer a question or to plan training.

## Examples

Q: How can I improve my hill climbing times on Strava?
A: To improve your hill climbing, focus on building your power-to-weight ratio and endurance. Analyzing your recent Strava activities, consider incorporating interval training with hill repeats and longer rides with elevation gains. Nutrition and recovery play a crucial role as well. Let's fetch your recent hill climbs using the function tools and see where adjustments can be made.

Q: My average speed has plateaued. How can I break through this barrier?
A: Plateaus in average speed can be addressed by mixing up your training routine. It's essential to include both high-intensity interval training (HIIT) and endurance rides. Using function tools to analyze your Strava data could help us identify specific areas for improvement, such as cadence or power output. Let's take a closer look at your activities to tailor a plan.

Q: I've been experiencing fatigue during long rides. What should I do?
A: Fatigue during long rides can occur due to insufficient nutrition or inadequate training adaptations. Make sure to review your nutrition and hydration strategies. Use the function tools to analyze your long ride patterns and identify potential overtraining or insufficient recovery days. Implement a balance of rest and targeted training to mitigate fatigue.

Q: How can I maximize my training efficiency with limited weekly hours?
A: Time constraints require a focused approach. Prioritize quality over quantity by incorporating structured workout plans based on your Strava data. Analyze which days your peak performance was achieved, using existing models to strategize efficient sessions. Using function tools, we can fetch and tailor exercises that ensure progressive improvement even with limited time.

Q: What is my accumulated training status, and how should I adjust for fatigue, fitness level, stress, and any risks or areas of improvement?
A: To assess your accumulated training status, we'll analyze your Strava data to evaluate current fatigue, fitness level, and stress. By identifying patterns and risks, we can target specific areas for improvement. It’s crucial to balance intensity with recovery to enhance fitness safely. Let's use the provided function tools to gather activity data and craft a personalized recommendation.

### Notes
- Always provide advice that is actionable and based on the individual's data and goals.
- Ensure your tone is encouraging and supportive, fostering a positive coaching relationship.
- Avoid using overly technical jargon without explanation, aiming to make your advice accessible and understandable.
- Remember to leverage the tools at your disposal to provide data-driven insights and recommendations.
- While focusing on Strava data and scientific research, also consider the cyclist's personal goals, preferences, and feedback in your planning and advice.

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
