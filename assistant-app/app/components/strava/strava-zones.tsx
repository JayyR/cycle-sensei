"use client";

import React, { useEffect } from "react";
import {
    Table,
    TableHeader,
    TableBody,
    TableColumn,
    TableRow,
    TableCell,
    Chip,
    Avatar
} from "@nextui-org/react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBoltLightning, faHeartPulse } from '@fortawesome/free-solid-svg-icons';

function calculateHRZones(zones) {
    const zoneNames = ["Endurance", "Moderate", "Tempo", "Threshold", "Anaerobic"];

    // Calculate zone percentages and assign to zones with zone number
    let prevMax = 0;

    const calculatedZones = zones.map((zone, index) => {
        const zoneNumber = index + 1;
        const zoneName = zoneNames[index];
        const minPercentage = prevMax;
        let maxPercentage;

        if (zone.max === -1) {
            maxPercentage = 100;
        } else {
            maxPercentage = zone.max;
        }

        prevMax = maxPercentage;

        return {
            zone: zoneNumber,
            name: zoneName,
            min: zone.min || 0,
            max: zone.max,
            minPercentage,
            maxPercentage
        };
    });

    // Attach zone info to each item
    return zones.map((zone, index) => {
        return {
            ...zone,
            ...calculatedZones[index] // Add the calculated zone info
        };
    });
}

function calculatePowerZones(zones) {
    const zoneNames = [
        "Active Recovery",
        "Endurance",
        "Tempo",
        "Threshold",
        "VO2Max",
        "Anaerobic",
        "Neuromuscular",
    ];

    let prevMax = 0;

    const calculatedZones = zones.map((zone, index) => {
        const zoneNumber = index + 1;
        const zoneName = zoneNames[index];
        const minPercentage = prevMax;
        let maxPercentage;

        if (zone.max === -1) {
            maxPercentage = 100;
        } else {
            maxPercentage = zone.max;
        }

        prevMax = maxPercentage;

        return {
            zone: zoneNumber,
            name: zoneName,
            min: zone.min || 0,
            max: zone.max,
            minPercentage,
            maxPercentage,
        };
    });

    return zones.map((zone, index) => {
        return {
            ...zone,
            ...calculatedZones[index],
        };
    });
}

const HRZonesTable = ({ zones }) => {
    const hrZones = calculateHRZones(zones);

    return (
        <Table aria-label="Heart Rate Zones">
            <TableHeader>
                <TableColumn>Zone</TableColumn>
                <TableColumn>Name</TableColumn>
                <TableColumn>Min HR</TableColumn>
                <TableColumn>Max HR</TableColumn>
                <TableColumn>Min %</TableColumn>
                <TableColumn>Max %</TableColumn>
            </TableHeader>
            <TableBody>
                {hrZones.map((zone) => (
                    <TableRow key={zone.zone}>
                        <TableCell>{zone.zone}</TableCell>
                        <TableCell>{zone.name}</TableCell>
                        <TableCell>{zone.min}</TableCell>
                        <TableCell>{zone.max === -1 ? "Max" : zone.max}</TableCell>
                        <TableCell>{zone.minPercentage}%</TableCell>
                        <TableCell>{zone.maxPercentage}%</TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}

const PowerZonesTable = ({ zones }) => {
    const powerZones = calculatePowerZones(zones);

    return (
        <Table aria-label="Power Zones">
            <TableHeader>
                <TableColumn>Zone</TableColumn>
                <TableColumn>Name</TableColumn>
                <TableColumn>Min Power</TableColumn>
                <TableColumn>Max Power</TableColumn>
                <TableColumn>Min %</TableColumn>
                <TableColumn>Max %</TableColumn>
            </TableHeader>
            <TableBody>
                {powerZones.map((zone) => (
                    <TableRow key={zone.zone}>
                        <TableCell>{zone.zone}</TableCell>
                        <TableCell>{zone.name}</TableCell>
                        <TableCell>{zone.min}</TableCell>
                        <TableCell>{zone.max === -1 ? "Max" : zone.max}</TableCell>
                        <TableCell>{zone.minPercentage}%</TableCell>
                        <TableCell>{zone.maxPercentage}%</TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}

const StravaZones = () => {

    const [zones, setZones] = React.useState(null);

    useEffect(() => {
        const fetchZones = async () => {
            const currentTime = Date.now();
            const response = await fetch("/api/strava/athlete/zones");
            const data = await response.json();
            localStorage.setItem("stravaZones", JSON.stringify(data));
            localStorage.setItem("stravaZonesTime", currentTime.toString());
            setZones(data);
        };

        const storedZones = localStorage.getItem("stravaZones");
        const storedTime = localStorage.getItem("stravaZonesTime");
        const isExpired = storedTime && (Date.now() - parseInt(storedTime, 10)) > 7 * 24 * 60 * 60 * 1000;

        if (storedZones && !isExpired) {
            setZones(JSON.parse(storedZones));
        } else {
            fetchZones();
        }

    }, []);

    return (
        <div>
            {zones ? (
                <div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px', flex: 1 }}>
                        <div style={{ flex: 1 }}>
                            <Chip radius="sm" color="danger"
                                avatar={<Avatar icon={<FontAwesomeIcon icon={faHeartPulse} size="lg" style={{ color: "#ff0000", }} />} />}
                            >Heart Rate</Chip>
                            <div style={{ marginTop: '10px' }}></div>
                            <HRZonesTable zones={zones.heart_rate.zones} />
                        </div>
                        <div style={{ flex: 1 }}>
                            <Chip radius="sm" color="warning"
                                avatar={<Avatar icon={<FontAwesomeIcon icon={faBoltLightning} size="lg" />} />}
                            >Power </Chip>
                            <div style={{ marginTop: '10px' }}></div>
                            <PowerZonesTable zones={zones.power.zones} />
                        </div>
                    </div>
                </div>

            ) : (
                <div>Loading...</div>
            )}
        </div>
    );
}



export default StravaZones;