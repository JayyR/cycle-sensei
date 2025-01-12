//Get athletes activity from Strava
export async function GET(request: Request, { params: { activityId } }) {

    try {
        const response = await fetch(`http://localhost:8080/athlete/activities/${activityId}`);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        return new Response(JSON.stringify(data), {
            headers: { 'Content-Type': 'application/json' },
        });
    } catch (error) {
        return new Response(null, { status: 500, statusText: error.message });
    }
}
