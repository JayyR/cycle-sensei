"use client";

import React from "react";
import styles from "./page.module.css"; // use simple styles for demonstration purposes

import { CardBody, Card, Accordion, AccordionItem, Spacer } from "@nextui-org/react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faAnglesDown, faChartLine, faPersonBiking, faZ } from '@fortawesome/free-solid-svg-icons';
import Chat from "../components/chat/chat";
import AthleteActivities from "../components/athlete/activities";
import AthleteStats from "../components/athlete/stats";
import AthleteZones from "../components/athlete/zones";

const Home = () => {
  return (
    <main className={styles.main}>
      <div className={styles.container}>
        <div className={styles.column}>
          <Accordion variant="splitted" defaultExpandedKeys={["0"]}>

            <AccordionItem key="0" aria-label="Recent Activities"
              title={
                <div className="flex" style={{ justifyContent: "flex-start" }}>
                  <p>Recent Activities</p>
                  <Spacer x={4} />
                </div>
              }
              startContent={<FontAwesomeIcon icon={faPersonBiking} />}
              indicator={<FontAwesomeIcon icon={faAnglesDown} />}
            >
              <AthleteActivities />
            </AccordionItem>

            <AccordionItem key="1" aria-label="Athlete Stats"
              title={
                <div className="flex" style={{ justifyContent: "flex-start" }}>
                  <p>Stats</p>
                  <Spacer x={4} />
                </div>
              }
              startContent={<FontAwesomeIcon icon={faChartLine} />}
              indicator={<FontAwesomeIcon icon={faAnglesDown} />}
            >
              <Card>
                <CardBody>
                  <AthleteStats />
                </CardBody>
              </Card>
            </AccordionItem>


            <AccordionItem key="2" aria-label="Athlete Zones"
              title={
                <div className="flex" style={{ justifyContent: "flex-start" }}>
                  <p>Athlete Zones</p>
                  <Spacer x={4} />
                </div>
              }
              startContent={<FontAwesomeIcon icon={faZ} />}
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
              <AthleteZones />
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