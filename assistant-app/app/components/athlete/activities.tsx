"use client";

import React, { useState } from "react";
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

} from "@nextui-org/react";
import AthleteActivity from "./activity";

const fetcher = (url: string) => fetch(url).then((res) => res.json());

interface Activity {
    id: number;
    name: string;
    start_date: string;
    sport_type: string;
}

// Utility function to format the date
const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { weekday: 'short', year: 'numeric', month: 'numeric', day: 'numeric' };
    return new Date(dateString).toLocaleDateString(undefined, options);
};


const AthleteActivities = () => {
    const [page, setPage] = useState(1);
    const [selectedActivityId, setSelectedActivityId] = useState<number | null>(null);
    const itemsPerPage = 10;
    const athleteId = sessionStorage.getItem("stravaAthleteId");
    const { data, isLoading } = useSWR(`/api/athlete/${athleteId}/activities?page=${page}&per_page=${itemsPerPage}`, fetcher, {
        keepPreviousData: true,
    });
    console.log(data);
    const loadingState = isLoading || data?.length === 0 ? "loading" : "idle";
    const pages = 10;

    // Filter the data by sport_type using useMemo
    const filteredData = React.useMemo(() => {
        return (data ?? []).filter((item: Activity) =>
            ["VirtualRide", "Ride", "GravelRide", "MountainBikeRide"].includes(item.sport_type)
        );
    }, [data]);
    // Function to handle row actions
    const handleRowAction = (key: React.Key) => {
        setSelectedActivityId(Number(key));
        setIsOpen(true)
    };
    const [isOpen, setIsOpen] = React.useState(false);

    return (
        <>
            <Table
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
                onRowAction={handleRowAction}
            >
                <TableHeader>
                    <TableColumn key="start_date">Date</TableColumn>
                    <TableColumn key="name">Title</TableColumn>
                    <TableColumn key="sport_type">Type</TableColumn>
                </TableHeader>
                <TableBody
                    items={filteredData ?? []}
                    loadingContent={<Spinner />}
                    loadingState={loadingState}
                >
                    {(item: Activity) => (
                        <TableRow key={item.id}>
                            {(columnKey) => (
                                <TableCell>
                                    {columnKey === "start_date" ? formatDate(item.start_date) : item[columnKey]}
                                </TableCell>
                            )}
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