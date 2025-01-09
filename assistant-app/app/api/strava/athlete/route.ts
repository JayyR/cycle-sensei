
// Get athletes from Strava
export async function GET() {
    try {
        const response = await fetch('http://localhost:8080/athlete');
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        return Response.json(data);
    } catch (error) {
       return new Response(null, { status: 500, statusText: error.message });
    }
}