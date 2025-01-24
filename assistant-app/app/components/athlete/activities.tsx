"use client";

import React, { useState, useEffect } from "react";
import useSWR from "swr";
import {
    Table,
    TableHeader,
    TableColumn,
    TableBody,
    TableRow,
    TableCell,
    Pagination,
    Spinner,
    Button,
    Dropdown,
    DropdownTrigger,
    DropdownMenu,
    DropdownItem,
    DateRangePicker,
    Chip,
} from "@nextui-org/react";
import AthleteActivity from "./activity";
import { faEye, faRefresh, faAngleDown } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { parseDate } from "@internationalized/date";

const fetcher = (url: string) => fetch(url).then((res) => res.json());

interface Activity {
    id: number;
    name: string;
    start_date: string;
    sport_type: string;
}

interface SelectedActivity {
    id: number;
    name: string;
}

// Utility function to format the date
const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { weekday: 'short', year: 'numeric', month: 'numeric', day: 'numeric' };
    return new Date(dateString).toLocaleDateString(undefined, options);
};


const AthleteActivities = () => {
    const [page, setPage] = useState(1);
    const [selectedActivityId, setSelectedActivityId] = useState<number | null>(null);
    const [selectedActivities, setSelectedActivities] = useState<SelectedActivity[]>([]);
    const [selectedSportType, setSelectedSportType] = useState<string | null>(null);
    const itemsPerPage = -1; // Load all activities from backend
    const athleteId = sessionStorage.getItem("stravaAthleteId");
    const { data, isLoading, mutate } = useSWR(`/api/athlete/${athleteId}/activities?page=1&per_page=${itemsPerPage}`, fetcher, {
        keepPreviousData: true,
    });

    const loadingState = isLoading || data?.length === 0 ? "loading" : "idle";
    const pages = Math.ceil((data?.length ?? 0) / 10); // Local pagination with 10 items per page

    const [dateRange, setDateRange] = useState<{ startDate: Date | null, endDate: Date | null }>({ startDate: null, endDate: null });

    // Filter the data by sport_type and date range using useMemo
    const filteredData = React.useMemo(() => {
        return (data ?? []).filter((item: Activity) => {
            const validSportTypes = ["VirtualRide", "Ride", "GravelRide", "MountainBikeRide"];
            const isSportTypeValid = validSportTypes.includes(item.sport_type);
            const matchesSportType = selectedSportType ? item.sport_type === selectedSportType : true;
            const isDateInRange = dateRange.startDate && dateRange.endDate
                ? new Date(item.start_date) >= dateRange.startDate && new Date(item.start_date) <= dateRange.endDate
                : true;
            return isSportTypeValid && matchesSportType && isDateInRange;
        });
    }, [data, dateRange, selectedSportType]);
    console.log(filteredData);

    const [isOpen, setIsOpen] = React.useState(false);

    const minDate = data?.length ? new Date(data[data.length - 1].start_date) : null;
    const maxDate = data?.length ? new Date(data[0].start_date) : null;

    const minCalendarDate = minDate ? parseDate(minDate.toISOString().split('T')[0]) : null;
    const maxCalendarDate = maxDate ? parseDate(maxDate.toISOString().split('T')[0]) : null;

    const handleRefresh = async () => {
        try {
            const response = await fetch(`/api/athlete/${athleteId}/activities`, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${sessionStorage.getItem("stravaAuthToken")}`,
                },
            });

            if (response.ok) {
                mutate(); // Refresh the data
            } else {
                console.error("Failed to refresh activities");
            }
        } catch (error) {
            console.error("Error refreshing activities:", error);
        }
    };

    const handleView = (id: number) => {
        setSelectedActivityId(id);
        setIsOpen(true);
    };

    useEffect(() => {
        const storedSelections = sessionStorage.getItem("selectedActivities");
        if (storedSelections) {
            const parsedSelections: SelectedActivity[] = JSON.parse(storedSelections);
            setSelectedActivities(parsedSelections);
            setSelectedKeys(new Set(parsedSelections.map((activity) => activity.id.toString())));
        }
    }, []);

    useEffect(() => {
        sessionStorage.setItem("selectedActivities", JSON.stringify(selectedActivities));
        // Also store in cookies for server-side access
        document.cookie = `selectedActivities=${JSON.stringify(selectedActivities)}; path=/`;
    }, [selectedActivities]);

    const [selectedKeys, setSelectedKeys] = React.useState(new Set([]));

    const handleSelectionChange = (keys: "all" | Set<React.Key>) => {
        let selectedKeysSet: Set<React.Key>;
        if (keys === "all") {
            selectedKeysSet = new Set(data?.map((item: Activity) => item.id.toString()) ?? []);
        } else {
            selectedKeysSet = keys;
        }

        setSelectedKeys(selectedKeysSet);

        const selected = Array.from(selectedKeysSet).map((key) => {
            const activity = data?.find((item: Activity) => item.id.toString() === key);
            return activity ? { id: activity.id, name: activity.name } : null;
        }).filter(Boolean) as SelectedActivity[];
        setSelectedActivities(selected);
    };

    const handleClearSelection = () => {
        setSelectedKeys(new Set());
        setSelectedActivities([]);
    };


    const handleSelectionByDays = (days: number) => {
        const now = dateRange.endDate ? new Date(dateRange.endDate) : new Date();
        const pastDate = new Date(now);
        pastDate.setDate(now.getDate() - days);

        const selectedKeysSet = new Set(
            (filteredData ?? [])
                .filter((item: Activity) => new Date(item.start_date) >= pastDate)
                .map((item: Activity) => item.id.toString())
        );

        setSelectedKeys(selectedKeysSet);

        const selected = Array.from(selectedKeysSet).map((key) => {
            const activity = filteredData?.find((item: Activity) => item.id.toString() === key);
            return activity ? { id: activity.id, name: activity.name } : null;
        }).filter(Boolean) as SelectedActivity[];
        setSelectedActivities(selected);
    };

    const paginatedData = React.useMemo(() => {
        const startIndex = (page - 1) * 10;
        const endIndex = startIndex + 10;
        return filteredData.slice(startIndex, endIndex);
    }, [filteredData, page]);

    const sportTypeAbbreviations: { [key: string]: string } = {
        VirtualRide: 'VR',
        Ride: 'R',
        GravelRide: 'GR',
        MountainBikeRide: 'MTB'
    };

    const renderCell = React.useCallback((item: Activity, columnKey: React.Key) => {
        let cellContent;
        if (columnKey === "start_date") {
            cellContent = formatDate(item.start_date);
        } else if (columnKey === "actions") {
            cellContent = (
                <div style={{ display: "flex", gap: "8px" }}>
                    <Button
                        key={`view-${item.id}`}
                        color="primary"
                        isIconOnly
                        endContent={<FontAwesomeIcon icon={faEye} />}
                        onPress={() => handleView(item.id)}
                        size="sm"
                    />
                </div>
            );
        } else if (columnKey === "sport_type") {
            cellContent = sportTypeAbbreviations[item.sport_type] || item.sport_type;
        } else {
            cellContent = item[columnKey as keyof Activity];
        }
        return <TableCell>{cellContent}</TableCell>;
    }, []);

    return (
        <>
            <Table isStriped
                removeWrapper
                aria-label="Example table with client async pagination"
                bottomContent={
                    pages > 0 ? (
                        <div className="flex w-full justify-center">
                            <Pagination
                                isCompact
                                showControls
                                showShadow
                                color="primary"
                                page={page}
                                total={pages}
                                onChange={(page) => setPage(page)}
                            />
                        </div>
                    ) : null
                }
                topContent={
                    <div className="flex flex-col gap-4 w-full" > {/* Added w-full */}
                        <div className="flex flex-wrap justify-between gap-3 items-end"> {/* Added flex-wrap */}
                            <div className="flex gap-3">
                                <Dropdown aria-label="Select Activities within">
                                    <DropdownTrigger>
                                        <Button
                                            endContent={<FontAwesomeIcon icon={faAngleDown} />}
                                            size="sm"
                                            variant="bordered"
                                        >Select Activities within</Button>
                                    </DropdownTrigger>
                                    <DropdownMenu aria-label="Date Group Selection"
                                        onAction={(key) => {
                                            if (key === "clearAll") {
                                                handleClearSelection();
                                            } else if (key === "7days") {
                                                handleSelectionByDays(7);
                                            } else if (key === "30days") {
                                                handleSelectionByDays(30);
                                            } else if (key === "60days") {
                                                handleSelectionByDays(60);
                                            } else if (key === "90days") {
                                                handleSelectionByDays(90);
                                            }
                                        }}
                                        variant="bordered"
                                    >
                                        <DropdownItem key="7days">7 Days</DropdownItem>
                                        <DropdownItem key="30days">30 Days</DropdownItem>
                                        <DropdownItem key="60days">60 Days</DropdownItem>
                                        <DropdownItem key="90days">90 Days</DropdownItem>
                                        <DropdownItem key="clearAll" color="warning">
                                            Clear selection
                                        </DropdownItem>
                                    </DropdownMenu>
                                </Dropdown>
                                
                                <Dropdown aria-label="Filter by Sport Type">
                                    <DropdownTrigger>
                                        <Button
                                            endContent={<FontAwesomeIcon icon={faAngleDown} />}
                                            size="sm"
                                            variant="bordered"
                                        >
                                            {selectedSportType ? sportTypeAbbreviations[selectedSportType] : "All Types"}
                                        </Button>
                                    </DropdownTrigger>
                                    <DropdownMenu 
                                        aria-label="Sport Type Selection"
                                        variant="bordered"
                                        onAction={(key) => {
                                            if (key === "all") {
                                                setSelectedSportType(null);
                                            } else {
                                                setSelectedSportType(key.toString());
                                            }
                                        }}
                                    >
                                        <DropdownItem key="all">All Types</DropdownItem>
                                        <DropdownItem key="VirtualRide">Virtual Ride (VR)</DropdownItem>
                                        <DropdownItem key="Ride">Ride (R)</DropdownItem>
                                        <DropdownItem key="GravelRide">Gravel Ride (GR)</DropdownItem>
                                        <DropdownItem key="MountainBikeRide">Mountain Bike (MTB)</DropdownItem>
                                    </DropdownMenu>
                                </Dropdown>
                                
                                <Chip color="success" variant="light">
                                    {selectedActivities.length} selected
                                </Chip>
                            </div>

                            <div className="flex flex-wrap gap-3 items-center"> {/* Added items-center */}
                                <div className="flex items-center gap-2"> {/* New container to keep items in same row */}
                                    <DateRangePicker
                                        label="Season"
                                        pageBehavior="single"
                                        visibleMonths={3}
                                        variant="underlined"
                                        size="sm"
                                        minValue={minCalendarDate}
                                        maxValue={maxCalendarDate}
                                        onChange={(range) => {
                                            const startDate = range?.start ? new Date(range.start.toString()) : null;
                                            const endDate = range?.end ? new Date(range.end.toString()) : null;
                                            setDateRange({ startDate, endDate });
                                        }}
                                    />
                                    <Button color="primary"
                                        isIconOnly
                                        endContent={<FontAwesomeIcon icon={faRefresh} />}
                                        onPress={handleRefresh} 
                                    />
                                </div>
                            </div>
                        </div>
                    </div>
                }
                selectionMode="multiple"
                onSelectionChange={handleSelectionChange}
                color="success"
                selectedKeys={selectedKeys}
            >
                <TableHeader>
                    <TableColumn key="start_date">Date</TableColumn>
                    <TableColumn key="name">Title</TableColumn>
                    <TableColumn key="sport_type">Type</TableColumn>
                    <TableColumn key="actions">Actions</TableColumn>
                </TableHeader>
                <TableBody
                    items={paginatedData}
                    loadingContent={<Spinner />}
                    loadingState={loadingState}
                    emptyContent={"No Activities found"}

                >
                    {(item: Activity) => (
                        <TableRow key={item.id} >
                            {(columnKey: string) => renderCell(item, columnKey)}
                        </TableRow>
                    )}
                </TableBody>
            </Table>
            {Boolean(selectedActivityId) && (
                <AthleteActivity id={selectedActivityId} isOpen={isOpen} onOpenChange={setIsOpen} />
            )}
        </>

    );
};

export default AthleteActivities;