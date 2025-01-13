"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, User } from "@nextui-org/react";
import Image from 'next/image';
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faLocationDot } from "@fortawesome/free-solid-svg-icons";

const StravaAuth = () => {
    const router = useRouter();
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [athlete, setAthlete] = useState({ firstName: '', lastName: '', city: '', state: '', image: '' });
    const [isLoading, setIsLoading] = useState(true);

    const handleStravaAuth = async () => {
        setIsLoading(true);
        const redirectUri = process.env.NEXT_PUBLIC_STRAVA_REDIRECT_URI;
        const currentPath = window.location.pathname;
        sessionStorage.setItem('originalPath', currentPath);
        const scope = process.env.NEXT_PUBLIC_STRAVA_SCOPE;
        window.location.href = `https://www.strava.com/oauth/authorize?client_id=${process.env.NEXT_PUBLIC_STRAVA_CLIENT_ID}&response_type=code&redirect_uri=${redirectUri}&scope=${scope}`;
    };

    const handleLogout = async () => {
        try {
            const response = await fetch('/api/strava/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (response.ok) {
                setIsAuthenticated(false);
                sessionStorage.removeItem('stravaAthlete');
                sessionStorage.removeItem('stravaAthleteId'); // Clear athlete ID from session
                router.refresh();
            } else {
                console.error('Failed to log out');
            }
        } catch (error) {
            console.error('Error during logout:', error);
        }
    };

    useEffect(() => {
        const fetchAthleteInfo = async () => {
            try {
                const cachedAthlete = sessionStorage.getItem('stravaAthlete');
                if (cachedAthlete) {
                    const parsedAthlete = JSON.parse(cachedAthlete);
                    setAthlete({ 
                        firstName: parsedAthlete.firstname, 
                        lastName: parsedAthlete.lastname, 
                        city: parsedAthlete.city, 
                        state: parsedAthlete.state, 
                        image: parsedAthlete.profile_medium 
                    });
                    sessionStorage.setItem('stravaAthleteId', parsedAthlete.id); // Store athlete ID in session
                    return;
                }

                const response = await fetch('/api/strava/athlete', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    if (data.firstname && data.lastname && data.city && data.state && data.profile_medium) {
                        setAthlete({ 
                            firstName: data.firstname, 
                            lastName: data.lastname, 
                            city: data.city, 
                            state: data.state, 
                            image: data.profile_medium 
                        });
                        sessionStorage.setItem('stravaAthlete', JSON.stringify(data));
                        sessionStorage.setItem('stravaAthleteId', data.id); // Store athlete ID in session
                    } else {
                        console.error('Incomplete athlete data:', data);
                    }
                } else {
                    console.error('Failed to fetch athlete info');
                }
            } catch (error) {
                console.error('Error fetching athlete info:', error);
            }
        };

        const checkAuthentication = async () => {
            try {
                const response = await fetch('/api/strava/token', {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    if (data.accessToken) {
                        setIsAuthenticated(true);
                        fetchAthleteInfo();
                    }
                }
                setIsLoading(false);
            } catch (error) {
                console.error('Error checking authentication:', error);
                setIsLoading(false);
            }
        };

        const handleStravaAuth = async () => {
            const urlParams = new URLSearchParams(window.location.search);
            const code = urlParams.get('code');

            if (code) {
                try {
                    const response = await fetch('/api/strava/auth', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ code }),
                    });

                    if (response.ok) {
                        const data = await response.json();
                        document.cookie = `stravaAccessToken=${data.accessToken}; Path=/;`;
                        setIsAuthenticated(true);
                        const originalPath = sessionStorage.getItem('originalPath') || '/';
                        router.replace(originalPath);
                        fetchAthleteInfo();
                    } else {
                        console.error('Failed to exchange code for token');
                    }
                } catch (error) {
                    console.error('Error during token exchange:', error);
                } finally {
                    setIsLoading(false);
                }
            } else {
                checkAuthentication();
            }
        };

        handleStravaAuth();
    }, [router]);

    return (
        <div>
            {!isLoading && (
                isAuthenticated ? (
                    <div>
                        <Dropdown placement="bottom-start">
                            <DropdownTrigger>
                                <User
                                    as="button"
                                    avatarProps={{
                                        isBordered: true,
                                        src: athlete.image,
                                    }}
                                    className="transition-transform"
                                    description={
                                        <>
                                            <FontAwesomeIcon icon={faLocationDot} /> {`${athlete.city}, ${athlete.state}`}
                                        </>
                                    }
                                    name={`${athlete.firstName} ${athlete.lastName}`}
                                />
                            </DropdownTrigger>
                            <DropdownMenu aria-label="User Actions" variant="flat">
                                <DropdownItem key="profile" className="h-14 gap-2">
                                    <p className="font-bold">Signed in as</p>
                                    <p className="font-bold">{`${athlete.firstName} ${athlete.lastName}`}</p>
                                </DropdownItem>
                                <DropdownItem key="logout" color="danger" onPress={handleLogout}>
                                    Log Out
                                </DropdownItem>
                            </DropdownMenu>
                        </Dropdown>
                    </div>
                ) : (
                    <Image
                        alt="NextUI hero Image"
                        src="/btn_strava_connect.svg"
                        layout="intrinsic"
                        width={300}
                        height={50}
                        onClick={handleStravaAuth}
                        style={{ height: 'auto', width: 'auto', maxHeight: '100%' }}
                    />
                )
            )}
        </div>
    );
};

export default StravaAuth;

