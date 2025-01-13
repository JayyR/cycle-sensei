import { fetchWithAuth } from "@/utils/fetchWithAuth";

//Get athletes activities from Strava
export async function GET(request: Request, context: { params: { athleteId: string } }) {
    const { athleteId } = context.params;
    const url = new URL(request.url);
    const page = url.searchParams.get('page') || '1';
    const perPage = url.searchParams.get('per_page') || '10';

    try {
        const data = await fetchWithAuth(`http://localhost:8080/athlete/${athleteId}/activities?page=${page}&perPage=${perPage}`,request);
        
        return new Response(JSON.stringify(data), {
            headers: { 'Content-Type': 'application/json' },
        });
    } catch (error) {
        return new Response(null, { status: 500, statusText: error.message });
    }
}