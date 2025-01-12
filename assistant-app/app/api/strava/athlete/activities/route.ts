//Get athletes activities from Strava
export async function GET(request: Request) {
    const url = new URL(request.url);
    const page = url.searchParams.get('page') || '1';
    const perPage = url.searchParams.get('per_page') || '10';

    try {
        const response = await fetch(`http://localhost:8080/athlete/activities?page=${page}&perPage=${perPage}`);
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