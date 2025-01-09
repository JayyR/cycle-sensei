import React, { useEffect, useState } from "react";
import {Button} from "@nextui-org/button";

const StravaZones = () => {
    const [zones, setZones] = useState(null);

    useEffect(() => {
        const fetchZones = async () => {
            const storedZones = localStorage.getItem("stravaZones");
            const storedTime = localStorage.getItem("stravaZonesTime");
            const currentTime = new Date().getTime();

            if (storedZones && storedTime && currentTime - parseInt(storedTime) < 24 * 60 * 60 * 1000) {
                setZones(JSON.parse(storedZones));
            } else {
                try {
                    const response = await fetch("/api/strava/athlete/zones");
                    const data = await response.json();
                    localStorage.setItem("stravaZones", JSON.stringify(data));
                    localStorage.setItem("stravaZonesTime", currentTime.toString());
                    setZones(data);
                } catch (error) {
                    console.error("Error fetching zones:", error);
                }
            }
        };

        fetchZones();
    }, []);

    return (
        <div>
            {zones ? (
                <Button
                className="bg-gradient-to-tr from-pink-500 to-yellow-500 text-white shadow-lg"
                radius="full"
              >
                Button
              </Button>
            ) : (
                <div>Loading...</div>
            )}
        </div>
    );
}



export default StravaZones;