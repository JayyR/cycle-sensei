import React, { useEffect } from "react";
import { Skeleton, Tabs, Card, CardBody, Tab } from "@nextui-org/react";

const AthleteStats = () => {

  //load athlete data from Strava when component mounts
  const [athlete, setAthlete] = React.useState(null);

  useEffect(() => {
    const fetchAthlete = async () => {
      const athleteId = sessionStorage.getItem("stravaAthleteId");
      const response = await fetch(`/api/athlete/${athleteId}`);
      const data = await response.json();
      localStorage.setItem("athleteStats", JSON.stringify(data));
      localStorage.setItem("athleteStatsTimestamp", Date.now().toString());
      setAthlete(data);
      console.log(data);
    };

    const athleteData = localStorage.getItem("athleteStats");
    const athleteTimestamp = localStorage.getItem("athleteStatsTimestamp");
    const isExpired = athleteTimestamp && (Date.now() - parseInt(athleteTimestamp, 10)) > 7 * 24 * 60 * 60 * 1000;

    if (athleteData && !isExpired) {
      setAthlete(JSON.parse(athleteData));
    } else {
      fetchAthlete();
    }
  }, []);


  return (
    <div>
      {athlete && (
        <Tabs aria-label="Options">
          <Tab key="weekly" title="Weekly">
            <Card>
              <CardBody>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
                incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud
                exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
              </CardBody>
            </Card>
          </Tab>
          <Tab key="yeartodate" title="Year to Date">
            <Card>
              <CardBody>
                Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex
                ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse
                cillum dolore eu fugiat nulla pariatur.
              </CardBody>
            </Card>
          </Tab>
          <Tab key="alltime" title="All Time">
            <Card>
              <CardBody>
                Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt
                mollit anim id est laborum.
              </CardBody>
            </Card>
          </Tab>
        </Tabs>
      )}

      {!athlete && (
        <div className="max-w-[300px] w-full flex items-center gap-3">
          <div>
            <Skeleton className="flex rounded-full w-12 h-12" />
          </div>
          <div className="w-full flex flex-col gap-2">
            <Skeleton className="h-3 w-3/5 rounded-lg" />
            <Skeleton className="h-3 w-4/5 rounded-lg" />
          </div>
        </div>
      )}
    </div>
  );
};

export default AthleteStats;