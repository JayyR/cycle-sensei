import { fetchWithAuth } from "@/utils/fetchWithAuth";

//Get athletes zones from Strava
export async function GET(request: Request, context: { params: { athleteId: string } }) {
    const { athleteId } = context.params;
    try {
        const data = await fetchWithAuth(`http://localhost:8080/athlete/${athleteId}/zones`,request);
        return Response.json(data);
    } catch (error) {
       return new Response(null, { status: 500, statusText: error.message });
    }
}