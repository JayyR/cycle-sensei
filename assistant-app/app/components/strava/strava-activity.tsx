"use client";

import { faAnglesRight } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Card, CardBody, Drawer, DrawerBody, DrawerContent, DrawerFooter, DrawerHeader, Tooltip } from "@nextui-org/react";
import { useEffect, useState } from "react";

const StravaActivity = ({ id, isOpen, onOpenChange }) => {
    const [activity, setActivity] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchActivity = async () => {
            setLoading(true); // Reset loading state to true when id changes
            try {
                const response = await fetch(`/api/strava/athlete/activities/${id}`);
                const data = await response.json();
                setActivity(data);
            } catch (error) {
                console.error("Error fetching activity:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchActivity();
    }, [id]);



    return (
        <>
            {loading && (
                <div className="loading-overlay">
                    <div className="loading-spinner">Loading...</div>
                </div>
            )}
            <Drawer isOpen={isOpen} onOpenChange={onOpenChange}
                backdrop="blur"
                classNames={{
                    base: "data-[placement=right]:sm:m-2 data-[placement=left]:sm:m-2  rounded-medium",
                }}

            >
                <DrawerContent>
                    {(onClose) => (
                        <>
                            <DrawerHeader className="absolute top-0 inset-x-0 z-50 flex flex-row gap-2 px-2 py-2 border-b border-default-200/50 justify-between bg-content1/50 backdrop-saturate-150 backdrop-blur-lg">
                                <Tooltip content="Close">
                                    <Button
                                        isIconOnly
                                        className="text-default-400"
                                        size="sm"
                                        variant="light"
                                        onPress={onClose}
                                    >
                                        <FontAwesomeIcon icon={faAnglesRight} />
                                    </Button>
                                </Tooltip>
                            </DrawerHeader>
                            <DrawerBody className="pt-16">
                                <Card>
                                    <CardBody>
                                        {activity ? (
                                            <div>
                                                {/* Render the activity data here */}
                                                <pre>{JSON.stringify(activity, null, 2)}</pre>
                                            </div>
                                        ) : (
                                            <div>No activity data available</div>
                                        )}
                                    </CardBody>
                                </Card>
                            </DrawerBody>
                            <DrawerFooter className="flex flex-col gap-1">
                                <Button color="danger" variant="light" onPress={onClose}>
                                    Close
                                </Button>
        
                            </DrawerFooter>
                        </>
                    )}
                </DrawerContent>
            </Drawer>
        </>
    );
};

export default StravaActivity;