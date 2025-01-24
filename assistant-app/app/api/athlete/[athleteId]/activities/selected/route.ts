import { NextResponse } from "next/server";
import { cookies } from 'next/headers';

export async function GET(
  request: Request,
  { params }: { params: { athleteId: string } }
) {
  try {
    // Get selected activities from cookies (server side)
    const cookieStore = await cookies();
    const selectedActivitiesStr = cookieStore.get('selectedActivities')?.value;
    
    if (!selectedActivitiesStr) {
      return NextResponse.json(
        { error: "No selected activities found" },
        { status: 400 }
      );
    }

    const selectedActivities = JSON.parse(selectedActivitiesStr);
    const activityIds = selectedActivities.map((activity: { id: number }) => activity.id).join(',');

    const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";
    const response = await fetch(
      `${backendUrl}/athlete/${params.athleteId}/selected-activities?activityIds=${activityIds}`,
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      throw new Error(`Backend returned ${response.status}: ${response.statusText}`);
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error("Error fetching selected activities:", error);
    return NextResponse.json(
      { error: "Failed to fetch selected activities" },
      { status: 500 }
    );
  }
}
