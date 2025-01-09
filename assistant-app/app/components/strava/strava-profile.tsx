import React, { useEffect } from "react";
import { User, Skeleton } from "@nextui-org/react";



const StravaProfile = () => {

  //load athlete data from Strava when component mounts
  const [athlete, setAthlete] = React.useState(null);

  useEffect(() => {
    const fetchAthlete = async () => {
      const response = await fetch("/api/strava/athlete");
      const data = await response.json();
      localStorage.setItem("athlete", JSON.stringify(data));
      localStorage.setItem("athleteTimestamp", Date.now().toString());
      setAthlete(data);
      console.log(data);
    };

    const athleteData = localStorage.getItem("athlete");
    const athleteTimestamp = localStorage.getItem("athleteTimestamp");
    const isExpired = athleteTimestamp && (Date.now() - parseInt(athleteTimestamp, 10)) > 24 * 60 * 60 * 1000;

    if (athleteData && !isExpired) {
      setAthlete(JSON.parse(athleteData));
    } else {
      fetchAthlete();
    }
  }, []);


  return (
    <div>

      {athlete && (
        <div className="max-w-[300px] w-full flex items-center gap-3">
          <User
            avatarProps={{
              size: "lg",
              src: "https://dgalywyr863hv.cloudfront.net/pictures/athletes/41803353/19718900/8/large.jpg",
            }}
            description={`${athlete.city}, ${athlete.state}`}
            name={`${athlete.firstname} ${athlete.lastname}`}
          />
        </div>
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

export default StravaProfile;