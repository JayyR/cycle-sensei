"use client";

import React from "react";
import styles from "./page.module.css"; // use simple styles for demonstration purposes
import Chat from "../components/chat/chat";
import StravaProfile from "../components/strava/strava-profile";
import StravaZones from "../components/strava/strava-zones";

const Home = () => {
  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={styles.column}>
          <StravaProfile />
          <StravaZones />
        </div>

        <div className={styles.chatContainer}>
          <div className={styles.chat}>
            <Chat />
          </div>
        </div>

      </div>

    </main>
  );
};

export default Home;