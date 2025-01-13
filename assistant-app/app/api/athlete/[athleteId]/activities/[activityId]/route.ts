import { fetchWithAuth } from "@/utils/fetchWithAuth";

//Get athletes activity from Strava
export async function GET(request: Request,context: { params: { athleteId: string , activityId:string} }) {
    const { athleteId, activityId } = context.params;
    try {
        const data = await fetchWithAuth(`http://localhost:8080/athlete/${athleteId}/activities/${activityId}`,request);
        return new Response(JSON.stringify(data), {
            headers: { 'Content-Type': 'application/json' },
        });
    } catch (error) {
        return new Response(null, { status: 500, statusText: error.message });
    }
}
