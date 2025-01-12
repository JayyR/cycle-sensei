"use client";

import React from "react";
import styles from "./page.module.css"; // use simple styles for demonstration purposes

import { CardBody, Card, Accordion, AccordionItem, Spacer } from "@nextui-org/react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faAnglesDown, faPersonBiking, faZ } from '@fortawesome/free-solid-svg-icons';
import Chat from "../components/chat/chat";
import StravaProfile from "../components/strava/strava-profile";
import StravaZones from "../components/strava/strava-zones";
import StravaActivities from "../components/strava/strava-activities";

const Home = () => {
  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={styles.column}>
          <Card>
            <CardBody>
              <StravaProfile />
            </CardBody>
          </Card>
          <Accordion variant="splitted" defaultExpandedKeys={["2"]}>
            <AccordionItem key="1" aria-label="Athlete Zones"
              title={
                <div className="flex" style={{ justifyContent: "flex-start" }}>
                  <p>Athlete Zones</p>
                  <Spacer x={4} />
                </div>
              }
              startContent={<FontAwesomeIcon icon={faZ} beat />}
              subtitle={
                <div className="flex" style={{ justifyContent: "space-between", width: "30%" }}>
                  <p>
                    Max HR <span className="text-primary ml-0">191</span>
                  </p>
                  <Spacer x={4} />
                  <p>
                    FTP <span className="text-primary ml-0">234</span>
                  </p>
                </div>
              }
              indicator={<FontAwesomeIcon icon={faAnglesDown} />}
            >
              <StravaZones />
            </AccordionItem>

            <AccordionItem key="2" aria-label="Recent Activities" 
            title={
              <div className="flex" style={{ justifyContent: "flex-start" }}>
                <p>Recent Activities</p>
                <Spacer x={4} />
              </div>
            }
            startContent={<FontAwesomeIcon icon={faPersonBiking} beat />}
            indicator={<FontAwesomeIcon icon={faAnglesDown} />}>
              <StravaActivities />
            </AccordionItem>
          </Accordion>
        </div>
        <div className={styles.chatContainer}>
          <Card className={styles.chat}>
            <CardBody>
              <Chat />
            </CardBody>
          </Card>
        </div>
      </div>
    </main>
  );
};

export default Home;